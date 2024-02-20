const serverAddr = "192.168.0.107:8080"
const apiServer = "http://" + serverAddr + "/api"
const websocketAddr = "ws://" + serverAddr + "/ws"

const stunServer = "stun:192.168.0.107:3478"
// stun无法解决地址相关的NAT映射，需要部署TURN服务器中转
// deploy coturn: https://github.com/coturn/coturn
const turnServer = "turn:192.168.0.107:3478"
const turnUser = "xxjay"
const turnCredential = "123456"