package save

import (
	"errors"
	"github.com/stellarisJAY/nesgo/web/model/db"
	"github.com/stellarisJAY/nesgo/web/util"
	"gorm.io/gorm"
	"time"
)

type Save struct {
	Id        int64     `json:"id" gorm:"column:id;primary key;AUTO_INCREMENT"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;"`
	Game      string    `json:"game" gorm:"column:game;"`
	Storage   string    `json:"storage" gorm:"column:storage;"`
	Path      string    `json:"path" gorm:"column:path;"`
	RoomId    int64     `json:"roomId" gorm:"column:room_id;index"`
}

const (
	MaxSaveCount = 10
)

var ErrMaxSaveCountReached = errors.New("reached save file limit")

func init() {
	d := db.GetDB()
	if err := d.AutoMigrate(&Save{}); err != nil {
		panic(err)
	}
}

func ListSaves(roomId int64, page, pageSize int) ([]*Save, error) {
	var saves []*Save
	if err := db.GetDB().
		Model(&Save{}).
		Where("room_id=?", roomId).
		Scopes(util.Page(page, pageSize)).
		Find(&saves).
		Error; err != nil {
		return nil, err
	}
	return saves, nil
}

func GetLastSave(roomId int64) (*Save, error) {
	var s Save
	if err := db.GetDB().
		Model(&Save{}).
		Where("room_id=?", roomId).
		Order("created_at DESC").
		First(&s).
		Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func GetSave(db *gorm.DB, id int64) (*Save, error) {
	var s Save
	if err := db.
		Model(&Save{}).
		Where("id=?", id).
		First(&s).
		Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func CountSaves(db *gorm.DB, roomId int64) (int64, error) {
	var res int64
	if err := db.
		Model(&Save{}).
		Where("room_id=?", roomId).
		Count(&res).
		Error; err != nil {
		return 0, err
	}
	return res, nil
}

func DeleteSave(db *gorm.DB, id int64) error {
	return db.
		Where("id=?", id).
		Delete(&Save{}).
		Error
}
