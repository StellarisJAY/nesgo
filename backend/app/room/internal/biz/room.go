package biz

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

type Room struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Host         int64     `json:"host"`
	HostName     string    `json:"hostName"`
	Private      bool      `json:"private"`
	Password     string    `json:"password"`
	PasswordHash string    `json:"passwordHash"`
	MemberCount  int       `json:"memberCount"`
	MemberLimit  int       `json:"memberLimit"`
	CreateTime   time.Time `json:"createTime"`
}

type JoinedRoom struct {
	Room
	UserId int64       `json:"userId"`
	Role   v1.RoomRole `json:"role"`
}

type RoomMember struct {
	Id       int64       `json:"id"`
	UserId   int64       `json:"userId"`
	Role     v1.RoomRole `json:"role"`
	RoomId   int64       `json:"roomId"`
	JoinedAt time.Time   `json:"joinedTime"`
}

type RoomSession struct {
	RoomId   int64  `json:"roomId"`
	Endpoint string `json:"endpoint"`
}

type RoomRepo interface {
	CreateRoom(ctx context.Context, room *Room) error
	GetRoom(ctx context.Context, id int64) (*Room, error)
	ListRooms(ctx context.Context, page int, pageSize int) ([]*Room, int, error)
	ListJoinedRooms(ctx context.Context, userId int64, page int, pageSize int) ([]*JoinedRoom, int, error)
	GetRoomMember(ctx context.Context, roomId int64, userId int64) (*RoomMember, error)
	AddRoomMember(ctx context.Context, member *RoomMember, room *Room) error
	GetOrCreateRoomSession(ctx context.Context, roomId int64, game string) (*RoomSession, bool, error)
	GetRoomSession(ctx context.Context, roomId int64) (*RoomSession, error)
	RemoveRoomSession(ctx context.Context, roomId int64) error
	GetOwnedRoom(ctx context.Context, name string, host int64) (*Room, error)
	CountMember(ctx context.Context, roomId int64) (int64, error)
	ListMembers(ctx context.Context, roomId int64) ([]*RoomMember, error)
	UpdateRoom(ctx context.Context, room *Room) error
	DeleteRoom(ctx context.Context, roomId int64) error
	UpdateMember(ctx context.Context, member *RoomMember) error
	DeleteMember(ctx context.Context, roomId, userId int64) error
}

var ErrMemberLimitReached = errors.New("member limit reached")

type RoomUseCase struct {
	rr     RoomRepo
	logger *log.Helper
}

func NewRoomUseCase(rr RoomRepo, logger log.Logger) *RoomUseCase {
	return &RoomUseCase{
		rr:     rr,
		logger: log.NewHelper(log.With(logger, "module", "biz/room")),
	}
}

func (uc *RoomUseCase) CreateRoom(ctx context.Context, room *Room) error {
	r, _ := uc.rr.GetOwnedRoom(ctx, room.Name, room.Host)
	if r != nil {
		return v1.ErrorCreateRoomFailed("room already exists")
	}
	if room.Private {
		room.Password = generatePassword(4)
	}
	// TODO 根据用户级别设置房间人数上限
	room.MemberLimit = 4
	err := uc.rr.CreateRoom(ctx, room)
	if err != nil {
		return v1.ErrorCreateRoomFailed("database error: %v", err)
	}
	return nil
}

func (uc *RoomUseCase) GetRoom(ctx context.Context, id int64) (*Room, error) {
	room, err := uc.rr.GetRoom(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, v1.ErrorRoomNotFound("room not found")
	}
	if err != nil {
		return nil, v1.ErrorGetRoomFailed("database error: %v", err)
	}
	memberCount, err := uc.rr.CountMember(ctx, room.Id)
	if err != nil {
		return nil, v1.ErrorGetRoomFailed("database error: %v", err)
	}
	room.MemberCount = int(memberCount)
	return room, nil
}

func (uc *RoomUseCase) ListRooms(ctx context.Context, page int, pageSize int) ([]*Room, int, error) {
	rooms, total, err := uc.rr.ListRooms(ctx, page, pageSize)
	if err != nil {
		return nil, 0, v1.ErrorGetRoomFailed("database error: %v", err)
	}
	for _, room := range rooms {
		memberCount, err := uc.rr.CountMember(ctx, room.Id)
		if err != nil {
			return nil, 0, v1.ErrorGetRoomFailed("database error: %v", err)
		}
		room.MemberCount = int(memberCount)
	}
	return rooms, total, nil
}

func (uc *RoomUseCase) ListJoinedRooms(ctx context.Context, userId int64, page int, pageSize int) ([]*JoinedRoom, int, error) {
	rooms, total, err := uc.rr.ListJoinedRooms(ctx, userId, page, pageSize)
	if err != nil {
		return nil, 0, v1.ErrorGetRoomFailed("database error: %v", err)
	}
	for _, room := range rooms {
		memberCount, err := uc.rr.CountMember(ctx, room.Id)
		if err != nil {
			return nil, 0, v1.ErrorGetRoomFailed("database error: %v", err)
		}
		room.MemberCount = int(memberCount)
	}
	return rooms, total, nil
}

