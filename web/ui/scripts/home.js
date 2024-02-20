let joinedRooms = []

onload = ()=> {
    listJoinedRooms(1, 10)
    listRooms(1, 10)

    document.getElementById("refresh-all-rooms-btn").onclick = _ => listRooms(1, 10)
    document.getElementById("refresh-joined-rooms-btn").onclick = _=>listJoinedRooms(1, 10)
}

memberTypes = ["Owner", "Gamer", "Watcher"]

function listJoinedRooms(page,pageSize) {
    get("/room/list/joined?page=" + page + "&pageSize=" + pageSize, null)
        .then(data=>{
            renderRoomsList(data, "joined-rooms-rows", ["private", "name", "memberType", "owner"],(r)=>{
                const enterButton = document.createElement("button");
                enterButton.className = "btn btn-primary"
                enterButton.type = "button"
                enterButton.innerText = "Enter"
                enterButton.onclick = _ => window.location = "/room/"+r["id"]
                return enterButton
            })
        })
        .catch(err=>{
            console.log(err)
        })
}

function listRooms(page, pageSize) {
    get("/room/list?page=" + page + "&pageSize=" + pageSize)
        .then(data=>{
            console.log(data)
            renderRoomsList(data, "list-rooms-rows", ["private", "name", "owner"], r=>{
                const joinButton = document.createElement("button")
                joinButton.className = "btn btn-primary"
                joinButton.type = "button"
                joinButton.innerText = "Join"
                joinButton.onclick = _ => tryJoinRoom(r)
                return joinButton
            })
        })
        .catch(err=>{
            console.log(err)
        })
}

function renderRoomsList(data, elementId, columns, buttonCreator) {
    const rows = document.getElementById(elementId)
    rows.innerHTML = ""
    data.forEach(r=>{
        const tr = document.createElement("tr")
        columns.forEach(col=>{
            const td = document.createElement("td")
            if (col === "memberType") {
                td.innerHTML = memberTypes[r[col]]
            }else {
                td.innerHTML = r[col]
            }
            tr.appendChild(td)
        })
        tr.appendChild(buttonCreator(r))
        rows.appendChild(tr)
    })
}

function tryJoinRoom(room) {
    get("/room/"+room["id"]+"/member", null)
        .then(data=>{
            window.location = "/room/" + room["id"]
        })
        .catch(resp=>{
            if (resp.status === 403) {
                if (room["private"]) {
                    const modal = new bootstrap.Modal(document.getElementById("join-room-password-modal"), {
                        keyboard: false
                    })
                    modal.show()
                    document.getElementById("join-room-modal-button").onclick = _ => {
                        joinRoom(room)
                    }
                }else {
                    joinRoom(room)
                }
            }else {
                console.log(resp.message)
            }
        })
}

function joinRoom(room) {
    const pwd = document.getElementById("join-room-password").value
    post("/room/" + room["id"] + "/join?password=" + pwd, {})
        .then(data=>{
            window.location = "/room/" + room["id"]
        })
        .catch(resp=>{
            alert(resp.message)
        })
}

function createRoom() {
    const name = document.getElementById("create-room-name").value
    const isPrivate = document.getElementById("create-room-private").value === "on"
    if (name === "" || name.length > 16) {
        alert("invalid name")
        return
    }
    post("/room", JSON.stringify({"name": name, "private": isPrivate}))
        .then(data=>{
            console.log(data)
        })
        .catch(resp=>{
            if (resp.status === 400) {
                alert(resp.message)
            }
            console.log(resp.message)
        })
}