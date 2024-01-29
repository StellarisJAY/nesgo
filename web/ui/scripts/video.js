let rtcSession = {
}

let roomId

onload = ev=>{
    roomId = window.location.pathname.substring(6)
}

const MessageSDPOffer = 0
const MessageSDPAnswer = 1
const MessageICECandidate = 2
const MessageGameButtonPressed = 3
const MessageGameButtonReleased = 4

function log(msg) {
    document.getElementById("log").innerText += msg + "<br>"
}

function connect() {
    const ws = new WebSocket(wsURL+"/room/"+roomId+"/rtc?auth=" + getToken())
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
        pc.ontrack = ev=>{
            console.log(ev)
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
}

function onConnected(ws) {
    window.onkeydown = ev=> {
        const code = ev.code
        if (code === "KeyA" || code === "KeyD" || code === "KeyW" || code === "KeyS" ||
            code === "Space" || code === "Enter" || code === "KeyJ") {
            ws.send(JSON.stringify({
                "type": MessageGameButtonPressed,
                "data": btoa(code),
            }))
            roomProperties.ws.send(JSON.stringify({"KeyCode": code, "Action": 0}))
        }
    }

    window.onkeyup = ev=> {
        const code = ev.code
        if (code === "KeyA" || code === "KeyD" || code === "KeyW" || code === "KeyS" ||
            code === "Space" || code === "Enter" || code === "KeyJ") {
            ws.send(JSON.stringify({
                "type": MessageGameButtonReleased,
                "data": btoa(code),
            }))
        }
    }
}