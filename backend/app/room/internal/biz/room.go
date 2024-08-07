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
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Host         int64  `json:"host"`
	HostName     string `json:"hostName"`
	Private      bool   `json:"private"`
	Password     string `json:"password"`
	PasswordHash string `json:"passwordHash"`
	MemberCount  int    `json:"memberCount"`
}

type JoinedRoom struct {
	Room
	UserId int64 `json:"userId"`
	Role   int   `json:"role"`
}

type RoomMember struct {
	Id       int64     `json:"id"`
	UserId   int64     `json:"userId"`
	Role     int       `json:"role"`
	RoomId   int64     `json:"roomId"`
	JoinedAt time.Time `json:"joinedTime"`
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
	AddRoomMember(ctx context.Context, member *RoomMember) error
	GetOrCreateRoomSession(ctx context.Context, roomId int64) (*RoomSession, bool, error)
	GetRoomSession(ctx context.Context, roomId int64) (*RoomSession, error)
	RemoveRoomSession(ctx context.Context, roomId int64) error
}

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
	if room.Private {
		room.Password = generatePassword(4)
	}
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
	return room, nil
}

func (uc *RoomUseCase) ListRooms(ctx context.Context, page int, pageSize int) ([]*Room, int, error) {
	rooms, total, err := uc.rr.ListRooms(ctx, page, pageSize)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return make([]*Room, 0), 0, nil
	}
	if err != nil {
		return nil, 0, v1.ErrorGetRoomFailed("database error: %v", err)
	}
	return rooms, total, nil
}

func (uc *RoomUseCase) ListJoinedRooms(ctx context.Context, userId int64, page int, pageSize int) ([]*JoinedRoom, int, error) {
	rooms, total, err := uc.rr.ListJoinedRooms(ctx, userId, page, pageSize)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return make([]*JoinedRoom, 0), 0, nil
	}
	if err != nil {
		return nil, 0, v1.ErrorGetRoomFailed("database error: %v", err)
	}
	return rooms, total, nil
}

func (uc *RoomUseCase) JoinRoom(ctx context.Context, userId int64, roomId int64, password string) error {
	member, _ := uc.rr.GetRoomMember(ctx, userId, roomId)
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
			Id:       0,
			UserId:   userId,
			Role:     1,
			RoomId:   roomId,
			JoinedAt: time.Now(),
		})
		if err != nil {
			return v1.ErrorCreateRoomFailed("database error: %v", err)
		}
	}
	return nil
}

func (uc *RoomUseCase) GetRoomSession(ctx context.Context, roomId, userId int64, game string) (*RoomSession, error) {
	session, err := uc.rr.GetRoomSession(ctx, roomId)
	if session != nil {
		return session, nil
	}
	if err != nil {
		return nil, v1.ErrorGetRoomFailed("consul kv error: %v", err)
	}
	session, _, err = uc.rr.GetOrCreateRoomSession(ctx, roomId)
	if err != nil {
		return nil, v1.ErrorGetRoomFailed("create session error: %v", err)
	}
	if session == nil {
		return nil, v1.ErrorCreateRoomSessionFailed("no available service to create session")
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

func generatePassword(length int) string {
	dict := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456"
	builder := strings.Builder{}
	for i := 0; i < length; i++ {
		builder.WriteByte(dict[rand.Intn(len(dict))])
	}
	return builder.String()
}
