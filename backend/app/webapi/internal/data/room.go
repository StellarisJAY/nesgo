package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	roomAPI "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
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

func (r *roomRepo) GetRoomSession(ctx context.Context, roomId, userId int64) (*biz.RoomSession, error) {
	response, err := r.data.rc.GetRoomSession(ctx, &roomAPI.GetRoomSessionRequest{
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
