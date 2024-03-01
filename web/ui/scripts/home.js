let joinedRooms = []
let listRoomsPage = 1
let joinedRoomsPage = 1
const listRoomsPageSize = 10

onload = ()=> {
    listJoinedRooms(joinedRoomsPage, listRoomsPageSize)
    listRooms(listRoomsPage, listRoomsPageSize)

    document.getElementById("refresh-all-rooms-btn").onclick = _ => listJoinedRooms(joinedRoomsPage, listRoomsPageSize)
    document.getElementById("refresh-joined-rooms-btn").onclick = _=>listRooms(listRoomsPage, listRoomsPageSize)
    document.getElementById("list-rooms-next").onclick = _=>listRooms(listRoomsPage+1, listRoomsPageSize)
    document.getElementById("list-rooms-prev").onclick = _=>listRooms(listRoomsPage-1, listRoomsPageSize)
    document.getElementById("joined-rooms-next").onclick = _=>listJoinedRooms(joinedRoomsPage+1, listRoomsPageSize)
    document.getElementById("joined-rooms-prev").onclick = _=>listJoinedRooms(joinedRoomsPage-1, listRoomsPageSize)
}

function listJoinedRooms(page,pageSize) {
    get("/room/list/joined?page=" + page + "&pageSize=" + pageSize, null)
        .then(data=>{
            renderRoomsList(data, "joined-rooms-rows", ["private", "name", "role", "host"],(r)=>{
                const enterButton = document.createElement("button");
                enterButton.className = "btn btn-primary"
                enterButton.type = "button"
                enterButton.innerText = "Enter"
                enterButton.onclick = _ => window.location = "/room/"+r["id"]
                const deleteButton = document.createElement("button")
                deleteButton.type = "button"
                deleteButton.innerText = "delete"
                deleteButton.className = "btn btn-danger"
                deleteButton.onclick = _=> deleteRoom(r["id"])
                return [enterButton, deleteButton]
            })
            if (data.length === 0) {
                return null
            }
            joinedRoomsPage = page
            return data
        })
        .then(data=>{
            document.getElementById("joined-rooms-prev").hidden = joinedRoomsPage === 1
            document.getElementById("joined-rooms-next").hidden = data === null
            document.getElementById("joined-rooms-page").innerText = ""+ joinedRoomsPage
        })
        .catch(err=>{
            console.log(err)
        })
}

function listRooms(page, pageSize) {
    if (page === 0) {
        return
    }
    get("/room/list?page=" + page + "&pageSize=" + pageSize)
        .then(data=>{
            if (data.length === 0) {
                return null
            }
            listRoomsPage = page
            renderRoomsList(data, "list-rooms-rows", ["private", "name", "host"], r=>{
                const joinButton = document.createElement("button")
                joinButton.className = "btn btn-primary"
                joinButton.type = "button"
                joinButton.innerText = "Join"
                joinButton.onclick = _ => tryJoinRoom(r)
                return [joinButton]
            })
            return data
        })
        .then(data=>{
            document.getElementById("list-rooms-prev").hidden = listRoomsPage=== 1
            document.getElementById("list-rooms-next").hidden = data === null
            document.getElementById("list-rooms-page").innerText = ""+ listRoomsPage
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
            if (col === "role") {
                td.innerHTML = roles[r[col]]
            }else {
                td.innerHTML = r[col]
            }
            tr.appendChild(td)
        })
        const td = document.createElement("td")
        const buttons = buttonCreator(r)
        buttons.forEach(b=>td.appendChild(b))
        tr.appendChild(td)
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
    const isPrivate = document.getElementById("create-room-private").checked
    if (name === "" || name.length > 16) {
        alert("invalid name")
        return
    }
    post("/room", JSON.stringify({"name": name, "private": isPrivate}))
        .then(data=>{
            const modal = new bootstrap.Modal(document.getElementById("create-room-modal"))
            modal.dispose()
        })
        .then(_=>{
            listJoinedRooms(1, 10)
        })
        .catch(resp=>{
            if (resp.status === 400) {
                alert(resp.message)
            }
            console.log(resp.message)
        })
}

function deleteRoom(roomId ) {
    post("/room/" + roomId + "/delete", null)
        .then(_=>{
            listJoinedRooms(1, 10)
        })
        .catch(resp=>{
            if (resp.status !== 500) {
                alert(resp.message)
            }
            console.log(resp.message)
        })
}