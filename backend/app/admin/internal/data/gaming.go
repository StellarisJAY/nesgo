package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	gamingAPI "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/biz"
	"net/url"
	"slices"
	"strings"
)

type gamingRepo struct {
	data   *Data
	logger *log.Helper
}

func NewGamingRepo(data *Data, logger log.Logger) biz.GamingRepo {
	return &gamingRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/gaming")),
	}
}

func (g *gamingRepo) GetRoomGameStats(ctx context.Context, roomId int64, endpoint string) (*biz.RoomGameStats, error) {
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(endpoint))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	gamingCli := gamingAPI.NewGamingClient(conn)
	response, err := gamingCli.GetGameInstanceStats(ctx, &gamingAPI.GetGameInstanceStatsRequest{RoomId: roomId})
	if err != nil {
		return nil, err
	}
	return &biz.RoomGameStats{
		Connections:       response.Stats.Connections,
		ActiveConnections: response.Stats.ActiveConnections,
		Game:              response.Stats.Game,
		Uptime:            response.Stats.Uptime,
	}, nil
}

func (g *gamingRepo) ListGamingServiceEndpoints(ctx context.Context, page, pageSize int32) ([]*biz.ServiceEndpoint, int32, error) {
	instances, err := g.data.discovery.GetService(ctx, "nesgo.service.gaming")
	if err != nil {
		return nil, 0, err
	}
	slices.SortFunc(instances, func(x, y *registry.ServiceInstance) int {
		return strings.Compare(x.ID, y.ID)
	})
	endpoints := make([]*biz.ServiceEndpoint, 0, pageSize)
	total := int32(len(instances))
	for i := page * pageSize; i < total && i < (page+1)*pageSize; i++ {
		u, _ := url.Parse(instances[i].Endpoints[0])
		address := u.Host
		endpoint := &biz.ServiceEndpoint{
			Address: address,
			Id:      instances[i].ID,
		}
		conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(address))
		if err != nil {
			g.logger.Errorf("get endponint stats error: %v", err)
			endpoints = append(endpoints, endpoint)
			continue
		}
		stats, err := gamingAPI.NewGamingClient(conn).GetEndpointStats(ctx, &gamingAPI.GetEndpointStatsRequest{})
		if err == nil {
			endpoint.CpuUsage = stats.CpuUsage
			endpoint.MemoryUsed = stats.MemoryUsed
			endpoint.MemoryTotal = stats.MemoryTotal
			endpoint.Uptime = stats.Uptime
			endpoint.EmulatorCount = stats.EmulatorCount
		} else {
			g.logger.Errorf("get endpoint stats error: %v", err)
		}
		_ = conn.Close()
		endpoints = append(endpoints, endpoint)
	}
	return endpoints, total, nil
}

func (g *gamingRepo) ListActiveRoomsOnEndpoint(ctx context.Context, id string) ([]*biz.RoomStats, error) {
	instances, err := g.data.discovery.GetService(ctx, "nesgo.service.gaming")
	if err != nil {
		return nil, err
	}
	idx := slices.IndexFunc(instances, func(instance *registry.ServiceInstance) bool {
		return instance.ID == id
	})
	if idx == -1 {
		return nil, nil
	}
	u, _ := url.Parse(instances[idx].Endpoints[0])
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(u.Host))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	gamingCli := gamingAPI.NewGamingClient(conn)
	response, err := gamingCli.ListGameInstances(ctx, &gamingAPI.ListGameInstancesRequest{})
	if err != nil {
		return nil, err
	}
	stats := make([]*biz.RoomStats, 0, len(response.Instances))
	for _, instance := range response.Instances {
		stats = append(stats, &biz.RoomStats{
			Room: biz.Room{Id: instance.RoomId},
			RoomGameStats: biz.RoomGameStats{
				Connections:       instance.Connections,
				ActiveConnections: instance.ActiveConnections,
				Game:              instance.Game,
				Uptime:            instance.Uptime,
			},
		})
	}
	return stats, nil
}