func (uc *RoomUseCase) JoinRoom(ctx context.Context, userId int64, roomId int64, password string) error {
	member, _ := uc.rr.GetRoomMember(ctx, roomId, userId)
	if member == nil {
		room, err := uc.rr.GetRoom(ctx, roomId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return v1.ErrorRoomNotFound("room not found")
		}
		if err != nil {
			return v1.ErrorGetRoomFailed("database error: %v", err)
		}
		if room.Private {
			sum := hex.EncodeToString(md5.New().Sum([]byte(password)))
			if sum != room.PasswordHash {
				return v1.ErrorRoomNotAccessible("wrong password")
			}
		}
		err = uc.rr.AddRoomMember(ctx, &RoomMember{
			UserId:   userId,
			Role:     v1.RoomRole_Observer,
			RoomId:   roomId,
			JoinedAt: time.Now(),
		}, room)
		if errors.Is(err, ErrMemberLimitReached) {
			return v1.ErrorRoomNotAccessible("room is full")
		}
		if err != nil {
			return v1.ErrorCreateRoomFailed("database error: %v", err)
		}
	}
	return nil
}

func (uc *RoomUseCase) GetCreateRoomSession(ctx context.Context, roomId, userId int64, game string) (*RoomSession, error) {
	member, _ := uc.rr.GetRoomMember(ctx, roomId, userId)
	if member == nil {
		return nil, v1.ErrorRoomNotAccessible("not a member of this room")
	}
	session, err := uc.rr.GetRoomSession(ctx, roomId)
	if session != nil {
		return session, nil
	}
	if err != nil {
		return nil, v1.ErrorGetRoomFailed("database error: %v", err)
	}
	session, _, err = uc.rr.GetOrCreateRoomSession(ctx, roomId, game)
	if err != nil {
		return nil, v1.ErrorGetRoomFailed("create session error: %v", err)
	}
	if session == nil {
		return nil, v1.ErrorCreateRoomSessionFailed("no available service to create session")
	}
	return session, nil
}

func (uc *RoomUseCase) GetRoomSession(ctx context.Context, roomId int64) (*RoomSession, error) {
	session, err := uc.rr.GetRoomSession(ctx, roomId)
	if err != nil {
		return nil, v1.ErrorGetRoomFailed("database error: %v", err)
	}
	if session == nil {
		return nil, v1.ErrorGetRoomFailed("session not available")
	}
	return session, nil
}

func (uc *RoomUseCase) RemoveRoomSession(ctx context.Context, roomId int64) error {
	err := uc.rr.RemoveRoomSession(ctx, roomId)
	if err != nil {
		return v1.ErrorGetRoomFailed("remove room session failed: %v", err)
	}
	return nil
}

func (uc *RoomUseCase) ListRoomMembers(ctx context.Context, roomId int64) ([]*RoomMember, error) {
	members, err := uc.rr.ListMembers(ctx, roomId)
	if err != nil {
		return nil, v1.ErrorGetRoomFailed("database error: %v", err)
	}
	return members, nil
}

func (uc *RoomUseCase) GetRoomMember(ctx context.Context, roomId, userId int64) (*RoomMember, error) {
	member, err := uc.rr.GetRoomMember(ctx, roomId, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, v1.ErrorGetRoomMemberFailed("member not found")
	}
	if err != nil {
		return nil, v1.ErrorGetRoomMemberFailed("database error: %v", err)
	}
	return member, nil
}

func (uc *RoomUseCase) DeleteRoom(ctx context.Context, roomId, userId int64) error {
	room, err := uc.rr.GetRoom(ctx, roomId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return v1.ErrorRoomNotFound("room not found")
	}
	if err != nil {
		return v1.ErrorGetRoomFailed("database error: %v", err)
	}
	if room.Host != userId {
		return v1.ErrorDeleteRoomFailed("you are not the owner of this room")
	}
	return uc.rr.DeleteRoom(ctx, roomId)
}

func (uc *RoomUseCase) UpdateRoom(ctx context.Context, room *Room, userId int64) error {
	r, err := uc.rr.GetRoom(ctx, room.Id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return v1.ErrorRoomNotFound("room not found")
	}
	if err != nil {
		return v1.ErrorUpdateRoomFailed("database error: %v", err)
	}
	if r.Host != userId {
		return v1.ErrorUpdateRoomFailed("you are not the owner of this room")
	}
	if r.Private == room.Private {
		room.Password = r.Password
	} else if !r.Private && room.Private {
		room.Password = generatePassword(4)
		room.PasswordHash = hex.EncodeToString(md5.New().Sum([]byte(room.Password)))
	} else if !room.Private {
		room.Password = ""
		room.PasswordHash = ""
	}
	err = uc.rr.UpdateRoom(ctx, room)
	if err != nil {
		return v1.ErrorUpdateRoomFailed("database error: %v", err)
	}
	return nil
}

func generatePassword(length int) string {
	dict := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456"
	builder := strings.Builder{}
	for i := 0; i < length; i++ {
		builder.WriteByte(dict[rand.Intn(len(dict))])
	}
	return builder.String()
}

func (uc *RoomUseCase) UpdateMember(ctx context.Context, member *RoomMember) error {
	old, err := uc.rr.GetRoomMember(ctx, member.RoomId, member.UserId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return v1.ErrorUpdateMemberFailed("room member not found")
	}
	if err != nil {
		return v1.ErrorUpdateMemberFailed("database error: %v", err)
	}
	if old.Role == member.Role {
		return nil
	}
	if old.Role == v1.RoomRole_Host {
		return nil
	}
	return uc.rr.UpdateMember(ctx, member)
}

func (uc *RoomUseCase) DeleteMember(ctx context.Context, roomId, userId int64) error {
	return uc.rr.DeleteMember(ctx, roomId, userId)
}
