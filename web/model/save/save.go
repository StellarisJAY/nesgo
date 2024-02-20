package save

import (
	"github.com/stellarisJAY/nesgo/web/model/db"
	"github.com/stellarisJAY/nesgo/web/util"
	"gorm.io/gorm"
)

type Save struct {
	gorm.Model
	Game    string `json:"game" gorm:"column:game;"`
	Storage string `json:"storage" gorm:"column:storage;"`
	Path    string `json:"path" gorm:"column:path;"`
	RoomId  int64  `json:"roomId" gorm:"column:room_id;index"`
}

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
