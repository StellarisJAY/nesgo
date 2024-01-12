function getToken() {
    return sessionStorage.getItem("access-token")
}

function setToken(token) {
    sessionStorage.setItem("access-token", token)
}

function getAuthorizedHeader() {
    return {"Authorization": getToken()}
}

