package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/web/middleware"
	"github.com/stellarisJAY/nesgo/web/model/db"
	"github.com/stellarisJAY/nesgo/web/model/user"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserService struct{}

type RegisterForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Avatar   string `json:"avatar"`
}

type LoginForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) Register(c *gin.Context) {
	var form RegisterForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(400, JSONResp{
			Status:  400,
			Message: "Bad request form",
		})
		return
	}
	_, err := user.GetUserByName(form.Name)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, JSONResp{
			Status:  200,
			Message: "User Name already inuse",
		})
		return
	}
	if err := user.CreateUser(&user.User{
		Name:      form.Name,
		Password:  form.Password,
		AvatarURL: form.Avatar,
	}); err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{})
}

func (u *UserService) Login(c *gin.Context) {
	var loginForm LoginForm
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(400, JSONResp{
			Status:  400,
			Message: "Bad request form",
		})
		return
	}
	usr, err := user.GetUserByName(loginForm.Name)
	if err == nil {
		if loginForm.Password != usr.Password {
			c.JSON(401, JSONResp{
				Status:  401,
				Message: "Wrong password or username",
			})
			return
		}
		token, _ := middleware.GenerateToken(usr)
		client := db.GetRedis()
		data, _ := json.Marshal(usr)
		if _, err := client.Set("user_"+strconv.FormatInt(usr.Id, 10), data, time.Hour*24).Result(); err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{
			"token": token,
		})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, JSONResp{
			Status:  404,
			Message: "User not found",
		})
		return
	}
	panic(err)
}
