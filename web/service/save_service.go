package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/web/model/db"
	"github.com/stellarisJAY/nesgo/web/model/save"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"
)

func (rs *RoomService) QuickSave(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	rs.m.Lock()
	session, ok := rs.rtcSessions[roomId]
	rs.m.Unlock()
	if !ok {
		c.JSON(200, JSONResp{Status: 404, Message: "room session not found"})
		return
	}
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
		if count, err := save.CountSaves(d, roomId); err != nil {
			return err
		} else if count == save.MaxSaveCount {
			return save.ErrMaxSaveCountReached
		}
		if err := d.Create(&s).Error; err != nil {
			return err
		}
		if err := rs.fileStorage.Store(path, data); err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, save.ErrMaxSaveCountReached) {
		c.JSON(200, JSONResp{Status: 400, Message: err.Error()})
		return
	}
	if err != nil {
		c.JSON(200, JSONResp{Status: 500, Message: err.Error()})
		return
	}
	c.JSON(200, JSONResp{Status: 200, Message: "OK"})
}

func (rs *RoomService) QuickLoad(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	saveId, _ := strconv.ParseInt(c.Param("saveId"), 10, 64)
	s, err := save.GetSave(db.GetDB(), saveId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, JSONResp{Status: 404, Message: "No save found"})
		return
	} else if err != nil {
		panic(err)
	}
	data, err := rs.fileStorage.Load(s.Path)
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
	if err := session.LoadSavedGame(s.Game, data); err != nil {
		c.JSON(200, JSONResp{Status: 500, Message: err.Error()})
		return
	}
	c.JSON(200, JSONResp{Status: 200, Message: "ok"})
}

type SaveVO struct {
	Game      string `json:"game"`
	CreatedAt string `json:"createdAt"`
	Id        int64  `json:"id"`
}

func (rs *RoomService) ListSaves(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	saves, err := save.ListSaves(roomId, 1, 10)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, JSONResp{Status: 200, Message: "ok", Data: nil})
		return
	} else if err != nil {
		panic(err)
	}

	slices.SortFunc(saves, func(a, b *save.Save) int {
		return int(b.CreatedAt.Sub(a.CreatedAt).Seconds())
	})

	vos := make([]*SaveVO, 0, len(saves))
	for _, s := range saves {
		vos = append(vos, &SaveVO{s.Game, s.CreatedAt.Format(time.DateTime), s.Id})
	}
	c.JSON(200, JSONResp{200, "ok", vos})
}

func (rs *RoomService) DeleteSave(c *gin.Context) {
	saveId, _ := strconv.ParseInt(c.Param("saveId"), 10, 64)
	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		s, err := save.GetSave(tx, saveId)
		if err != nil {
			return err
		}
		if err := save.DeleteSave(tx, saveId); err != nil {
			return err
		}
		err = rs.fileStorage.Delete(s.Path)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		} else {
			return err
		}
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, JSONResp{Status: 404, Message: "save record not found"})
		return
	} else if err != nil {
		panic(err)
	}
	c.JSON(200, JSONResp{Status: 200, Message: "ok"})
}

func getStoragePath(roomId int64, game string, timestamp time.Time) string {
	return filepath.Join(fmt.Sprintf("room_%d", roomId), game, strconv.FormatInt(timestamp.UnixMilli(), 16)+".save")
}
