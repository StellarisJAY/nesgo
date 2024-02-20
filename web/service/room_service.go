package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stellarisJAY/nesgo/web/config"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"github.com/stellarisJAY/nesgo/web/model/user"
	"github.com/stellarisJAY/nesgo/web/util/fs"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type RoomService struct {
	m           sync.Mutex
	rtcSessions map[int64]*RTCRoomSession
	fileStorage fs.FileStorage
}

type CreateRoomForm struct {
	Name    string `json:"name" binding:"required"`
	Private bool   `json:"private" binding:"required"`
}

type CreateRoomResp struct {
	RoomId   int64  `json:"roomId"`
	Password string `json:"password"`
}

type RoomVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Owner    int64  `json:"owner"`
	Password string `json:"password"`
}

type ListJoinedRoomVO struct {
	MemberType byte `json:"memberType"`
	RoomListVO
}

type RoomListVO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	OwnerName   string `json:"owner"`
	Private     bool   `json:"private"`
	MemberCount int    `json:"memberCount"`
}

type RoomMemberVO struct {
	Id         int64 `json:"id"`
	MemberType byte  `json:"memberType"`
}

func NewRoomService() *RoomService {
	storage, err := fs.NewFileStorage(config.GetConfig().FileStorageType)
	if err != nil {
		panic(err)
	}
	return &RoomService{
		m:           sync.Mutex{},
		rtcSessions: make(map[int64]*RTCRoomSession),
		fileStorage: storage,
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
		Owner:    userId,
		Name:     form.Name,
		Password: generatePassword(),
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
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data: CreateRoomResp{
			RoomId:   r.Id,
			Password: r.Password,
		},
	})
}

func (rs *RoomService) ListJoinedRooms(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	page, pageSize, err := getPageQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid page query params"})
		return
	}
	rooms, err := room.GetJoinedRooms(userId, page, pageSize)
	if err != nil {
		panic(err)
	}
	userNames := make(map[int64]string)
	joinedRoomVOs := make([]*ListJoinedRoomVO, 0, len(rooms))
	for _, r := range rooms {
		joinedRoomVO := &ListJoinedRoomVO{
			MemberType: r.MemberType,
		}
		joinedRoomVO.Name = r.Name
		log.Println(r.Password == "")
		joinedRoomVO.Private = r.Password != ""
		joinedRoomVO.Id = r.Id
		if name, ok := userNames[r.Owner]; ok {
			joinedRoomVO.OwnerName = name
		} else {
			if u, err := user.GetUserById(r.Owner); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
				panic(err)
			} else {
				joinedRoomVO.OwnerName = u.Name
				userNames[r.Owner] = u.Name
			}
		}
		count, err := room.GetMemberCount(r.Id)
		if err != nil {
			panic(err)
		}
		joinedRoomVO.MemberCount = count
		joinedRoomVOs = append(joinedRoomVOs, joinedRoomVO)
	}
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data:    joinedRoomVOs,
	})
}

func (rs *RoomService) ListAllRooms(c *gin.Context) {
	page, pageSize, err := getPageQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid page query params"})
		return
	}
	rooms, err := room.ListAllRooms(page, pageSize)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, JSONResp{
				Status:  200,
				Message: "ok",
				Data:    []*RoomListVO{},
			})
			return
		}
		panic(err)
	}
	userNames := make(map[int64]string)
	roomVOs := make([]*RoomListVO, 0, len(rooms))
	for _, r := range rooms {
		vo := &RoomListVO{
			Id:      r.Id,
			Name:    r.Name,
			Private: r.Password != "",
		}
		if name, ok := userNames[r.Owner]; ok {
			vo.OwnerName = name
		} else {
			u, err := user.GetUserById(r.Owner)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			} else if err != nil {
				panic(err)
			}
			vo.OwnerName = u.Name
			userNames[r.Owner] = u.Name
		}
		count, err := room.GetMemberCount(r.Id)
		if err != nil {
			panic(err)
		}
		vo.MemberCount = count
		roomVOs = append(roomVOs, vo)
	}
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data:    roomVOs,
	})
}

func (rs *RoomService) JoinRoom(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	roomId := c.Param("roomId")
	password := c.Query("password")
	if roomId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request: invalid roomId"})
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
	if r.Password == "" || r.Password == password {
		err := room.AddMember(&room.Member{
			RoomId:     id,
			UserId:     userId,
			MemberType: room.MemberTypeWatcher,
		})
		if err != nil {
			panic(err)
		}
		log.Println("inserted new member")
		c.JSON(200, JSONResp{
			Status:  200,
			Message: "Success",
		})
	} else {
		c.JSON(200, JSONResp{
			Status:  500,
			Message: "wrong password",
		})
		return
	}
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

func (rs *RoomService) ConnectRTCRoomSession(c *gin.Context) {
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
	var session *RTCRoomSession
	// check if room's game session is created
	if s, ok := rs.rtcSessions[roomId]; !ok {
		// Only owner can create session
		if member.MemberType != room.MemberTypeOwner {
			rs.m.Unlock()
			c.JSON(200, JSONResp{
				Status:  http.StatusForbidden,
				Message: "only owner can start game session",
			})
			return
		}

		game := c.Query("game")
		if game == "" {
			rs.m.Unlock()
			c.JSON(200, JSONResp{
				Status:  http.StatusBadRequest,
				Message: "invalid game name",
			})
			return
		}
		newSession, err := NewRTCRoomSession(game)
		if err != nil {
			panic(err)
		}
		rs.rtcSessions[roomId] = newSession
		ctx, cancelFunc := context.WithCancel(context.Background())
		newSession.cancel = cancelFunc
		go newSession.ControlLoop(ctx)
		session = newSession
	} else {
		session = s
	}
	rs.m.Unlock()
	// handle room websocket conn
	conn, err := websocket.Upgrade(c.Writer, c.Request, http.Header{}, 1024, 1024)
	if err != nil {
		panic(err)
	}
	session.signalChan <- Signal{
		Type: SignalNewConnection,
		Data: &WebsocketConn{
			Member: member,
			Conn:   conn,
		},
	}
}

func (rs *RoomService) Restart(c *gin.Context) {
	roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{
			Status:  400,
			Message: "invalid room id",
		})
		return
	}
	// check membership
	if m, ok := rs.IsRoomMember(roomId, userId); !ok {
		c.JSON(200, JSONResp{
			Status:  http.StatusForbidden,
			Message: "not a member of this room",
		})
		return
	} else if m.MemberType != room.MemberTypeOwner {
		c.JSON(200, JSONResp{
			Status:  http.StatusForbidden,
			Message: "only owner can restart emulator",
		})
		return
	}

	rs.m.Lock()
	if session, ok := rs.rtcSessions[roomId]; !ok {
		rs.m.Unlock()
		c.JSON(200, JSONResp{
			Status:  http.StatusNotFound,
			Message: "game session not found",
		})
		return
	} else {
		rs.m.Unlock()
		if game := c.Query("game"); game == "" {
			c.JSON(200, JSONResp{
				Status:  http.StatusBadRequest,
				Message: "invalid game name",
			})
			return
		} else {
			err := session.restart(game)
			if err != nil {
				panic(err)
			}
			c.JSON(200, JSONResp{Status: http.StatusOK, Message: "success"})
		}
	}
}

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generatePassword() string {
	codeLen := 4
	sb := strings.Builder{}
	for i := 0; i < codeLen; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
