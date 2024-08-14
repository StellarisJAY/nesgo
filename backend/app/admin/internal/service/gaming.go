package service

import (
	"context"
	adminAPI "github.com/stellarisJAY/nesgo/backend/api/app/admin/v1"
)

func (s *AdminService) ListGamingServiceEndpoints(ctx context.Context, request *adminAPI.ListEndpointsRequest) (*adminAPI.ListEndpointsResponse, error) {
	endpoints, total, err := s.g.ListGamingServiceEndpoints(ctx, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	result := make([]*adminAPI.GamingServiceEndpoint, 0, len(endpoints))
	for _, endpoint := range endpoints {
		result = append(result, &adminAPI.GamingServiceEndpoint{
			Address:       endpoint.Address,
			Id:            endpoint.Id,
			CpuUsage:      endpoint.CpuUsage,
			MemoryUsed:    endpoint.MemoryUsed,
			MemoryTotal:   endpoint.MemoryTotal,
			EmulatorCount: endpoint.EmulatorCount,
			Uptime:        endpoint.Uptime,
		})
	}
	return &adminAPI.ListEndpointsResponse{Endpoints: result, Total: total}, nil
}

func (s *AdminService) ListActiveRooms(ctx context.Context, req *adminAPI.ListActiveRoomsRequest) (*adminAPI.ListActiveRoomsResponse, error) {
	rooms, err := s.g.ListActiveRoomsOnEndpoint(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	result := make([]*adminAPI.ActiveRoom, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, &adminAPI.ActiveRoom{
			RoomId:            room.Id,
			Name:              room.Name,
			Host:              room.Host,
			HostName:          room.HostName,
			MemberCount:       room.MemberCount,
			MemberLimit:       room.MemberLimit,
			Connections:       room.Connections,
			ActiveConnections: room.ActiveConnections,
			Game:              room.Game,
			Uptime:            room.Uptime,
		})
	}
	return &adminAPI.ListActiveRoomsResponse{Rooms: result}, nil
}
