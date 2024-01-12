const baseURL = "http://localhost:8080"

async function request(path, method, data) {
    let args = {
        method: method,
        headers: {
            contentType: "application/json"
        },
        body: data
    }
    if (getToken() !== "") {
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