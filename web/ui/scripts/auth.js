function getToken() {
    return localStorage.getItem("access-token")
}

function setToken(token) {
    localStorage.setItem("access-token", token)
}

function getAuthorizedHeader() {
    return {"Authorization": getToken()}
}

