const Configs = {
    ApiServer: "http://localhost:8080/", // 设置后端API服务器地址
    StunServer: "stun:localhost:3478",   // 设置STUN服务器地址
    TurnServer: {
        Host: "turn:localhost:3478", // 设置TURN服务器地址
        Username: "turn_user",       // 设置TURN服务器用户名
        Password: "turn_password",   // 设置TURN服务器密码
    }
};

export default Configs;