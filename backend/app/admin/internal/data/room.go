package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	roomAPI "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/biz"
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

func (r *roomRepo) ListRooms(ctx context.Context, page int32, pageSize int32) ([]*biz.Room, int32, error) {
	response, err := r.data.roomCli.ListRooms(ctx, &roomAPI.ListRoomsRequest{
		Page:     page,
		PageSize: pageSize,
		Joined:   false,
	})
	if err != nil {
		return nil, 0, err
	}

	rooms := make([]*biz.Room, 0, len(response.Rooms))
	for _, room := range response.Rooms {
		rooms = append(rooms, &biz.Room{
			Id:          room.Id,
			Name:        room.Name,
			Host:        room.Host,
			Private:     room.Private,
			MemberCount: room.MemberCount,
			MemberLimit: room.MemberLimit,
		})
	}
	return rooms, response.Total, nil
}

func (r *roomRepo) GetRoom(ctx context.Context, id int64) (*biz.Room, error) {
	room, err := r.data.roomCli.GetRoom(ctx, &roomAPI.GetRoomRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &biz.Room{
		Id:          room.Id,
		Name:        room.Name,
		Host:        room.Host,
		Private:     room.Private,
		MemberCount: room.MemberCount,
		MemberLimit: room.MemberLimit,
	}, nil
}

func (r *roomRepo) GetRoomEndpoint(ctx context.Context, id int64) (string, error) {
	session, err := r.data.roomCli.GetRoomSession(ctx, &roomAPI.GetRoomSessionRequest{RoomId: id})
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", nil
	}
	return session.Endpoint, nil
}
