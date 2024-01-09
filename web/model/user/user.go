package user

import (
	"github.com/stellarisJAY/nesgo/web/model/db"
)

type User struct {
	Id        int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;" json:"id"`
	Name      string `gorm:"column:name;unique;'" json:"name"`
	Password  string `gorm:"column:password" json:"password"`
	AvatarURL string `gorm:"column:avatar_url" json:"avatarURL"`
}

func init() {
	d := db.GetDB()
	if err := d.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}

func GetUserById(id int64) (*User, error) {
	d := db.GetDB()
	user := User{}
	err := d.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByName(name string) (*User, error) {
	d := db.GetDB()
	user := User{}
	if err := d.Where("name=?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *User) error {
	d := db.GetDB()
	return d.Create(user).Error
}
