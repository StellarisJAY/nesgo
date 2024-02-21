package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/web/model/db"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"github.com/stellarisJAY/nesgo/web/model/save"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func (rs *RoomService) QuickSave(c *gin.Context) {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid room id",
		})
		return
	}
	// check membership
	if m, ok := rs.IsRoomMember(roomId, userId); !ok {
		c.JSON(200, JSONResp{
			Status:  http.StatusForbidden,
			Message: "not a member of this room",
		})
		return
	} else if m.MemberType != room.MemberTypeOwner {
		c.JSON(200, JSONResp{
			Status:  http.StatusForbidden,
			Message: "not a owner of this room",
		})
		return
	}

	rs.m.Lock()
	if session, ok := rs.rtcSessions[roomId]; ok {
		rs.m.Unlock()
		path := getStoragePath(roomId, "", time.Now())
		data, err := session.Save()
		if err != nil {
			panic(err)
		}
		s := &save.Save{
			Game:    session.game,
			Storage: rs.fileStorage.Type(),
			Path:    path,
			RoomId:  roomId,
		}
		// 避免脏数据
		err = db.GetDB().Transaction(func(d *gorm.DB) error {
			if err := d.Create(&s).Error; err != nil {
				return err
			}
			if err := rs.fileStorage.Store(path, data); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			c.JSON(200, JSONResp{
				Status:  500,
				Message: err.Error(),
			})
			return
		}
		c.JSON(200, JSONResp{Status: 200, Message: "OK"})
	} else {
		rs.m.Unlock()
		c.JSON(200, JSONResp{
			Status:  404,
			Message: "room session not found",
		})
	}
}

func (rs *RoomService) QuickLoad(c *gin.Context) {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid room id",
		})
		return
	}
	// check membership
	if m, ok := rs.IsRoomMember(roomId, userId); !ok {
		c.JSON(200, JSONResp{
			Status:  http.StatusForbidden,
			Message: "not a member of this room",
		})
		return
	} else if m.MemberType != room.MemberTypeOwner {
		c.JSON(200, JSONResp{
			Status:  http.StatusForbidden,
			Message: "not a owner of this room",
		})
		return
	}

	lastSave, err := save.GetLastSave(roomId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, JSONResp{Status: 404, Message: "No save found"})
		return
	} else if err != nil {
		panic(err)
	}
	data, err := rs.fileStorage.Load(lastSave.Path)
	if err != nil {
		c.JSON(200, JSONResp{Status: 500, Message: err.Error()})
		return
	}

	rs.m.Lock()
	session, ok := rs.rtcSessions[roomId]
	rs.m.Unlock()
	if !ok {
		c.JSON(200, JSONResp{Status: 404, Message: "Room session not found"})
		return
	}
	if err := session.LoadSavedGame(data); err != nil {
		c.JSON(200, JSONResp{Status: 500, Message: err.Error()})
		return
	}
	c.JSON(200, JSONResp{Status: 200, Message: "ok"})
}

func (rs *RoomService) Save(c *gin.Context) {

}

func (rs *RoomService) ListSaves(c *gin.Context) {

}

func getStoragePath(roomId int64, game string, timestamp time.Time) string {
	return filepath.Join(fmt.Sprintf("room_%d", roomId), game, strconv.FormatInt(timestamp.UnixMilli(), 16)+".save")
}
