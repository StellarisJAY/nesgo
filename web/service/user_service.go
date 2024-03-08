package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/web/middleware"
	"github.com/stellarisJAY/nesgo/web/model/user"
	"gorm.io/gorm"
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
			Status:  400,
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
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
	})
}

func (u *UserService) Login(c *gin.Context) {
	var loginForm LoginForm
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(200, JSONResp{Status: 400, Message: "Bad request form"})
		return
	}
	usr, err := user.GetUserByName(loginForm.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, JSONResp{Status: 404, Message: "User not found"})
		return
	} else if err != nil {
		panic(err)
	}
	if loginForm.Password != usr.Password {
		c.JSON(200, JSONResp{Status: 401, Message: "Wrong password or username"})
		return
	}
	auth := middleware.Authorization{
		UserId:    usr.Id,
		Ip:        c.RemoteIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}
	token, err := middleware.StoreAuthorization(auth)
	if err != nil {
		panic(err)
	}
	c.JSON(200, JSONResp{Status: 200, Message: "ok", Data: token})
}
