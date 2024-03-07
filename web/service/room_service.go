package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/web/config"
	"github.com/stellarisJAY/nesgo/web/model/db"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"github.com/stellarisJAY/nesgo/web/model/save"
	"github.com/stellarisJAY/nesgo/web/model/user"
	"github.com/stellarisJAY/nesgo/web/util/fs"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type RoomService struct {
	m           sync.Mutex
	rtcSessions map[int64]*WebRTCRoomSession
	fileStorage fs.FileStorage
}

type CreateRoomForm struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

type CreateRoomResp struct {
	RoomId   int64  `json:"roomId"`
	Password string `json:"password"`
}

type RoomVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Host     int64  `json:"host"`
	Password string `json:"password"`
}

type ListJoinedRoomVO struct {
	Role byte `json:"role"`
	RoomListVO
}

type RoomListVO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	HostName    string `json:"host"`
	Private     bool   `json:"private"`
	MemberCount int    `json:"memberCount"`
}

type FullRoomInfoVO struct {
	Id       int64  `json:"id"`
	Host     int64  `json:"host"`
	Password string `json:"password"`
	Private  bool   `json:"private"`
}

func NewRoomService() *RoomService {
	storage, err := fs.NewFileStorage(config.GetConfig().FileStorageType)
	if err != nil {
		panic(err)
	}
	return &RoomService{
		m:           sync.Mutex{},
		rtcSessions: make(map[int64]*WebRTCRoomSession),
		fileStorage: storage,
	}
}

func (rs *RoomService) CreateRoom(c *gin.Context) {
	userId := c.GetInt64("uid")
	var form CreateRoomForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{Status: http.StatusBadRequest, Message: "bad request form"})
		return
	}
	_, err := room.GetRoomByNameAndHost(form.Name, userId)
	if err == nil {
		c.JSON(200, JSONResp{Status: http.StatusBadRequest, Message: "room name already inuse"})
		return
	}
	password := ""
	if form.Private {
		password = generatePassword()
	}
	r := room.Room{
		Host:     userId,
		Name:     form.Name,
		Password: password,
	}

	err = db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := room.CreateRoom(tx, &r); err != nil {
			return err
		}
		if err := room.AddMember(tx, &room.Member{
			RoomId: r.Id,
			UserId: userId,
			Role:   room.RoleHost,
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
	})
}

func (rs *RoomService) ListJoinedRooms(c *gin.Context) {
	userId := c.GetInt64("uid")
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
			Role: r.Role,
		}
		joinedRoomVO.Name = r.Name
		joinedRoomVO.Private = r.Password != ""
		joinedRoomVO.Id = r.Id
		if name, ok := userNames[r.Host]; ok {
			joinedRoomVO.HostName = name
		} else {
			if u, err := user.GetUserById(r.Host); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
				panic(err)
			} else {
				joinedRoomVO.HostName = u.Name
				userNames[r.Host] = u.Name
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
			c.JSON(200, JSONResp{Status: 200, Message: "ok", Data: []*RoomListVO{}})
			return
		}
		panic(err)
	}
	roomVOs := roomDOToRoomListVO(rooms)
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data:    roomVOs,
	})
}

func (rs *RoomService) JoinRoom(c *gin.Context) {
	userId := c.GetInt64("uid")
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, JSONResp{Status: 404, Message: "room not found"})
		return
	} else if err != nil {
		panic(err)
	}

	count, err := room.GetMemberCount(id)
	if err != nil {
		panic(err)
	}

	if count == room.MaxMemberCount {
		c.JSON(200, JSONResp{Status: 400, Message: "room already full"})
		return
	}

	if !r.Private || r.Password == password {
		err := room.AddMember(db.GetDB(), &room.Member{
			RoomId: id,
			UserId: userId,
			Role:   room.RoleObserver,
		})
		if err != nil {
			panic(err)
		}
		c.JSON(200, JSONResp{Status: 200, Message: "Success"})
	} else {
		c.JSON(200, JSONResp{Status: 403, Message: "wrong password"})
		return
	}
}

