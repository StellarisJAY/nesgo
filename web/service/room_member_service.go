package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"github.com/stellarisJAY/nesgo/web/model/user"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type RoomMemberVO struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	MemberType byte   `json:"memberType"`
}

type AlterRoomMemberRequest struct {
	RoomId        int64 `json:"roomId"`
	MemberId      int64 `json:"memberId" binding:"required"`
	MemberType    byte  `json:"memberType"`
	Kick          bool  `json:"kick"`
	SetController int   `json:"setController"`
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

func (rs *RoomService) ListRoomMembers(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	members, _ := room.ListRoomMembers(roomId)
	result := make([]*RoomMemberVO, 0, len(members))
	for _, member := range members {
		u, err := user.GetUserById(member.UserId)
		if err != nil {
			continue
		}
		result = append(result, &RoomMemberVO{
			Id:         member.UserId,
			Name:       u.Name,
			MemberType: member.MemberType,
		})
	}
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "OK",
		Data:    result,
	})
}

func (rs *RoomService) GetMemberType(c *gin.Context) {
	m, _ := c.Get("optMember")
	member := m.(*room.Member)
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data:    member.MemberType,
	})
}

func (rs *RoomService) GetRoomMemberSelf(c *gin.Context) {
	m, _ := c.Get("optMember")
	member := m.(*room.Member)
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data:    member,
	})
}

func (rs *RoomService) KickMember(c *gin.Context) {
	roomId := c.GetInt64("roomId")
	userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	operator, _ := c.Get("optMember")
	operator = operator.(*room.Member)
	var req AlterRoomMemberRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, JSONResp{Status: 400, Message: err.Error()})
		return
	}
	if userId == req.MemberId {
		c.JSON(200, JSONResp{Status: 400, Message: "can not kick yourself"})
		return
	}
	_, ok := rs.IsRoomMember(roomId, req.MemberId)
	if !ok {
		c.JSON(200, JSONResp{Status: 404, Message: "room member not found"})
		return
	}
	err = room.DeleteMember(roomId, req.MemberId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(200, JSONResp{Status: 404, Message: "room member not found"})
		return
	} else if err != nil {
		panic(err)
	}

	rs.m.Lock()
	session, ok := rs.rtcSessions[roomId]
	rs.m.Unlock()
	if ok {
		session.wsBroadcast(MessageRoomMemberChange, RoomMemberChangeNotification{
			MemberId: req.MemberId,
			Kick:     true,
		})
	}
	c.JSON(200, JSONResp{Status: 200, Message: "ok"})
}

// RoomMemberVerifier this middleware checks operator's identity.
// Passes {"roomId":roomId, "optMember":*Member} to Next
func (rs *RoomService) RoomMemberVerifier(accessibleMemberType []byte) func(*gin.Context) {
	return func(c *gin.Context) {
		roomId, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
		userId, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, JSONResp{
				Status:  400,
				Message: "invalid room id",
			})
			return
		}
		m, ok := rs.IsRoomMember(roomId, userId)
		if !ok {
			c.AbortWithStatusJSON(200, JSONResp{
				Status:  403,
				Message: "not a member of this room",
			})
			return
		}
		accessible := false
		for _, mType := range accessibleMemberType {
			if mType == m.MemberType {
				accessible = true
				break
			}
		}
		if !accessible {
			c.AbortWithStatusJSON(200, JSONResp{
				Status:  403,
				Message: "not accessible",
			})
			return
		}
		c.Set("optMember", m)
		c.Set("roomId", roomId)
		c.Next()
	}
}

func (rs *RoomService) OwnerAccessible() func(*gin.Context) {
	return rs.RoomMemberVerifier([]byte{room.MemberTypeOwner})
}

func (rs *RoomService) MemberAccessible() func(*gin.Context) {
	return rs.RoomMemberVerifier([]byte{room.MemberTypeOwner, room.MemberTypeWatcher, room.MemberTypeGamer})
}