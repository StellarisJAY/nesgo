package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/room/internal/biz"
)

type RoomService struct {
	v1.UnimplementedRoomServer
	ruc    *biz.RoomUseCase
	logger *log.Helper
}

func NewRoomService(ruc *biz.RoomUseCase, logger log.Logger) *RoomService {
	return &RoomService{
		ruc:    ruc,
		logger: log.NewHelper(logger),
	}
}

func (r *RoomService) CreateRoom(ctx context.Context, request *v1.CreateRoomRequest) (*v1.CreateRoomResponse, error) {
	room := &biz.Room{
		Id:           0,
		Name:         request.Name,
		Host:         request.Host,
		Private:      request.Private,
		EmulatorType: request.EmulatorType,
	}
	err := r.ruc.CreateRoom(ctx, room)
	if err != nil {
		return nil, err
	}
	return &v1.CreateRoomResponse{
		Id:          room.Id,
		Private:     room.Private,
		Password:    room.Password,
		MemberLimit: int32(room.MemberLimit),
	}, nil
}

func (r *RoomService) GetRoom(ctx context.Context, request *v1.GetRoomRequest) (*v1.GetRoomResponse, error) {
	room, err := r.ruc.GetRoom(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetRoomResponse{
		Id:           room.Id,
		Name:         room.Name,
		Host:         room.Host,
		Private:      room.Private,
		Password:     room.Password,
		MemberCount:  int32(room.MemberCount),
		MemberLimit:  int32(room.MemberLimit),
		CreateTime:   room.CreateTime.UnixMilli(),
		EmulatorType: room.EmulatorType,
	}, nil
}

func (r *RoomService) ListRoomMembers(ctx context.Context, request *v1.ListRoomMemberRequest) (*v1.ListRoomMemberResponse, error) {
	members, err := r.ruc.ListRoomMembers(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.RoomMember, 0, len(members))
	for _, m := range members {
		result = append(result, &v1.RoomMember{
			UserId:   m.UserId,
			Role:     m.Role,
			JoinedAt: m.JoinedAt.UnixMilli(),
		})
	}
	return &v1.ListRoomMemberResponse{Members: result}, nil
}

func (r *RoomService) ListRooms(ctx context.Context, request *v1.ListRoomsRequest) (*v1.ListRoomsResponse, error) {
	if request.Joined {
		rooms, total, err := r.ruc.ListJoinedRooms(ctx, request.UserId, int(request.Page), int(request.PageSize))
		if err != nil {
			return nil, err
		}
		result := make([]*v1.GetRoomResponse, 0, len(rooms))
		for _, room := range rooms {
			result = append(result, &v1.GetRoomResponse{
				Id:           room.Id,
				Name:         room.Name,
				Host:         room.Host,
				Private:      room.Private,
				MemberCount:  int32(room.MemberCount),
				MemberLimit:  int32(room.MemberLimit),
				Role:         room.Role,
				CreateTime:   room.CreateTime.UnixMilli(),
				EmulatorType: room.EmulatorType,
			})
		}
		return &v1.ListRoomsResponse{Rooms: result, Total: int32(total)}, nil
	} else {
		rooms, total, err := r.ruc.ListRooms(ctx, int(request.Page), int(request.PageSize))
		if err != nil {
			return nil, err
		}
		result := make([]*v1.GetRoomResponse, 0, len(rooms))
		for _, room := range rooms {
			result = append(result, &v1.GetRoomResponse{
				Id:           room.Id,
				Name:         room.Name,
				Host:         room.Host,
				Private:      room.Private,
				MemberCount:  int32(room.MemberCount),
				MemberLimit:  int32(room.MemberLimit),
				CreateTime:   room.CreateTime.UnixMilli(),
				EmulatorType: room.EmulatorType,
			})
		}
		return &v1.ListRoomsResponse{Rooms: result, Total: int32(total)}, nil
	}
}

func (r *RoomService) JoinRoom(ctx context.Context, request *v1.JoinRoomRequest) (*v1.JoinRoomResponse, error) {
	err := r.ruc.JoinRoom(ctx, request.UserId, request.Id, request.Password)
	if err != nil {
		return nil, err
	}
	return &v1.JoinRoomResponse{}, nil
}

func (r *RoomService) GetRoomSession(ctx context.Context, request *v1.GetRoomSessionRequest) (*v1.GetRoomSessionResponse, error) {
	session, err := r.ruc.GetRoomSession(ctx, request.RoomId)
	if err != nil {
		return nil, err
	}
	return &v1.GetRoomSessionResponse{
		RoomId:   session.RoomId,
		Endpoint: session.Endpoint,
	}, nil
}

func (r *RoomService) GetCreateRoomSession(ctx context.Context, request *v1.GetCreateRoomSessionRequest) (*v1.GetRoomSessionResponse, error) {
	session, err := r.ruc.GetCreateRoomSession(ctx, request.RoomId, request.UserId, request.SelectedGame)
	if err != nil {
		return nil, err
	}
	return &v1.GetRoomSessionResponse{
		RoomId:   session.RoomId,
		Endpoint: session.Endpoint,
	}, nil
}

func (r *RoomService) RemoveRoomSession(ctx context.Context, request *v1.RemoveRoomSessionRequest) (*v1.RemoveRoomSessionResponse, error) {
	err := r.ruc.RemoveRoomSession(ctx, request.RoomId)
	if err != nil {
		return nil, err
	}
	return &v1.RemoveRoomSessionResponse{}, nil
}

func (r *RoomService) UpdateRoom(ctx context.Context, request *v1.UpdateRoomRequest) (*v1.UpdateRoomResponse, error) {
	room := &biz.Room{
		Id:      request.RoomId,
		Name:    request.Name,
		Private: request.Private,
	}
	err := r.ruc.UpdateRoom(ctx, room, request.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRoomResponse{
		RoomId:   room.Id,
		Name:     room.Name,
		Private:  room.Private,
		Password: room.Password,
	}, nil
}

func (r *RoomService) DeleteRoom(ctx context.Context, request *v1.DeleteRoomRequest) (*v1.DeleteRoomResponse, error) {
	err := r.ruc.DeleteRoom(ctx, request.RoomId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteRoomResponse{}, nil
}

func (r *RoomService) GetRoomMember(ctx context.Context, request *v1.GetRoomMemberRequest) (*v1.GetRoomMemberResponse, error) {
	member, err := r.ruc.GetRoomMember(ctx, request.RoomId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetRoomMemberResponse{
		Member: &v1.RoomMember{
			UserId:   member.UserId,
			Role:     member.Role,
			JoinedAt: member.JoinedAt.UnixMilli(),
		},
	}, nil
}

func (r *RoomService) UpdateMember(ctx context.Context, request *v1.UpdateMemberRequest) (*v1.UpdateMemberResponse, error) {
	err := r.ruc.UpdateMember(ctx, &biz.RoomMember{
		UserId: request.UserId,
		Role:   request.Role,
		RoomId: request.RoomId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateMemberResponse{}, nil
}

func (r *RoomService) DeleteMember(ctx context.Context, request *v1.DeleteMemberRequest) (*v1.DeleteMemberResponse, error) {
	err := r.ruc.DeleteMember(ctx, request.RoomId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteMemberResponse{}, nil
}

func (r *RoomService) AddDeleteRoomSessionTask(ctx context.Context, request *v1.AddDeleteRoomSessionTaskRequest) (*v1.AddDeleteRoomSessionTaskResponse, error) {
	err := r.ruc.AddDeleteRoomSessionTask(ctx, request.RoomId, request.InstanceId)
	if err != nil {
		return nil, err
	}
	return &v1.AddDeleteRoomSessionTaskResponse{}, nil
}
