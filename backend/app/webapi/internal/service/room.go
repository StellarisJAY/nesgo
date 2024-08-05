package service

import (
	"context"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
)

func (ws *WebApiService) ListMyRooms(ctx context.Context, request *v1.ListMyRoomsRequest) (*v1.ListMyRoomResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ws *WebApiService) GetRoomSession(ctx context.Context, request *v1.GetRoomSessionRequest) (*v1.GetRoomSessionResponse, error) {
	session, err := ws.rc.GetRoomSession(ctx, request.RoomId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetRoomSessionResponse{
		RoomId:   session.RoomId,
		Endpoint: session.Endpoint,
	}, nil
}
