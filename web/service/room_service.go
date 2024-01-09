package service

import (
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"math/rand"
	"strconv"
	"strings"
)

type RoomService struct{}

type CreateRoomForm struct {
	Name string `json:"name" binding:"required"`
}

type CreateRoomResp struct {
	JSONResp
	RoomId     int64  `json:"roomId"`
	InviteCode string `json:"inviteCode"`
}

type RoomVO struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Owner      int64  `json:"owner"`
	InviteCode string `json:"inviteCode"`
}

type RoomMemberVO struct {
	Id         int64 `json:"id"`
	MemberType byte  `json:"memberType"`
}

func NewRoomService() *RoomService {
	return &RoomService{}
}

func (rs *RoomService) CreateRoom(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	var form CreateRoomForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(401, JSONResp{
			Status:  401,
			Message: "bad request form",
		})
		return
	}
	_, err := room.GetRoomByNameAndOwner(form.Name, userId)
	if err == nil {
		c.JSON(200, JSONResp{
			Status:  401,
			Message: "room name already inuse",
		})
		return
	}
	r := room.Room{
		Owner:      userId,
		Name:       form.Name,
		InviteCode: generateInviteCode(),
	}
	if err := room.CreateRoom(&r); err != nil {
		panic(err)
	}
	c.JSON(200, CreateRoomResp{
		RoomId:     r.Id,
		InviteCode: r.InviteCode,
	})
}

func (rs *RoomService) ListOwningRooms(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	rooms, err := room.GetRoomsByOwnerId(userId)
	if err != nil {
		panic(err)
	}
	roomsVO := make([]RoomVO, 0, len(rooms))
	for _, r := range rooms {
		roomsVO = append(roomsVO, RoomVO{
			Id:         r.Id,
			Name:       r.Name,
			Owner:      r.Owner,
			InviteCode: r.InviteCode,
		})
	}
	c.JSON(200, JSONResp{
		Data: roomsVO,
	})
}

func (rs *RoomService) JoinRoom(c *gin.Context) {

}

func (rs *RoomService) ListRoomMembers(c *gin.Context) {
	
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateInviteCode() string {
	codeLen := 8
	sb := strings.Builder{}
	for i := 0; i < codeLen; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
