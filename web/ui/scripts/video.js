let rtcSession = {
}

let roomId

let configs = {
    controlButtonMapping: {
        "button-up": "Up",
        "button-down": "Down",
        "button-left": "Left",
        "button-right": "Right",
        "button-a": "A",
        "button-b": "B",
        "button-select": "Select",
        "button-start": "Start",
    },
    keyboardMapping: {
        "KeyA": "Left",
        "KeyD": "Right",
        "KeyW": "Up",
        "KeyS": "Down",
        "KeyJ": "B",
        "Space": "A",
        "Enter": "Start",
        "Tab": "Select",
    },
    existingGames: {},
}

const MessageSDPOffer = 0
const MessageSDPAnswer = 1
const MessageICECandidate = 2
const MessageGameButtonPressed = 3
const MessageGameButtonReleased = 4

const MemberTypeOwner = 0
const MemberTypeGamer = 1
const MemberTypeWatcher = 2

onload = ev=>{
    roomId = window.location.pathname.substring(6)
    // 连接之前禁用控制按钮
    setControlButtonsDisabled(true)
    getRoomMemberSelf()
    listGames()
}

function connect() {
    const connectButton = document.getElementById("connect-button")
    connectButton.disabled = true
    const selectedGame = document.getElementById("select-game").value;
    const ws = new WebSocket(websocketAddr+"/room/"+roomId+"/rtc?auth=" + getToken() + "&game="+selectedGame)
    ws.onopen = ev=> {
        const pc = new RTCPeerConnection({
            iceServers: [
                {urls: stunServer},
                {urls: turnServer, username: turnUser, credential: turnCredential}
            ],
            iceTransportPolicy: "relay",
        })
        pc.onicecandidate = ev => {
            if (ev.candidate !== null) {
                ws.send(JSON.stringify({
                    "type": MessageICECandidate,
                    "data": btoa(JSON.stringify(ev.candidate))
                }))
                pc.addIceCandidate(ev.candidate).then(_ => console.log("ice candidate: ", ev.candidate))
            }else {
                console.log(ev)
            }
        }

        pc.onconnectionstatechange = ev=>{
            console.log("peer conn state: " + pc.connectionState)
            switch (pc.connectionState) {
                case "connected":
                    break
                case "disconnected":
                    pc.close()
                    connectButton.disabled = false
                    break
                default:
                    break
            }
        }

        pc.oniceconnectionstatechange = ev=>{
            console.log("ice conn state: " + pc.iceConnectionState)
        }

        pc.ontrack = ev=>{
            console.log("track id:" + ev.streams[0].id)
            document.getElementById("video").srcObject = ev.streams[0]
            document.getElementById("video").autoplay = true
            document.getElementById("video").controls = true
        }
        rtcSession.pc = pc
    }

    ws.onclose = ev => {
        window.onkeydown = _=>{}
        window.onkeyup = _ => {}
        alert("websocket connection closed")
        document.getElementById("connect-button").disabled = false
    }

    ws.onerror = ev => {
        console.log(ev)
        ws.close()
    }

    ws.onmessage = ev => {
        const message = JSON.parse(ev.data)
        if (message["type"] === MessageSDPOffer) {
            const sdpOffer = JSON.parse(atob(message["data"]))
            const pc = rtcSession.pc
            pc.setRemoteDescription(sdpOffer)
                .then(_ => pc.createAnswer())
                .then(sdp => pc.setLocalDescription(sdp))
                .then(_ => {
                    console.log("remote sdp: ", pc.remoteDescription)
                    console.log("local sdl:  ", pc.localDescription)
                    ws.send(JSON.stringify({
                        "type": MessageSDPAnswer,
                        "data": btoa(JSON.stringify(pc.localDescription)),
                    }))
                    pc.addTransceiver("video")
                })
                .then(_=>onConnected())
                .catch(err=>{
                    console.log(err)
                })
        }
    }
    rtcSession.ws = ws
}

function restartEmulator() {
    const selectedGame = document.getElementById("select-game").value
    post("/room/"+roomId+"/restart?game=" + selectedGame, {})
        .catch(error=>{
            console.log(error)
        })
}

function sendAction(code, pressed) {
    rtcSession.ws.send(JSON.stringify({
        "type": pressed,
        "data": btoa(code),
    }))
}

function onConnected() {
    if (rtcSession.member["memberType"] !== MemberTypeWatcher) {
        window.onkeydown = ev=> {
            const button = configs.keyboardMapping[ev.code];
            if (button) {
                sendAction(button, MessageGameButtonPressed)
            }
        }

        window.onkeyup = ev=> {
            const button = configs.keyboardMapping[ev.code];
            if (button) {
                sendAction(button, MessageGameButtonReleased)
            }
        }
    }
    for (const id in configs.controlButtonMapping) {
        const button = document.getElementById(id)
        button.disabled = rtcSession.member["memberType"] === MemberTypeWatcher
        const code = configs.controlButtonMapping[id]
        button.addEventListener("mousedown", ()=>sendAction(code, MessageGameButtonPressed))
        button.addEventListener("mouseup", ()=>sendAction(code, MessageGameButtonReleased))
        button.addEventListener("touchstart", ()=>sendAction(code, MessageGameButtonPressed))
        button.addEventListener("touchend", ()=>sendAction(code, MessageGameButtonReleased))
    }

    document.getElementById("restart-button").disabled = false
}

function listGames() {
    get("/games", null)
        .then(data=>{
        const selector = document.getElementById("select-game");
        for (let game of data) {
            configs.existingGames[game.name] = game
            selector.innerHTML += "<option value=\"" + game.name + "\">" + game.name + "</option>"
        }
    })
        .catch(err => {
            console.log(err)
        })
}

function setControlButtonsDisabled(disabled) {
    for (const id in configs.controlButtonMapping) {
        document.getElementById(id).disabled = disabled
    }
}

function getRoomMemberSelf() {
    get("/room/"+roomId+"/member", null)
        .then(data=>{
            rtcSession.member = data
        })
        .catch(resp=>{
            if (resp.status === 403) {
                alert(resp.message)
            }else {
                console.log(resp.message)
            }
            window.location = "/home"
        })
}