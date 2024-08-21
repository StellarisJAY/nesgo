package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	roomAPI "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/biz"
)

type roomRepo struct {
	data   *Data
	logger *log.Helper
}

func NewRoomRepo(data *Data, logger log.Logger) biz.RoomRepo {
	return &roomRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/room")),
	}
}

func (r *roomRepo) AddDeleteRoomSessionTask(ctx context.Context, roomId int64, instanceId string) error {
	_, err := r.data.roomCli.AddDeleteRoomSessionTask(ctx, &roomAPI.AddDeleteRoomSessionTaskRequest{
		RoomId:     roomId,
		InstanceId: instanceId,
	})
	return err
}
