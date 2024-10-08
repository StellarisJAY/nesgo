package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
	roomAPI "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
)

func (ws *WebApiService) ListMyRooms(ctx context.Context, request *v1.ListRoomRequest) (*v1.ListRoomResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	rooms, total, err := ws.rc.ListJoinedRooms(ctx, c.UserId, int(request.Page), int(request.PageSize))
	if err != nil {
		return nil, err
	}
	result := make([]*v1.Room, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, &v1.Room{
			Id:           room.Id,
			Name:         room.Name,
			Host:         room.Host,
			Private:      room.Private,
			MemberCount:  room.MemberCount,
			MemberLimit:  room.MemberLimit,
			CreateTime:   room.CreateTime.UnixMilli(),
			HostName:     room.HostName,
			EmulatorType: room.EmulatorType,
		})
	}
	return &v1.ListRoomResponse{
		Rooms: result,
		Total: int32(total),
	}, nil
}

func (ws *WebApiService) CreateRoom(ctx context.Context, request *v1.CreateRoomRequest) (*v1.CreateRoomResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	room, err := ws.rc.CreateRoom(ctx, request.Name, request.Private, c.UserId, request.EmulatorType)
	if err != nil {
		return nil, err
	}
	return &v1.CreateRoomResponse{
		Id:          room.Id,
		Name:        room.Name,
		Host:        room.Host,
		Private:     room.Private,
		Password:    room.Password,
		MemberLimit: room.MemberLimit,
	}, nil
}

func (ws *WebApiService) GetRoom(ctx context.Context, request *v1.GetRoomRequest) (*v1.GetRoomResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	room, err := ws.rc.GetRoom(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	resp := &v1.GetRoomResponse{
		Id:           room.Id,
		Name:         room.Name,
		Host:         room.Host,
		Private:      room.Private,
		MemberCount:  room.MemberCount,
		HostName:     room.HostName,
		MemberLimit:  room.MemberLimit,
		CreateTime:   room.CreateTime.UnixMilli(),
		EmulatorType: room.EmulatorType,
	}
	if c.UserId == room.Host {
		resp.Password = room.Password
	}
	return resp, nil
}

func (ws *WebApiService) ListAllRooms(ctx context.Context, request *v1.ListRoomRequest) (*v1.ListRoomResponse, error) {
	rooms, total, err := ws.rc.ListRooms(ctx, int(request.Page), int(request.PageSize))
	if err != nil {
		return nil, err
	}
	result := make([]*v1.Room, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, &v1.Room{
			Id:           room.Id,
			Name:         room.Name,
			Host:         room.Host,
			Private:      room.Private,
			MemberCount:  room.MemberCount,
			HostName:     room.HostName,
			MemberLimit:  room.MemberLimit,
			CreateTime:   room.CreateTime.UnixMilli(),
			EmulatorType: room.EmulatorType,
		})
	}
	return &v1.ListRoomResponse{
		Rooms: result,
		Total: int32(total),
	}, nil
}

func (ws *WebApiService) ListMembers(ctx context.Context, request *v1.ListMemberRequest) (*v1.ListMemberResponse, error) {
	members, err := ws.rc.ListMembers(ctx, request.RoomId)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.Member, 0, len(members))
	for _, member := range members {
		result = append(result, &v1.Member{
			UserId:   member.Id,
			Name:     member.Name,
			Role:     member.Role,
			JoinedAt: member.JoinedAt.UnixMilli(),
			Player1:  member.Player1,
			Player2:  member.Player2,
		})
	}
	return &v1.ListMemberResponse{
		Members: result,
	}, nil
}

func (ws *WebApiService) JoinRoom(ctx context.Context, request *v1.JoinRoomRequest) (*v1.JoinRoomResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	err := ws.rc.JoinRoom(ctx, request.RoomId, c.UserId, request.Password)
	if err != nil {
		return nil, err
	}
	return &v1.JoinRoomResponse{
		RoomId: request.RoomId,
		UserId: c.UserId,
		Role:   roomAPI.RoomRole_Observer,
	}, nil
}

func (ws *WebApiService) UpdateRoom(ctx context.Context, request *v1.UpdateRoomRequest) (*v1.UpdateRoomResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	room := &biz.Room{
		Id:      request.RoomId,
		Name:    request.Name,
		Private: request.Private,
	}
	err := ws.rc.UpdateRoom(ctx, room, c.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRoomResponse{
		RoomId:   request.RoomId,
		Name:     request.Name,
		Private:  room.Private,
		Password: room.Password,
	}, nil
}

func (ws *WebApiService) DeleteRoom(ctx context.Context, request *v1.DeleteRoomRequest) (*v1.DeleteRoomResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	err := ws.rc.DeleteRoom(ctx, request.RoomId, c.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteRoomResponse{}, nil
}

func (ws *WebApiService) GetRoomMember(ctx context.Context, request *v1.GetRoomMemberRequest) (*v1.GetRoomMemberResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	member, err := ws.rc.GetRoomMember(ctx, request.RoomId, c.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetRoomMemberResponse{Member: &v1.Member{
		UserId:   member.Id,
		Name:     member.Name,
		Role:     member.Role,
		JoinedAt: member.JoinedAt.UnixMilli(),
	}}, nil
}

func (ws *WebApiService) UpdateMemberRole(ctx context.Context, request *v1.UpdateMemberRoleRequest) (*v1.UpdateMemberRoleResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	err := ws.rc.UpdateMemberRole(ctx, request.RoomId, request.UserId, request.Role, c.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateMemberRoleResponse{}, nil
}

func (ws *WebApiService) DeleteMember(ctx context.Context, request *v1.DeleteMemberRequest) (*v1.DeleteMemberResponse, error) {
	claims, _ := jwt.FromContext(ctx)
	c := claims.(*biz.LoginClaims)
	err := ws.rc.DeleteMember(ctx, request.RoomId, request.UserId, c.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteMemberResponse{}, nil
}
