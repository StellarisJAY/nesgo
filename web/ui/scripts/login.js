async function loginRequest(name, password) {
    const data = {name: name, password: password}
    const args = {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
            "Content-Type": "application/json"
        }
    }
    const resp  = await fetch("http://localhost:8080/user/login", args)
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
        })
        .catch(error=>{
            console.log(error)
        })
}