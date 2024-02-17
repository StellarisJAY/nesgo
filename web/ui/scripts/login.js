function login() {
    const name = document.getElementById("username").value
    const password = document.getElementById("password").value
    post("/user/login", JSON.stringify({name: name, password: password}))
        .then(data=>{
        setToken(data["token"])
        window.location.href = "/home"
         })
        .catch(error=>{
            console.log(error)
        })
}