func (rs *RoomService) GetRoomInfo(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	roomDO, err := room.GetRoomById(roomId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, JSONResp{Status: 404, Message: "room not found"})
		return
	} else if err != nil {
		panic(err)
	}
	host, _ := user.GetUserById(roomDO.Host)
	count, _ := room.GetMemberCount(roomDO.Id)
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data: RoomListVO{
			Id:          roomDO.Id,
			Name:        roomDO.Name,
			HostName:    host.Name,
			Private:     roomDO.Password != "",
			MemberCount: count,
		},
	})
}

func (rs *RoomService) GetRoomFullInfo(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	roomDO, err := room.GetRoomById(roomId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, JSONResp{Status: 404, Message: "room not found"})
		return
	} else if err != nil {
		panic(err)
	}
	c.JSON(200, JSONResp{Status: 200, Message: "ok", Data: FullRoomInfoVO{
		Id:       roomDO.Id,
		Host:     roomDO.Host,
		Password: roomDO.Password,
		Private:  roomDO.Private,
	}})
}

type AlterRoomRequest struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

func (rs *RoomService) AlterRoom(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	var req AlterRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, JSONResp{Status: 400, Message: "Bad request"})
		return
	}
	origin, err := room.GetRoomById(roomId)
	if err != nil {
		panic(err)
	}
	if req.Name != "" {
		origin.Name = req.Name
	}
	if !origin.Private && req.Private {
		origin.Password = generatePassword()
	}
	origin.Private = req.Private
	if err := room.UpdateRoom(origin); err != nil {
		panic(err)
	}
	c.JSON(200, JSONResp{200, "ok", FullRoomInfoVO{
		Id:       roomId,
		Host:     origin.Host,
		Password: origin.Password,
		Private:  origin.Private,
	}})
}

func (rs *RoomService) DeleteRoom(c *gin.Context) {
	roomId := c.GetInt64("roomId")

	rs.m.Lock()
	session, ok := rs.rtcSessions[roomId]
	if ok {
		delete(rs.rtcSessions, roomId)
		session.Close()
	}
	rs.m.Unlock()

	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := room.DeleteRoomMembers(tx, roomId); err != nil {
			return err
		}
		if err := room.DeleteRoom(tx, roomId); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	saves, err := save.ListSaves(roomId, 1, 10)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		goto RETURN
	}
	if err != nil {
		panic(err)
	}
	err = save.DeleteSave(db.GetDB(), roomId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		goto RETURN
	} else if err != nil {
		panic(err)
	}
	for _, s := range saves {
		_ = rs.fileStorage.Delete(s.Path)
	}
RETURN:
	c.JSON(200, JSONResp{Status: 200, Message: "ok"})
}

func (rs *RoomService) Search(c *gin.Context) {
	search := c.Query("search")
	if search == "" || len(search) >= 32 {
		c.JSON(200, JSONResp{Status: 400, Message: "invalid search text format"})
		return
	}
	rooms, err := room.SearchRoom(search)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, JSONResp{Status: 200, Message: "ok"})
		return
	} else if err != nil {
		panic(err)
	}
	roomVOs := roomDOToRoomListVO(rooms)
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data:    roomVOs,
	})
}

func roomDOToRoomListVO(rooms []*room.Room) []*RoomListVO {
	userNames := make(map[int64]string)
	roomVOs := make([]*RoomListVO, 0, len(rooms))
	for _, r := range rooms {
		vo := &RoomListVO{
			Id:      r.Id,
			Name:    r.Name,
			Private: r.Private,
		}
		if name, ok := userNames[r.Host]; ok {
			vo.HostName = name
		} else {
			u, err := user.GetUserById(r.Host)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			} else if err != nil {
				panic(err)
			}
			vo.HostName = u.Name
			userNames[r.Host] = u.Name
		}
		count, err := room.GetMemberCount(r.Id)
		if err != nil {
			panic(err)
		}
		vo.MemberCount = count
		roomVOs = append(roomVOs, vo)
	}
	return roomVOs
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
