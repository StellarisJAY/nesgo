package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type RoomService struct {
	m        sync.Mutex
	sessions map[int64]*RoomSession
}

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
	return &RoomService{
		sessions: make(map[int64]*RoomSession),
		m:        sync.Mutex{},
	}
}

func (rs *RoomService) CreateRoom(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	var form CreateRoomForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  http.StatusBadRequest,
			Message: "bad request form",
		})
		return
	}
	_, err := room.GetRoomByNameAndOwner(form.Name, userId)
	if err == nil {
		c.JSON(200, JSONResp{
			Status:  http.StatusBadRequest,
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
	if err := room.AddMember(&room.Member{
		RoomId:     r.Id,
		UserId:     userId,
		MemberType: room.MemberTypeOwner,
	}); err != nil {
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
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	roomId := c.Param("roomId")
	inviteCode := c.Query("inviteCode")
	if roomId == "" || inviteCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request: invalid roomId or inviteCode"})
		return
	}
	id, err := strconv.ParseInt(roomId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request: invalid roomId"})
		return
	}
	r, err := room.GetRoomById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, JSONResp{
				Status:  404,
				Message: "room not found",
			})
			return
		}
		panic(err)
	}
	if r.InviteCode == inviteCode {
		err := room.AddMember(&room.Member{
			RoomId:     id,
			UserId:     userId,
			MemberType: room.MemberTypeWatcher,
		})
		if err != nil {
			panic(err)
		}
		c.JSON(200, JSONResp{
			Status:  200,
			Message: "Success",
		})
	} else {
		c.JSON(200, JSONResp{
			Status:  500,
			Message: "wrong invite code",
		})
		return
	}
}

func (rs *RoomService) HandleWebsocket(c *gin.Context) {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid room id",
		})
		return
	}
	var member *room.Member
	// check membership
	if m, ok := rs.IsRoomMember(roomId, userId); !ok {
		c.JSON(200, JSONResp{
			Status:  http.StatusForbidden,
			Message: "not a member of this room",
		})
		return
	} else {
		member = m
	}

	rs.m.Lock()
	// check if room's game session is created
	if s, ok := rs.sessions[roomId]; !ok {
		rs.m.Unlock()
		c.JSON(200, JSONResp{Status: 400, Message: "emulator is not running"})
	} else {
		rs.m.Unlock()
		// handle room websocket conn
		conn, err := websocket.Upgrade(c.Writer, c.Request, http.Header{}, 1024, 1024)
		if err != nil {
			panic(err)
		}
		s.newConnChan <- &RoomConnection{
			conn: conn,
			m:    member,
		}
	}
}

func (rs *RoomService) StartGame(c *gin.Context) {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid room id",
		})
		return
	}
	member, err := room.GetMember(roomId, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, JSONResp{
				Status:  http.StatusForbidden,
				Message: "not a member of this room",
			})
			return
		} else {
			panic(err)
		}
	}
	if member.MemberType != room.MemberTypeOwner {
		c.JSON(http.StatusForbidden, JSONResp{Status: http.StatusForbidden, Message: "not owner of this room"})
		return
	}

	// create game session and start game
	rs.m.Lock()
	if _, ok := rs.sessions[roomId]; ok {
		rs.m.Unlock()
	} else {
		// todo select game file
		session := newRoomSession(roomId, "SuperMario.nes")
		go session.ControlLoop()
		session.StartGame()
		rs.sessions[roomId] = session
		rs.m.Unlock()
	}
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "emulator is running",
	})
}

func (rs *RoomService) RoomPage(c *gin.Context) {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid room id",
		})
		return
	}
	roomDO, err := room.GetRoomById(roomId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, JSONResp{
				Status:  404,
				Message: "room not found",
			})
			return
		} else {
			panic(err)
		}
	}
	if _, ok := rs.IsRoomMember(roomId, userId); !ok {
		c.JSON(200, JSONResp{
			Status:  http.StatusForbidden,
			Message: "not a member of this room",
		})
		return
	}

	c.HTML(200, "room.html", roomDO)
}

func (rs *RoomService) IsRoomMember(roomId, userId int64) (*room.Member, bool) {
	m, err := room.GetMember(roomId, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false
		}
		panic(err)
	}
	return m, true
}

func (rs *RoomService) GetRoomInfo(c *gin.Context) {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid room id",
		})
		return
	}
	roomDO, err := room.GetRoomById(roomId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, JSONResp{
				Status:  404,
				Message: "room not found",
			})
			return
		} else {
			panic(err)
		}
	}
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data:    roomDO,
	})
}

func (rs *RoomService) ListRoomMembers(c *gin.Context) {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid room id",
		})
		return
	}
	if _, ok := rs.IsRoomMember(roomId, userId); !ok {
		c.JSON(200, JSONResp{
			Status:  403,
			Message: "not member of this room",
		})
		return
	}

	memberIds, err := room.ListRoomMembers(roomId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, JSONResp{
				Status:  200,
				Message: "ok",
			})
			return
		} else {
			panic(err)
		}
	}
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "OK",
		Data:    memberIds,
	})
}

func (rs *RoomService) GetMemberType(c *gin.Context) {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid room id",
		})
		return
	}
	m, ok := rs.IsRoomMember(roomId, userId)
	if !ok {
		c.JSON(200, JSONResp{
			Status:  403,
			Message: "not a member of this room",
		})
	} else {
		c.JSON(200, JSONResp{
			Status:  200,
			Message: "ok",
			Data:    m.MemberType,
		})
	}
}

func (rs *RoomService) GetRoomMemberSelf(c *gin.Context) {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid room id",
		})
		return
	}
	m, err := room.GetMemberFull(roomId, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, JSONResp{
				Status:  403,
				Message: "not a member of this room",
			})
			return
		}
		panic(err)
	}
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data:    m,
	})
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
