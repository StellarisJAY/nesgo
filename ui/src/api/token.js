const tokenStorage = {
    getToken() {
        return sessionStorage.getItem("nesgo_token")
    },
    setToken(token) {
        sessionStorage.setItem("nesgo_token", token)
    }
}

export default tokenStorage