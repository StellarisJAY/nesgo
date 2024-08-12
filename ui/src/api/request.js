import axios from "axios";
import tokenStorage from "./token";
import { message } from "ant-design-vue";
import router from "../router/index"

// const a = axios.create({
//     baseURL: "http://192.168.0.107:8080/api",
// })

const v1 = axios.create({
    baseURL: "http://localhost:8030",
})

// const webSocketAddr = "ws://192.168.0.107:8080/ws"

v1.interceptors.request.use(config=>{
    const token = tokenStorage.getToken()
    if (token) {
        config.headers.set("Authorization", "Bearer " + token)
    }
    return config
})

v1.interceptors.response.use(response=>{
    if (response.status === 200) {
        return response.data
    }else if (response.status === 401) {
        return router.push("/login")
    }else {
        message.error(response.data.message).then()
        return Promise.reject(response)
    }
})

// a.interceptors.request.use(config=>{
//     const token = tokenStorage.getToken()
//     if (token) {
//         config.headers.set("Authorization", token)
//     }
//     return config
// })
//
// a.interceptors.response.use(response=>{
//     const resp = response.data
//     if (resp.status === 200) {
//         return resp
//     }else if (resp.status >= 500) {
//         message.error(resp.message)
//         return Promise.reject(resp)
//     }else {
//         return Promise.reject(resp)
//     }
// }, response=>{
//     if (response.response.status === 401) {
//         return router.push("/login")
//     }
//     return Promise.reject(response)
// })

const api = {
    axios: v1,
    get(path) {
        return this.axios.get(path)
    },
    post(path, data) {
        return this.axios.post(path, data)
    },
    put(path, data) {
        return this.axios.put(path, data)
    },
    delete(path) {
        return this.axios.delete(path)
    }
}

export default api