package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	roomAPI "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
)

type roomRepo struct {
	data    *Data
	roomCli roomAPI.RoomClient
	logger  *log.Helper
}

func NewRoomRepo(data *Data, r registry.Discovery, logger log.Logger) biz.RoomRepo {
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///nesgo.service.room"),
		grpc.WithDiscovery(r))
	if err != nil {
		panic(err)
	}
	roomCli := roomAPI.NewRoomClient(conn)
	return &roomRepo{
		data:    data,
		roomCli: roomCli,
		logger:  log.NewHelper(log.With(logger, "module", "data/room")),
	}
}

func (r *roomRepo) GetRoomSession(ctx context.Context, roomId, userId int64) (*biz.RoomSession, error) {
	response, err := r.roomCli.GetRoomSession(ctx, &roomAPI.GetRoomSessionRequest{
		RoomId: roomId,
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	return &biz.RoomSession{
		RoomId:   roomId,
		Endpoint: response.Endpoint,
	}, nil
}
