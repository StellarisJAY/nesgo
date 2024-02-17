let joinedRooms = []

onload = ()=> {
    listJoinedRooms()
}

memberTypes = ["Owner", "Gamer", "Watcher"]

function listJoinedRooms() {
    get("/room/list/joined", null)
        .then(data=>{
            const rows = document.getElementById("joined-rooms-rows")
            rows.innerHTML = ""
            const columnOrder = ["private", "name", "memberType", "owner"]
            data.forEach(r=>{
                const tr = document.createElement("tr")
                columnOrder.forEach(col=>{
                    const td = document.createElement("td")
                    if (col === "memberType") {
                        td.innerHTML = memberTypes[r[col]]
                    }else {
                        td.innerHTML = r[col]
                    }
                    tr.appendChild(td)
                })
                const enterButton = document.createElement("button");
                enterButton.className = "btn btn-primary"
                enterButton.type = "button"
                enterButton.innerText = "Enter"
                enterButton.onclick = _ => window.location = "/room/"+r["id"]
                tr.appendChild(enterButton)
                rows.appendChild(tr)
            })
        })
        .catch(err=>{
            console.log(err)
        })
}