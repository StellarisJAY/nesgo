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

onload = ev=>{
    roomId = window.location.pathname.substring(6)
    listGames()
}

function connect() {
    const selectedGame = document.getElementById("select-game").value;
    console.log(selectedGame)
    const ws = new WebSocket(wsURL+"/room/"+roomId+"/rtc?auth=" + getToken() + "&game="+selectedGame)
    ws.onopen = ev=> {
        const pc = new RTCPeerConnection({
            iceServers: [
                {urls: 'stun:stun.l.google.com:19302'}
            ]
        })
        pc.onicecandidate = ev => {
            if (ev.candidate !== null) {
                ws.send(JSON.stringify({
                    "type": MessageICECandidate,
                    "data": btoa(JSON.stringify(ev.candidate))
                }))
            }
        }

        pc.onconnectionstatechange = ev=>{
            console.log("peer conn state: " + pc.connectionState)
            switch (pc.connectionState) {
                case "connected":
                    document.getElementById("connect-button").disabled = true
                    break
                case "disconnected":
                    pc.close()
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
        }
        rtcSession.pc = pc
    }

    ws.onclose = ev => {
        window.onkeydown = _=>{}
        window.onkeyup = _ => {}
    }

    ws.onerror = ev => {

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
                .then(_=>onConnected(ws))
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
        .then(resp=>{
            return resp.json()
        })
        .then(data=>{
            if (data.status === 200) {
            }else if (data.status === 500) {
                document.getElementById("select-game").disabled = true
                document.getElementById("start-game-button").disabled = true
            }else {
                alert(data.message)
            }
        })
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

function onConnected(ws) {
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

    for (const id in configs.controlButtonMapping) {
        const button = document.getElementById(id)
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
        .then(resp=>{
            if (resp.status === 403) {
                window.location = "/login"
            }
            return resp.json()
        }).then(res=> {
        if (res.status !== 200) {
            throw new Error(res.message)
        }
        const selector = document.getElementById("select-game");
        for (let game of res.data) {
            configs.existingGames[game.name] = game
            selector.innerHTML += "<option value=\"" + game.name + "\">" + game.name + "</option>"
        }

        selector.addEventListener("change", ()=>{
        })
    })
        .catch(err => {
            console.log(err)
        })
}