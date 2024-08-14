import axios from "axios";
import tokenStorage from "./token";
import router from "../router/index"


const v1 = axios.create({
    baseURL: "http://localhost:8050",
});

v1.interceptors.request.use(config=>{
    const token = tokenStorage.getToken()
    if (token) {
        config.headers.set("Authorization", "Bearer " + token)
    }
    return config
});

v1.interceptors.response.use(r=>{
    if (r.status && r.status === 200) return r.data;
    if (!r["response"]) return Promise.reject(r);
    const response = r["response"];
    if (response.status === 401) return router.push("/login");
    return Promise.reject(response.data);
});

function errorHandler(err) {
    const response = err["response"];
    if (response && response.status) {
        if (response.status === 401) return router.push("/login");
    }
    return Promise.reject(response);
}

const api = {
    axios: v1,
    get(path) {
        return this.axios.get(path).catch(err=>errorHandler(err));
    },
    post(path, data) {
        return this.axios.post(path, data).catch(err=>errorHandler(err));
    },
    put(path, data) {
        return this.axios.put(path, data).catch(err=>errorHandler(err));
    },
    delete(path) {
        return this.axios.delete(path).catch(err=>errorHandler(err));
    }
}

export default api