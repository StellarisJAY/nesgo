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
const ws = new WebSocket("ws://localhost:8080/ws/"+gameId+"?auth=" + getToken())
let gameState = "running"

const pauseIcon = "<svg t=\"1703662386725\" class=\"icon\" viewBox=\"0 0 1024 1024\" version=\"1.1\" " +
    "xmlns=\"http://www.w3.org/2000/svg\" " +
    "p-id=\"8374\" width=\"20\" height=\"20\">" +
    "<path d=\"M304 176h80v672h-80zM712 176h-64c-4.4 0-8 3.6-8 8v656c0 4.4 3.6 8 8 8h64c4.4 0 8-3.6 8-8V184c0-4.4-3.6-8-8-8z\" " +
    "p-id=\"8375\">" +
    "</path></svg>"

const resumeIcon = "<svg t=\"1703662618428\" class=\"icon\" viewBox=\"0 0 1024 1024\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\" p-id=\"8513\" " +
    "width=\"20\" height=\"20\">" +
    "<path d=\"M715.8 493.5L335 165.1c-14.2-12.2-35-1.2-35 18.5v656.8c0 19.7 20.8 30.7 35 18.5l380.8-328.4c10.9-9.4 10.9-27.6 0-37z\" fill=\"#333333\" p-id=\"8514\">" +
    "</path></svg>"

let lastFrame

onload = function () {
    if (getToken() === "") {
        alert("please login")
        window.location.href = "/login"
    }
    document.getElementById("resumePauseButton").innerHTML = pauseIcon
}

ws.onmessage = function (event) {
    // blob to frame
    let reader = new FileReader()
    reader.onloadend = function () {
        lastFrame = reader.result
        drawCompressedFrame(reader.result)
    }
    reader.readAsArrayBuffer(event.data)
}

onkeydown = function (event) {
    const code = event.code
    if (gameState === "paused") {
        return
    }
    if (code === "KeyA" || code === "KeyD" || code === "KeyW" || code === "KeyS" ||
        code === "Space" || code === "Enter" || code === "KeyJ") {
        ws.send(JSON.stringify({"KeyCode": code, "Action": 0}))
    }
}

onkeyup = function (event) {
    const code = event.code
    if (code === "KeyP") {
        resumeOrPauseGame()
    }
    if (code === "Backspace") {
        reverseOnce()
    }
    if (gameState === "paused") {
        return
    }
    if (code === "KeyA" || code === "KeyD" || code === "KeyW" || code === "KeyS" ||
        code === "Space" || code === "Enter" || code === "KeyJ") {
        ws.send(JSON.stringify({"KeyCode": code, "Action": 1}))
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

function incScale() {
    gameConfigs.scale = Math.min(maxScale, gameConfigs.scale + 1)
    c.width = width * gameConfigs.scale
    c.height = height * gameConfigs.scale
    drawCompressedFrame(lastFrame)
}

function decScale() {
    gameConfigs.scale = Math.max(minScale, gameConfigs.scale - 1)
    c.width = width * gameConfigs.scale
    c.height = height * gameConfigs.scale
    drawCompressedFrame(lastFrame)
}

function incCPURate() {
    let oldRate = gameConfigs.boostRate
    gameConfigs.boostRate = Math.min(5.0, gameConfigs.boostRate + 0.5)
    if (oldRate !== gameConfigs.boostRate) {
        boostCPU(gameConfigs.boostRate)
    }

}

function decCPURate() {
    let oldRate = gameConfigs.boostRate
    gameConfigs.boostRate = Math.max(1.0, gameConfigs.boostRate - 0.5)
    if (oldRate !== gameConfigs.boostRate) {
        boostCPU(gameConfigs.boostRate)
    }
}

function boostCPU(rate) {
    post("/game/"+gameId+"/boost/"+rate, null)
        .then(data=>{
            return data.json()
        })
        .then(res=>{
            document.getElementById("boost").innerText = "Boost: " + gameConfigs.boostRate
        })
        .catch(error=>{
            console.error(error)
        })
}

function resumeOrPauseGame() {
    switch (gameState) {
        case "running":
            pauseGame()
            break
        case "paused":
            resumeGame()
            break
    }
}

function pauseGame() {
    post("/game/"+gameId+"/pause", null)
        .then(res=>{
            if (res.ok) {
                gameState = "paused"
                document.getElementById("resumePauseButton").innerHTML = resumeIcon
            }
        })
        .catch(err=>{
            console.error(err)
        })
}

function resumeGame() {
    if (gameState !== "paused") {
        return
    }
    post("/game/"+gameId+"/resume", null)
        .then(res=>{
            if (res.ok) {
                gameState = "running"
                document.getElementById("resumePauseButton").innerHTML = pauseIcon
            }
        })
        .catch(err=>{
            console.error(err)
        })
}

function reverseOnce() {
    post("/game/"+gameId+"/reverse", null)
        .then(res=>{

        })
        .catch(err=>{
            console.error(err)
        })
}