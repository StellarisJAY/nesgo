package service

import (
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"net/http"
	"strconv"
)

func (rs *RoomService) Restart(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	rs.m.Lock()
	if session, ok := rs.rtcSessions[roomId]; !ok {
		rs.m.Unlock()
		c.JSON(200, JSONResp{
			Status:  http.StatusNotFound,
			Message: "game session not found",
		})
		return
	} else {
		rs.m.Unlock()
		if game := c.Query("game"); game == "" {
			c.JSON(200, JSONResp{
				Status:  http.StatusBadRequest,
				Message: "invalid game name",
			})
			return
		} else {
			err := session.restart(game)
			if err != nil {
				panic(err)
			}
			c.JSON(200, JSONResp{Status: http.StatusOK, Message: "success"})
		}
	}
}

func (rs *RoomService) TransferControl(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	target, err := strconv.ParseInt(c.Query("target"), 10, 64)
	controlId, _ := strconv.Atoi(c.Query("control"))
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid target user id",
		})
		return
	}
	targetMember, ok := rs.IsRoomMember(roomId, target)
	if !ok || targetMember.MemberType == room.MemberTypeWatcher {
		c.JSON(200, JSONResp{
			Status:  http.StatusForbidden,
			Message: "can not give control to watcher or not a room member",
		})
		return
	}

	rs.m.Lock()
	session, ok := rs.rtcSessions[roomId]
	rs.m.Unlock()
	if !ok {
		c.JSON(200, JSONResp{
			Status:  http.StatusNotFound,
			Message: "room session not found",
		})
		return
	}
	if err := session.TransferControl(target, controlId); err != nil {
		c.JSON(200, JSONResp{Status: 400, Message: err.Error()})
	} else {
		c.JSON(200, JSONResp{Status: 200, Message: "ok"})
	}
}
