package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type RoomGameStats struct {
	Connections       int32  `json:"connections"`
	ActiveConnections int32  `json:"activeConnections"`
	Game              string `json:"game"`
	Uptime            int64  `json:"uptime"`
}

type ServiceEndpoint struct {
	Address string `json:"address"`
	Id      string `json:"id"`
}

type GamingRepo interface {
	GetRoomGameStats(ctx context.Context, roomId int64, endpoint string) (*RoomGameStats, error)
	ListGamingServiceEndpoints(ctx context.Context, page, pageSize int32) ([]*ServiceEndpoint, int32, error)
	ListActiveRoomsOnEndpoint(ctx context.Context, id string) ([]*RoomStats, error)
}

type GamingUseCase struct {
	repo   GamingRepo
	rr     RoomRepo
	ur     UserRepo
	logger *log.Helper
}

func NewGamingUseCase(repo GamingRepo, rr RoomRepo, ur UserRepo, logger log.Logger) *GamingUseCase {
	return &GamingUseCase{
		repo:   repo,
		rr:     rr,
		ur:     ur,
		logger: log.NewHelper(log.With(logger, "module", "biz/gaming")),
	}
}

func (g *GamingUseCase) ListGamingServiceEndpoints(ctx context.Context, page, pageSize int32) ([]*ServiceEndpoint, int32, error) {
	endpoints, total, err := g.repo.ListGamingServiceEndpoints(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return endpoints, total, nil
}

func (g *GamingUseCase) ListActiveRoomsOnEndpoint(ctx context.Context, id string) ([]*RoomStats, error) {
	rooms, err := g.repo.ListActiveRoomsOnEndpoint(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, stat := range rooms {
		r, _ := g.rr.GetRoom(ctx, stat.Id)
		if r != nil {
			stat.Name = r.Name
			stat.Private = r.Private
			stat.Host = r.Host
			stat.MemberCount = r.MemberCount
			stat.MemberLimit = r.MemberLimit
			user, _ := g.ur.GetUser(ctx, stat.Host)
			if user != nil {
				stat.HostName = user.Name
			}
		}
	}
	return rooms, nil
}
