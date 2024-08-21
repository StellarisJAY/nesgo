package service

import (
	"context"
	adminAPI "github.com/stellarisJAY/nesgo/backend/api/app/admin/v1"
)

func (s *AdminService) ListRooms(ctx context.Context, request *adminAPI.ListRoomsRequest) (*adminAPI.ListRoomsResponse, error) {
	rooms, total, err := s.r.ListRooms(ctx, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}

	result := make([]*adminAPI.Room, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, &adminAPI.Room{
			Id:          room.Id,
			Name:        room.Name,
			Host:        room.Host,
			HostName:    room.HostName,
			Private:     room.Private,
			MemberCount: room.MemberCount,
			MemberLimit: room.MemberLimit,
		})
	}
	return &adminAPI.ListRoomsResponse{
		Total: total,
		Rooms: result,
	}, nil
}

func (s *AdminService) GetRoomStats(ctx context.Context, request *adminAPI.GetRoomStatsRequest) (*adminAPI.GetRoomStatsResponse, error) {
	stats, err := s.r.GetRoomStats(ctx, request.RoomId)
	if err != nil {
		return nil, err
	}
	return &adminAPI.GetRoomStatsResponse{
		Stats: &adminAPI.ActiveRoom{
			RoomId:            stats.Id,
			Name:              stats.Name,
			Host:              stats.Host,
			HostName:          stats.HostName,
			MemberCount:       stats.MemberCount,
			MemberLimit:       stats.MemberLimit,
			Endpoint:          stats.Endpoint,
			Connections:       stats.Connections,
			ActiveConnections: stats.ActiveConnections,
			Game:              stats.Game,
			Uptime:            stats.Uptime,
		},
	}, nil
}
