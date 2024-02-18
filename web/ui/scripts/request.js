const baseURL = "http://192.168.0.107:8080"
const wsURL = "ws://192.168.0.107:8080/ws"
const api = baseURL + "/api"
async function request(path, method, data) {
    let args = {
        method: method,
        headers:{},
        body: data
    }
    if (getToken()) {
        args.headers["Authorization"] = getToken()
    }
    return fetch(api + path, args)
        .then(resp=>{
            if (resp.status === 403) {
                window.location = "/login"
                return
            }
            return resp.json()
        })
        .then(resp=>{
            if (resp.status === 200) {
                return resp.data
            }else {
                throw {"status": resp.status, "message": resp.message}
            }
        })
}

function post(path, data) {
    return request(path, "POST", data)
}

function get(path, data) {
    return request(path, "GET", data)
}