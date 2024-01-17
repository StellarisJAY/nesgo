const baseURL = "http://192.168.0.107:8080"
const wsURL = "ws://192.168.0.107:8080"

async function request(path, method, data) {
    let args = {
        method: method,
        headers:{},
        body: data
    }
    if (getToken()) {
        args.headers["Authorization"] = getToken()
    }
    return await fetch(baseURL + path, args)
}

function post(path, data) {
    return request(path, "POST", data)
}

function get(path, data) {
    return request(path, "GET", data)
}