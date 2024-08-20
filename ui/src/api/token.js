const tokenStorage = {
    getToken() {
        return sessionStorage.getItem("nesgo_token")
    },
    setToken(token) {
        sessionStorage.setItem("nesgo_token", token)
    },
    delToken() {
        sessionStorage.clear()
    }
}

export default tokenStorage