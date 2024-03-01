package room

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/web/model"
	"github.com/stellarisJAY/nesgo/web/model/db"
	"github.com/stellarisJAY/nesgo/web/util"
	"gorm.io/gorm"
)

type Room struct {
	Id       int64  `gorm:"column:id;primary key;AUTO_INCREMENT" json:"id"`
	Host     int64  `gorm:"column:host" json:"host"`
	Name     string `gorm:"column:name" json:"name"`
	Password string `gorm:"column:password" json:"password"`
}

type Member struct {
	RoomId int64 `gorm:"column:room_id;primary key" json:"roomId"`
	UserId int64 `gorm:"column:user_id;primary key" json:"userId"`
	Role   byte  `gorm:"column:role" json:"role"`
}

type UserMember struct {
	Id        int64  `gorm:"column:id" json:"id"`
	Name      string `gorm:"column:name" json:"name"`
	AvatarURL string `gorm:"column:avatar_url" json:"avatarURL"`
	Role      byte   `gorm:"column:role" json:"role"`
}

type JoinedRoom struct {
	Id       int64  `gorm:"column:id;primary key;AUTO_INCREMENT" json:"id"`
	Host     int64  `gorm:"column:host" json:"host"`
	Name     string `gorm:"column:name" json:"name"`
	Password string `gorm:"column:password" json:"password"`
	Role     byte   `gorm:"column:role" json:"role"`
}

const (
	RoleHost byte = iota
	RoleGamer
	RoleObserver
	MaxMemberCount = 4
)

func init() {
	d := db.GetDB()
	if err := d.AutoMigrate(&Room{}, &Member{}); err != nil {
		panic(err)
	}
}

func CreateRoom(db *gorm.DB, room *Room) error {
	if err := db.Create(room).Error; err != nil {
		return err
	}
	return nil
}

func GetRoomByNameAndHost(name string, host int64) (*Room, error) {
	d := db.GetDB()
	var r Room
	if err := d.Where("name=? AND host=?", name, host).
		First(&r).
		Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func GetRoomById(id int64) (*Room, error) {
	if r, err := model.CacheGet(CacheKeyForRoom(id), func(_ string) (*Room, error) {
		var r Room
		if err := db.GetDB().
			Where("id=?", id).
			First(&r).
			Error; err != nil {
			return nil, err
		}
		return &r, nil
	}); err != nil {
		return nil, err
	} else {
		return r, nil
	}
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

func AddMember(db *gorm.DB, member *Member) error {
	return db.Create(&member).Error
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

func GetMemberFull(roomId, userId int64) (*UserMember, error) {
	var um UserMember
	if err := db.GetDB().
		Select("id,name,avatar_url,role").
		Table("users").
		Joins("inner join members on users.id = members.user_id").
		Where("user_id=? AND room_id=?", userId, roomId).
		First(&um).Error; err != nil {
		return nil, err
	}
	return &um, nil
}

func GetJoinedRooms(userId int64, page, pageSize int) ([]*JoinedRoom, error) {
	var joinedRooms []*JoinedRoom
	if err := db.GetDB().
		Select("id, host, name, role, password").
		Table("rooms").
		Joins("inner join members on rooms.id=members.room_id").
		Where("members.user_id=?", userId).
		Scopes(util.Page(page, pageSize)).
		Find(&joinedRooms).
		Error; err != nil {
		return nil, err
	}
	return joinedRooms, nil
}

func ListAllRooms(page, pageSize int) ([]*Room, error) {
	var rooms []*Room
	if err := db.GetDB().
		Model(&Room{}).
		Scopes(util.Page(page, pageSize)).
		Find(&rooms).
		Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func GetMemberCount(roomId int64) (int, error) {
	var count int64 = 0
	if err := db.GetDB().
		Model(&Member{}).
		Where("room_id=?", roomId).
		Count(&count).
		Error; err != nil {
		return 0, err
	} else {
		return int(count), nil
	}
}

func DeleteMember(roomId, memberId int64) error {
	return db.GetDB().
		Where("room_id=? AND user_id=?", roomId, memberId).
		Delete(&Member{}).
		Error
}

func UpdateMember(member *Member) error {
	return db.GetDB().Model(&Member{}).
		Where("room_id=? AND user_id=?", member.RoomId, member.UserId).
		Updates(member).
		Error
}

func DeleteRoom(db *gorm.DB, roomId int64) error {
	return db.Where("id=?", roomId).Delete(&Room{}).Error
}

func DeleteRoomMembers(db *gorm.DB, roomId int64) error {
	return db.Where("room_id=?", roomId).Delete(&Member{}).Error
}

func CacheKeyForRoom(roomId int64) string {
	return fmt.Sprintf("nesgo_room_%d", roomId)
}
