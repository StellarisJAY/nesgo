package service

import (
	"errors"
	"github.com/gin-gonic/gin"
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

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generatePassword() string {
	codeLen := 4
	sb := strings.Builder{}
	for i := 0; i < codeLen; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
