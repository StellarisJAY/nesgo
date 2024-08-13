package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	gamingAPI "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	roomAPI "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	"time"
)

type RoomUseCase struct {
	repo   RoomRepo
	ur     UserRepo
	logger *log.Helper
}

type RoomSession struct {
	RoomId   int64  `json:"roomId"`
	Endpoint string `json:"endpoint"`
}

type Room struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Host        int64     `json:"host"`
	HostName    string    `json:"hostName"`
	Private     bool      `json:"private"`
	MemberCount int32     `json:"memberCount"`
	MemberLimit int32     `json:"memberLimit"`
	Password    string    `json:"password"`
	CreateTime  time.Time `json:"createTime"`
}

type Member struct {
	Id       int64            `json:"id"`
	Name     string           `json:"name"`
	Role     roomAPI.RoomRole `json:"role"`
	JoinedAt time.Time        `json:"joinedAt"`
	Player1  bool             `json:"player1"`
	Player2  bool             `json:"player2"`
}

type JoinedRoom struct {
	Room
	Role roomAPI.RoomRole `json:"role"`
}

type RoomRepo interface {
	GetRoomSession(ctx context.Context, roomId int64) (*RoomSession, error)
	GetCreateRoomSession(ctx context.Context, roomId, userId int64, selectedGame string) (*RoomSession, error)
	CreateRoom(ctx context.Context, room *Room) error
	GetRoom(ctx context.Context, roomId int64) (*Room, error)
	ListJoinedRooms(ctx context.Context, userId int64, page, pageSize int) ([]*JoinedRoom, int, error)
	ListRooms(ctx context.Context, page, pageSize int) ([]*Room, int, error)
	ListMembers(ctx context.Context, roomId int64) ([]*Member, error)
	JoinRoom(ctx context.Context, roomId, userId int64, password string) error
	UpdateRoom(ctx context.Context, room *Room, userId int64) error
	DeleteRoom(ctx context.Context, roomId, userId int64) error
	GetRoomMember(ctx context.Context, roomId, userId int64) (*Member, error)
}

func NewRoomUseCase(repo RoomRepo, ur UserRepo, logger log.Logger) *RoomUseCase {
	return &RoomUseCase{
		repo:   repo,
		ur:     ur,
		logger: log.NewHelper(log.With(logger, "module", "biz/room")),
	}
}

func (uc *RoomUseCase) CreateRoom(ctx context.Context, name string, private bool, userId int64) (*Room, error) {
	room := &Room{
		Name:    name,
		Private: private,
		Host:    userId,
	}
	err := uc.repo.CreateRoom(ctx, room)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (uc *RoomUseCase) GetRoom(ctx context.Context, roomId int64) (*Room, error) {
	room, err := uc.repo.GetRoom(ctx, roomId)
	if err != nil {
		return nil, err
	}
	user, err := uc.ur.GetUser(ctx, room.Host)
	if err != nil {
		return nil, err
	}
	room.HostName = user.Name
	return room, nil
}

func (uc *RoomUseCase) ListJoinedRooms(ctx context.Context, userId int64, page, pageSize int) ([]*JoinedRoom, int, error) {
	rooms, total, err := uc.repo.ListJoinedRooms(ctx, userId, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for _, room := range rooms {
		user, _ := uc.ur.GetUser(ctx, room.Host)
		if user != nil {
			room.HostName = user.Name
		}
	}
	return rooms, total, nil
}

func (uc *RoomUseCase) ListRooms(ctx context.Context, page, pageSize int) ([]*Room, int, error) {
	rooms, total, err := uc.repo.ListRooms(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for _, room := range rooms {
		user, _ := uc.ur.GetUser(ctx, room.Host)
		if user != nil {
			room.HostName = user.Name
		}
	}
	return rooms, total, nil
}

func (uc *RoomUseCase) ListMembers(ctx context.Context, roomId int64) ([]*Member, error) {
	members, err := uc.repo.ListMembers(ctx, roomId)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		user, _ := uc.ur.GetUser(ctx, member.Id)
		if user != nil {
			member.Name = user.Name
		}
	}

	session, _ := uc.repo.GetRoomSession(ctx, roomId)
	if session == nil {
		return members, nil
	}

	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(session.Endpoint))
	if err != nil {
		uc.logger.Errorf("list members connect gaming service failed: %v", err)
		return members, nil
	}
	defer conn.Close()
	response, err := gamingAPI.NewGamingClient(conn).GetControllers(ctx, &gamingAPI.GetControllersRequest{RoomId: roomId})
	if err != nil {
		uc.logger.Errorf("list members get player1&2 failed: %v", err)
		return members, nil
	}
	for _, member := range members {
		member.Player1 = member.Id == response.Controller1
		member.Player2 = member.Id == response.Controller2
	}
	return members, nil
}

func (uc *RoomUseCase) JoinRoom(ctx context.Context, roomId, userId int64, password string) error {
	return uc.repo.JoinRoom(ctx, roomId, userId, password)
}

func (uc *RoomUseCase) UpdateRoom(ctx context.Context, room *Room, userId int64) error {
	return uc.repo.UpdateRoom(ctx, room, userId)
}

func (uc *RoomUseCase) DeleteRoom(ctx context.Context, roomId, userId int64) error {
	return uc.repo.DeleteRoom(ctx, roomId, userId)
}

func (uc *RoomUseCase) GetRoomMember(ctx context.Context, roomId, userId int64) (*Member, error) {
	member, err := uc.repo.GetRoomMember(ctx, roomId, userId)
	if err != nil {
		return nil, err
	}
	user, err := uc.ur.GetUser(ctx, member.Id)
	if err != nil {
		return nil, err
	}
	member.Name = user.Name
	return member, nil
}
