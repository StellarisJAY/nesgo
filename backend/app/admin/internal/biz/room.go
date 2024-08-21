package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type Room struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Host        int64  `json:"host"`
	HostName    string `json:"hostName"`
	Private     bool   `json:"private"`
	MemberCount int32  `json:"memberCount"`
	MemberLimit int32  `json:"memberLimit"`
	Endpoint    string `json:"endpoint"`
}

type RoomStats struct {
	Room
	RoomGameStats
}

type RoomUseCase struct {
	repo   RoomRepo
	ur     UserRepo
	gr     GamingRepo
	logger *log.Helper
}

type RoomRepo interface {
	ListRooms(ctx context.Context, page int32, pageSize int32) ([]*Room, int32, error)
	GetRoom(ctx context.Context, id int64) (*Room, error)
	GetRoomEndpoint(ctx context.Context, id int64) (string, error)
}

func NewRoomUseCase(repo RoomRepo, ur UserRepo, gr GamingRepo, logger log.Logger) *RoomUseCase {
	return &RoomUseCase{
		repo:   repo,
		ur:     ur,
		gr:     gr,
		logger: log.NewHelper(log.With(logger, "module", "biz/room")),
	}
}

func (uc *RoomUseCase) ListRooms(ctx context.Context, page int32, pageSize int32) ([]*Room, int32, error) {
	rooms, total, err := uc.repo.ListRooms(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for _, room := range rooms {
		host, err := uc.ur.GetUser(ctx, room.Host)
		if err != nil {
			return nil, 0, err
		}
		room.HostName = host.Name
	}
	return rooms, total, nil
}

func (uc *RoomUseCase) GetRoomStats(ctx context.Context, roomId int64) (*RoomStats, error) {
	room, err := uc.repo.GetRoom(ctx, roomId)
	if err != nil {
		return nil, err
	}

	host, err := uc.ur.GetUser(ctx, room.Host)
	if err != nil {
		return nil, err
	}
	room.HostName = host.Name

	endpoint, _ := uc.repo.GetRoomEndpoint(ctx, room.Id)
	stats := &RoomStats{
		Room: *room,
	}
	if endpoint == "" {
		return stats, nil
	}
	stats.Endpoint = endpoint
	gameStats, err := uc.gr.GetRoomGameStats(ctx, roomId, endpoint)
	if err != nil {
		return nil, err
	}
	stats.RoomGameStats = *gameStats
	return stats, nil
}
