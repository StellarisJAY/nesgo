let roomProperties = {
    id: -1,
    game: "SuperMario.nes",
}

const c = document.getElementsByTagName("canvas").item(0)
const gameId = c.id
const maxScale = 3
const minScale = 1
const width = 256
const height = 240

let gameConfigs = {
    boostRate: 1.0,
    scale: 3,
}
c.width = width * gameConfigs.scale
c.height = height * gameConfigs.scale
let ctx = c.getContext('2d')

onload = function (ev) {
    console.log(window.location.pathname)
    roomProperties.id = window.location.pathname.substring(6)
    getRoomInfo()
}

function connect() {
    if (roomProperties.id === -1) {
        return
    }
    roomProperties.ws = new WebSocket("ws://localhost:8080/ws/room/"+roomProperties.id+"?auth=" + getToken())
    roomProperties.ws.onerror = function(event) {
        console.log(event)
    }
    roomProperties.ws.onmessage = function (event) {
        // blob to frame
        let reader = new FileReader()
        reader.onloadend = function () {
            lastFrame = reader.result
            drawCompressedFrame(reader.result)
        }
        reader.readAsArrayBuffer(event.data)
    }
}

function disconnect() {

}

function startGame() {
    if (roomProperties.id === -1) {
        return
    }
    post("/room/"+roomProperties.id+"/start?game=" + roomProperties.game, {})
        .then(resp=>{
            return resp.json()
        })
        .then(data=>{

        })
        .catch(error=>{
            console.log(error)
        })
}

function getRoomInfo() {
    get("/room/"+roomProperties.id+"/info", null)
        .then(resp=>{
            if (resp.status === 403) {
                window.location = "/login"
            }
            return resp.json()
        })
        .then(result=>{
            console.log(result.status)
            if (result.status !== 200) {
                throw new Error(result.message)
            }else {
                roomProperties.info = result.data
                document.getElementById("roomId").innerText = roomProperties.id
                document.getElementById("roomName").innerText = roomProperties.info["name"]
                document.getElementById("roomInviteCode").innerText = roomProperties.info["inviteCode"]
            }
        })
        .catch(err=>{
            alert(err)
        })
}

onkeydown = function (event) {
    const code = event.code
    if (code === "KeyA" || code === "KeyD" || code === "KeyW" || code === "KeyS" ||
        code === "Space" || code === "Enter" || code === "KeyJ") {
        roomProperties.ws.send(JSON.stringify({"KeyCode": code, "Action": 0}))
    }
}

onkeyup = function (event) {
    const code = event.code
    if (code === "KeyA" || code === "KeyD" || code === "KeyW" || code === "KeyS" ||
        code === "Space" || code === "Enter" || code === "KeyJ") {
        roomProperties.ws.send(JSON.stringify({"KeyCode": code, "Action": 1}))
    }
}

function drawCompressedFrame(frameData) {
    const view = new DataView(frameData)
    let imageData = ctx.getImageData(0, 0, width * gameConfigs.scale, height * gameConfigs.scale)
    const frameSize = width * height
    for (let i = 0; i < frameSize; i++) {
        let colorId = view.getUint8(i)
        setScaledPixel(i, view.getUint8(frameSize + colorId*3),
            view.getUint8(frameSize + colorId*3 + 1),
            view.getUint8(frameSize + colorId*3 + 2),
            imageData)
    }
    ctx.putImageData(imageData, 0, 0)
}

function setScaledPixel(pixelNum, r,g,b, imageData) {
    let scale = gameConfigs.scale
    let row = Math.floor(pixelNum / width)
    let col = pixelNum % width
    let x0 = col * scale
    let y0 = row * scale
    for (let p = 0; p < scale; p++) {
        for (let q = 0; q < scale; q++){
            let index = xyToFrameIndex(x0+p,y0+q)
            let i = index * 4
            imageData.data[i] = r
            imageData.data[i+1] = g
            imageData.data[i+2] = b
            imageData.data[i+3] = 255
        }
    }
}

function xyToFrameIndex(x, y) {
    let w = width * gameConfigs.scale
    return x + y * w
}