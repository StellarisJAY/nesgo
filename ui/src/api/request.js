import axios from "axios";
import tokenStorage from "./token";
import { message } from "ant-design-vue";
import router from "../router/index"

const a = axios.create({
    baseURL: "http://localhost:8080/api",
})

a.interceptors.request.use(config=>{
    const token = tokenStorage.getToken()
    if (token) {
        config.headers.set("Authorization", token)
    }
    return config
})

a.interceptors.response.use(response=>{
    const resp = response.data
    if (resp.status === 200) {
        return resp
    }else if (resp.status >= 500) {
        message.error(resp.message)
        return Promise.reject(resp)
    }else {
        message.warn(resp.message)
        return Promise.reject(resp)
    }
    return resp
}, response=>{
    if (response.status === 401) {
        router.go("/login")
    }
    return Promise.reject(response)
})

const api = {
    axios: a,
    get(path) {
        return a.get(path)
    },
    post(path, data) {
        return a.post(path, data)
    }
}


export default api