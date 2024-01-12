async function loginRequest(name, password) {
    const data = {name: name, password: password}
    const resp =  await post("/user/login", JSON.stringify(data));
    if (!resp.ok) {
        throw new Error("login request error")
    }
    return resp.json()
}

function login() {
    const name = document.getElementById("username").value
    const password = document.getElementById("password").value
    loginRequest(name, password)
        .then(resp=>{
            setToken(resp.data["token"])
            window.location.href = "/room/1"
        })
        .catch(error=>{
            console.log(error)
        })
}