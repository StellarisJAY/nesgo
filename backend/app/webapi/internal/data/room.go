package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	roomAPI "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
	"time"
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

func (r *roomRepo) GetRoomSession(ctx context.Context, roomId int64) (*biz.RoomSession, error) {
	response, err := r.data.rc.GetRoomSession(ctx, &roomAPI.GetRoomSessionRequest{
		RoomId: roomId,
	})
	if err != nil {
		return nil, err
	}
	return &biz.RoomSession{
		RoomId:   roomId,
		Endpoint: response.Endpoint,
	}, nil
}

func (r *roomRepo) GetCreateRoomSession(ctx context.Context, roomId, userId int64, selectedGame string) (*biz.RoomSession, error) {
	response, err := r.data.rc.GetCreateRoomSession(ctx, &roomAPI.GetCreateRoomSessionRequest{
		RoomId:       roomId,
		UserId:       userId,
		SelectedGame: selectedGame,
	})
	if err != nil {
		return nil, err
	}
	return &biz.RoomSession{
		RoomId:   roomId,
		Endpoint: response.Endpoint,
	}, nil
}

func (r *roomRepo) CreateRoom(ctx context.Context, room *biz.Room) error {
	response, err := r.data.rc.CreateRoom(ctx, &roomAPI.CreateRoomRequest{
		Name:    room.Name,
		Host:    room.Host,
		Private: room.Private,
	})
	if err != nil {
		return err
	}
	room.Id = response.Id
	room.Password = response.Password
	room.MemberLimit = response.MemberLimit
	return nil
}

func (r *roomRepo) GetRoom(ctx context.Context, roomId int64) (*biz.Room, error) {
	response, err := r.data.rc.GetRoom(ctx, &roomAPI.GetRoomRequest{Id: roomId})
	if err != nil {
		return nil, err
	}
	return &biz.Room{
		Id:          response.Id,
		Name:        response.Name,
		Host:        response.Host,
		Private:     response.Private,
		Password:    response.Password,
		MemberCount: response.MemberCount,
		MemberLimit: response.MemberLimit,
		CreateTime:  time.UnixMilli(response.CreateTime).Local(),
	}, nil
}

func (r *roomRepo) ListJoinedRooms(ctx context.Context, userId int64, page, pageSize int) ([]*biz.JoinedRoom, int, error) {
	response, err := r.data.rc.ListRooms(ctx, &roomAPI.ListRoomsRequest{
		UserId:   userId,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Joined:   true,
	})
	if err != nil {
		return nil, 0, err
	}
	result := make([]*biz.JoinedRoom, 0, len(response.Rooms))
	for _, room := range response.Rooms {
		result = append(result, &biz.JoinedRoom{
			Room: biz.Room{
				Id:          room.Id,
				Name:        room.Name,
				Host:        room.Host,
				Private:     room.Private,
				Password:    room.Password,
				MemberCount: room.MemberCount,
				MemberLimit: room.MemberLimit,
				CreateTime:  time.UnixMilli(room.CreateTime).Local(),
			},
			Role: room.Role,
		})
	}
	return result, int(response.Total), nil
}

func (r *roomRepo) ListRooms(ctx context.Context, page, pageSize int) ([]*biz.Room, int, error) {
	response, err := r.data.rc.ListRooms(ctx, &roomAPI.ListRoomsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Joined:   false,
	})
	if err != nil {
		return nil, 0, err
	}
	result := make([]*biz.Room, 0, len(response.Rooms))
	for _, room := range response.Rooms {
		result = append(result, &biz.Room{
			Id:          room.Id,
			Name:        room.Name,
			Host:        room.Host,
			Private:     room.Private,
			Password:    room.Password,
			MemberCount: room.MemberCount,
			MemberLimit: room.MemberLimit,
			CreateTime:  time.UnixMilli(room.CreateTime).Local(),
		})
	}
	return result, int(response.Total), nil
}

func (r *roomRepo) ListMembers(ctx context.Context, roomId int64) ([]*biz.Member, error) {
	members, err := r.data.rc.ListRoomMembers(ctx, &roomAPI.ListRoomMemberRequest{
		Id: roomId,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*biz.Member, 0, len(members.Members))
	for _, member := range members.Members {
		result = append(result, &biz.Member{
			Id:       member.UserId,
			Role:     member.Role,
			JoinedAt: time.UnixMilli(member.JoinedAt).Local(),
		})
	}
	return result, nil
}

func (r *roomRepo) JoinRoom(ctx context.Context, roomId, userId int64, password string) error {
	_, err := r.data.rc.JoinRoom(ctx, &roomAPI.JoinRoomRequest{
		Id:       roomId,
		UserId:   userId,
		Password: password,
	})
	return err
}

func (r *roomRepo) UpdateRoom(ctx context.Context, room *biz.Room, userId int64) error {
	response, err := r.data.rc.UpdateRoom(ctx, &roomAPI.UpdateRoomRequest{
		RoomId:  room.Id,
		Name:    room.Name,
		Private: room.Private,
		UserId:  userId,
	})
	if err != nil {
		return err
	}
	room.Password = response.Password
	return nil
}

func (r *roomRepo) DeleteRoom(ctx context.Context, roomId, userId int64) error {
	_, err := r.data.rc.DeleteRoom(ctx, &roomAPI.DeleteRoomRequest{
		RoomId: roomId,
		UserId: userId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *roomRepo) GetRoomMember(ctx context.Context, roomId, userId int64) (*biz.Member, error) {
	response, err := r.data.rc.GetRoomMember(ctx, &roomAPI.GetRoomMemberRequest{
		RoomId: roomId,
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	return &biz.Member{
		Id:       response.Member.UserId,
		Role:     response.Member.Role,
		JoinedAt: time.UnixMilli(response.Member.JoinedAt).Local(),
	}, nil
}

func (r *roomRepo) UpdateMember(ctx context.Context, roomId, userId int64, role roomAPI.RoomRole) error {
	_, err := r.data.rc.UpdateMember(ctx, &roomAPI.UpdateMemberRequest{
		RoomId: roomId,
		UserId: userId,
		Role:   role,
	})
	return err
}

func (r *roomRepo) DeleteMember(ctx context.Context, roomId, userId int64) error {
	_, err := r.data.rc.DeleteMember(ctx, &roomAPI.DeleteMemberRequest{
		RoomId: roomId,
		UserId: userId,
	})
	return err
}
