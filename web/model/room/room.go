package room

import (
	"github.com/stellarisJAY/nesgo/web/model/db"
	"log"
)

type Room struct {
	Id         int64  `gorm:"column:id;primary key;AUTO_INCREMENT" json:"id"`
	Owner      int64  `gorm:"column:owner" json:"owner"`
	Name       string `gorm:"column:name" json:"name"`
	InviteCode string `gorm:"column:invite_code" json:"inviteCode"`
}

type Member struct {
	RoomId     int64 `gorm:"column:room_id;primary key" json:"roomId"`
	UserId     int64 `gorm:"column:user_id;primary key" json:"userId"`
	MemberType byte  `gorm:"column:member_type" json:"memberType"`
}

const (
	MemberTypeOwner byte = iota
	MemberTypeGamer
	MemberTypeWatcher
)

func init() {
	d := db.GetDB()
	if err := d.AutoMigrate(&Room{}, &Member{}); err != nil {
		panic(err)
	}
}

func CreateRoom(room *Room) error {
	d := db.GetDB()
	if err := d.Create(room).Error; err != nil {
		return err
	}
	log.Println(room.Id)
	return nil
}

func GetRoomByNameAndOwner(name string, owner int64) (*Room, error) {
	d := db.GetDB()
	var r Room
	if err := d.Where("name=? AND owner=?", name, owner).
		First(&r).
		Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func GetRoomsByOwnerId(owner int64) ([]*Room, error) {
	d := db.GetDB()
	var rooms []*Room
	if err := d.Where("owner=?", owner).
		Find(&rooms).
		Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func GetRoomById(id int64) (*Room, error) {
	d := db.GetDB()
	var r Room
	err := d.Where("id=?", id).First(&r).Error
	return &r, err
}

func ListRoomMembers(roomId int64) ([]*Member, error) {
	d := db.GetDB()
	var members []*Member
	err := d.Where("room_id=?", roomId).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func AddMember(member *Member) error {
	return db.GetDB().Create(&member).Error
}

func GetMember(roomId, userId int64) (*Member, error) {
	var member Member
	if err := db.GetDB().
		Where("room_id=? AND user_id=?", roomId, userId).
		First(&member).
		Error; err != nil {
		return nil, err
	}
	return &member, nil
}
