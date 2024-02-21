function register() {
    const name = document.getElementById("input-name").value
    const pwd = document.getElementById("input-password").value
    const confirmPwd = document.getElementById("input-confirm-password").value
    if (name === "" || name.length > 16) {
        alert("invalid user name format")
        return
    }
    if (pwd.length < 6 || pwd.length > 32) {
        alert("password too long or too short")
        return
    }
    if (pwd !== confirmPwd) {
        alert("passwords input don't match")
        return
    }

    post("/user/register", JSON.stringify({"name": name, "password": pwd}))
        .then(_=>{
            alert("register success")
            window.location = "/login"
        })
        .catch(resp=>{
            if (resp.status === 400) {
                alert(resp.message)
            }
            console.log(resp)
        })
}