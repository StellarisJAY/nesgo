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
const MessageTurnServerInfo = 5

const RoleHost = 0
const RoleGamer = 1
const RoleObserver = 2

let roomMembers = {}

onload = ev=>{
    roomId = window.location.pathname.substring(6)
    // 连接之前禁用控制按钮
    setControlButtonsDisabled(true)
    getRoomMemberSelf()
}

function connect() {
    const connectButton = document.getElementById("connect-button")
    connectButton.disabled = true
    const selectedGame = document.getElementById("select-game").value;
    const ws = new WebSocket(websocketAddr+"/room/"+roomId+"/rtc?auth=" + getToken() + "&game="+selectedGame)
    document.getElementById("connect-button").disabled = true

    ws.onclose = ev => {
        console.log("websocket connection closed")
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
        }else if (message["type"] === MessageTurnServerInfo) {
            const turnInfo = JSON.parse(atob(message["data"]))
            rtcSession.turnAddress = turnInfo["address"]
            rtcSession.turnUser = turnInfo["username"]
            rtcSession.turnPassword = turnInfo["password"]
            createPeerConnection()
        }
    }
    rtcSession.ws = ws
}

function createPeerConnection() {
    const ws = rtcSession.ws
    const pc = new RTCPeerConnection({
        iceServers: [
            {urls: stunServer},
            {urls: rtcSession.turnAddress, username: rtcSession.turnUser, credential: rtcSession.turnPassword },
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
                rtcSession.ws.close()
                break
            case "disconnected":
                pc.close()
                document.getElementById("connect-button").disabled = false
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

    pc.ondatachannel = ev=>{
        const datachannel = ev.channel
        rtcSession.dataChannel = datachannel
        datachannel.onclose = _=>{
            window.onkeydown = _=>{}
            window.onkeyup = _ => {}
        }
        datachannel.onerror = err=>{
            console.log(err)
        }
        datachannel.onmessage = msg=>{
            console.log("unexpected dataChannel message:", msg)
        }
    }
    rtcSession.pc = pc
}

function restartEmulator() {
    const selectedGame = document.getElementById("select-game").value
    post("/room/"+roomId+"/restart?game=" + selectedGame, {})
        .catch(error=>{
            console.log(error)
        })
}

function sendAction(code, pressed) {
    rtcSession.dataChannel.send(JSON.stringify({
        "type": pressed,
        "data": btoa(code),
    }))
}

function onConnected() {
    if (rtcSession.member["role"] !== RoleObserver) {
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
        button.disabled = rtcSession.member["role"] === RoleObserver
        const code = configs.controlButtonMapping[id]
        button.addEventListener("mousedown", ()=>sendAction(code, MessageGameButtonPressed))
        button.addEventListener("mouseup", ()=>sendAction(code, MessageGameButtonReleased))
        button.addEventListener("touchstart", ()=>sendAction(code, MessageGameButtonPressed))
        button.addEventListener("touchend", ()=>sendAction(code, MessageGameButtonReleased))
    }

    document.getElementById("restart-button").disabled = false
    document.getElementById("save-button").disabled = false
    document.getElementById("load-button").disabled = false
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
    document.getElementById("save-button").disabled = disabled
    document.getElementById("load-button").disabled = disabled
}

function getRoomMemberSelf() {
    get("/room/"+roomId+"/member", null)
        .then(data=>{
            rtcSession.member = data
            listRoomMembers()
            listGames()
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

function quickSave() {
    post("/room/" + roomId + "/quickSave", null)
        .then(_=>{
        })
        .catch(resp=>{
            if (resp.status === 500) {
                alert("Internal server error")
                console.log(resp.message)
                return
            }
            alert(resp.message)
        })
}

function quickLoad() {
    post("/room/" + roomId + "/quickLoad", null)
        .then(_=>{})
        .catch(resp=>{
            if (resp.status === 500) {
                alert("Internal server error")
                console.log(resp.message)
                return
            }
            alert(resp.message)
        })
}

function listRoomMembers() {
    return get("/room/"+roomId+"/members", null)
        .then(data=> {
            roomMembers = {}
            data.forEach(m=>{
                roomMembers[m["id"]] = m
            })
        })
        .catch(resp=>{
            console.log(resp)
        })
}

function showRoomMembersModal() {
    const modal = new bootstrap.Modal(document.getElementById("room-members-modal"), {
        keyboard: false
    })
    listRoomMembers()
        .then(_=>renderMembersTable())
        .then(_=>modal.show())
}

function renderMembersTable() {
    const rows = document.getElementById("member-rows")
    rows.innerHTML = ""
    const isHost = rtcSession.member["role"] === RoleHost
    for (let k in roomMembers) {
        const member = roomMembers[k];
        const row = document.createElement("tr")
        row.innerHTML += "<td>" + member["name"] + "</td>"
        const td1 = document.createElement("td")
        const td2 = document.createElement("td")
        const td3 = document.createElement("td")
        const td4 = document.createElement("td")

        const p1 = document.createElement("input")
        const p2 = document.createElement("input")
        const gamer = document.createElement("input")
        const observer = document.createElement("input")
        p1.type="checkbox"
        p2.type = "checkbox"
        gamer.type = "checkbox"
        observer.type = "checkbox"
        p1.disabled=!isHost
        p2.disabled=!isHost
        gamer.disabled = !isHost
        observer.disabled = !isHost
        p1.checked = member["player1"]
        p2.checked = member["player2"]
        gamer.checked = member["role"]  === RoleGamer
        observer.checked = member["role"] === RoleObserver
        if (isHost) {
            p1.onchange = _ => transferControl(member["id"], true, false)
            p2.onchange = _=>transferControl(member["id"], false, true)
            gamer.onchange = _=>alterRole(member["id"], RoleGamer)
            observer.onchange = _=>alterRole(member["id"], RoleObserver)
        }
        console.log(member)
        td1.appendChild(p1)
        td2.appendChild(p2)
        td3.appendChild(gamer)
        td4.appendChild(observer)
        row.append(td1, td2, td3, td4)
        const kickButton = document.createElement("button")
        kickButton.innerText = "Kick"
        kickButton.type = "button"
        kickButton.className = "btn btn-primary kick-button"
        kickButton.style.height = "80%"
        kickButton.disabled = rtcSession.member["role"] !== 0
        if (rtcSession.member["role"] === RoleHost) {
            kickButton.onclick = _=>kick(member["id"])
        }
        row.appendChild(kickButton)
        rows.appendChild(row)
    }
}

function kick(memberId) {
    post("/room/"+roomId+"/member/kick", JSON.stringify({
        "memberId": memberId,
        "kick": true
    }))
        .then(data=>{
            delete roomMembers[memberId]
            renderMembersTable()
        })
        .catch(resp=>{
            if (resp.status!==500) {
                alert(resp.message)
            }
            console.log(resp.message)
        })
}

function transferControl(memberId, control1, control2) {
    if (roomMembers[memberId]["role"] === RoleObserver) {
        alert("can not give control to observer")
        return
    }
    post("/room/"+roomId+"/control/transfer", JSON.stringify({
        "memberId": memberId,
        "setController1": control1,
        "setController2": control2,
    }))
        .then(data=>{
            return listRoomMembers()
        })
        .then(_=>renderMembersTable())
        .catch(resp=>{
            if (resp.status !== 500) {
                alert(resp.message)
            }
            console.log(resp.message)
        })
}

function alterRole(memberId, role) {
    if (roomMembers[memberId]["role"] === role) {
        return
    }
    post("/room/" + roomId + "/role", JSON.stringify({
        "memberId": memberId,
        "role": role,
    }))
        .then(_=>{return listRoomMembers()})
        .then(_=>renderMembersTable())
        .catch(resp=>{
            if (resp.status !== 500) {
                alert(resp.message)
            }
            console.log(resp.message)
        })
}