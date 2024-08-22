const Configs = {
    ApiServer: "http://localhost:8080/", // 设置后端API服务器地址
    StunServer: "stun:localhost:3478",   // 设置STUN服务器地址
    TurnServer: {
        Host: "turn:localhost:3478", // 设置TURN服务器地址
        Username: "turn_user",       // 设置TURN服务器用户名
        Password: "turn_password",   // 设置TURN服务器密码
    },
    defaultKeyboardSetting: {
        "id": "0",
        "name": "默认设置",
        "bindings": [
            {
                "emulatorKey": "Left",
                "emulatorKeyTranslated": "Left",
                "buttons": ["KeyA"],
            },
            {
                "emulatorKey": "Right",
                "emulatorKeyTranslated": "Right",
                "buttons": ["KeyD"],
            },
            {
                "emulatorKey": "Up",
                "emulatorKeyTranslated": "Up",
                "buttons": ["KeyW"],
            },
            {
                "emulatorKey": "Down",
                "emulatorKeyTranslated": "Down",
                "buttons": ["KeyS"],
            },
            {
                "emulatorKey": "A",
                "emulatorKeyTranslated": "A",
                "buttons": ["Space"],
            },
            {
                "emulatorKey": "B",
                "emulatorKeyTranslated": "B",
                "buttons": ["KeyJ"],
            },
            {
                "emulatorKey": "Start",
                "emulatorKeyTranslated": "Start",
                "buttons": ["Enter"],
            },
            {
                "emulatorKey": "Select",
                "emulatorKeyTranslated": "Select",
                "buttons": ["Tab"],
            },
        ]
    }
};

export default Configs;