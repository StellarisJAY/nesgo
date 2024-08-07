// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.20.3
// source: room/service/v1/room.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Room_CreateRoom_FullMethodName        = "/room.v1.Room/CreateRoom"
	Room_GetRoom_FullMethodName           = "/room.v1.Room/GetRoom"
	Room_ListRoomMembers_FullMethodName   = "/room.v1.Room/ListRoomMembers"
	Room_ListRooms_FullMethodName         = "/room.v1.Room/ListRooms"
	Room_JoinRoom_FullMethodName          = "/room.v1.Room/JoinRoom"
	Room_GetRoomSession_FullMethodName    = "/room.v1.Room/GetRoomSession"
	Room_RemoveRoomSession_FullMethodName = "/room.v1.Room/RemoveRoomSession"
)

// RoomClient is the client API for Room service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomClient interface {
	CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomResponse, error)
	GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error)
	ListRoomMembers(ctx context.Context, in *ListRoomMemberRequest, opts ...grpc.CallOption) (*ListRoomMemberResponse, error)
	ListRooms(ctx context.Context, in *ListRoomsRequest, opts ...grpc.CallOption) (*ListRoomsResponse, error)
	JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...grpc.CallOption) (*JoinRoomResponse, error)
	GetRoomSession(ctx context.Context, in *GetRoomSessionRequest, opts ...grpc.CallOption) (*GetRoomSessionResponse, error)
	RemoveRoomSession(ctx context.Context, in *RemoveRoomSessionRequest, opts ...grpc.CallOption) (*RemoveRoomSessionResponse, error)
}

type roomClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomClient(cc grpc.ClientConnInterface) RoomClient {
	return &roomClient{cc}
}

func (c *roomClient) CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateRoomResponse)
	err := c.cc.Invoke(ctx, Room_CreateRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomResponse)
	err := c.cc.Invoke(ctx, Room_GetRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) ListRoomMembers(ctx context.Context, in *ListRoomMemberRequest, opts ...grpc.CallOption) (*ListRoomMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListRoomMemberResponse)
	err := c.cc.Invoke(ctx, Room_ListRoomMembers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) ListRooms(ctx context.Context, in *ListRoomsRequest, opts ...grpc.CallOption) (*ListRoomsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListRoomsResponse)
	err := c.cc.Invoke(ctx, Room_ListRooms_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...grpc.CallOption) (*JoinRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(JoinRoomResponse)
	err := c.cc.Invoke(ctx, Room_JoinRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) GetRoomSession(ctx context.Context, in *GetRoomSessionRequest, opts ...grpc.CallOption) (*GetRoomSessionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomSessionResponse)
	err := c.cc.Invoke(ctx, Room_GetRoomSession_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomClient) RemoveRoomSession(ctx context.Context, in *RemoveRoomSessionRequest, opts ...grpc.CallOption) (*RemoveRoomSessionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveRoomSessionResponse)
	err := c.cc.Invoke(ctx, Room_RemoveRoomSession_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomServer is the server API for Room service.
// All implementations must embed UnimplementedRoomServer
// for forward compatibility.
type RoomServer interface {
	CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error)
	GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error)
	ListRoomMembers(context.Context, *ListRoomMemberRequest) (*ListRoomMemberResponse, error)
	ListRooms(context.Context, *ListRoomsRequest) (*ListRoomsResponse, error)
	JoinRoom(context.Context, *JoinRoomRequest) (*JoinRoomResponse, error)
	GetRoomSession(context.Context, *GetRoomSessionRequest) (*GetRoomSessionResponse, error)
	RemoveRoomSession(context.Context, *RemoveRoomSessionRequest) (*RemoveRoomSessionResponse, error)
	mustEmbedUnimplementedRoomServer()
}

// UnimplementedRoomServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRoomServer struct{}

func (UnimplementedRoomServer) CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedRoomServer) GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoom not implemented")
}
func (UnimplementedRoomServer) ListRoomMembers(context.Context, *ListRoomMemberRequest) (*ListRoomMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRoomMembers not implemented")
}
func (UnimplementedRoomServer) ListRooms(context.Context, *ListRoomsRequest) (*ListRoomsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRooms not implemented")
}
func (UnimplementedRoomServer) JoinRoom(context.Context, *JoinRoomRequest) (*JoinRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinRoom not implemented")
}
func (UnimplementedRoomServer) GetRoomSession(context.Context, *GetRoomSessionRequest) (*GetRoomSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomSession not implemented")
}
func (UnimplementedRoomServer) RemoveRoomSession(context.Context, *RemoveRoomSessionRequest) (*RemoveRoomSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveRoomSession not implemented")
}
func (UnimplementedRoomServer) mustEmbedUnimplementedRoomServer() {}
func (UnimplementedRoomServer) testEmbeddedByValue()              {}

// UnsafeRoomServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomServer will
// result in compilation errors.
type UnsafeRoomServer interface {
	mustEmbedUnimplementedRoomServer()
}

func RegisterRoomServer(s grpc.ServiceRegistrar, srv RoomServer) {
	// If the following call pancis, it indicates UnimplementedRoomServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Room_ServiceDesc, srv)
}

func _Room_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_CreateRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).CreateRoom(ctx, req.(*CreateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_GetRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).GetRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_GetRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).GetRoom(ctx, req.(*GetRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_ListRoomMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRoomMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).ListRoomMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_ListRoomMembers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).ListRoomMembers(ctx, req.(*ListRoomMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_ListRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).ListRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_ListRooms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).ListRooms(ctx, req.(*ListRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_JoinRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).JoinRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_JoinRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).JoinRoom(ctx, req.(*JoinRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_GetRoomSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).GetRoomSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_GetRoomSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).GetRoomSession(ctx, req.(*GetRoomSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Room_RemoveRoomSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRoomSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServer).RemoveRoomSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Room_RemoveRoomSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServer).RemoveRoomSession(ctx, req.(*RemoveRoomSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Room_ServiceDesc is the grpc.ServiceDesc for Room service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Room_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "room.v1.Room",
	HandlerType: (*RoomServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRoom",
			Handler:    _Room_CreateRoom_Handler,
		},
		{
			MethodName: "GetRoom",
			Handler:    _Room_GetRoom_Handler,
		},
		{
			MethodName: "ListRoomMembers",
			Handler:    _Room_ListRoomMembers_Handler,
		},
		{
			MethodName: "ListRooms",
			Handler:    _Room_ListRooms_Handler,
		},
		{
			MethodName: "JoinRoom",
			Handler:    _Room_JoinRoom_Handler,
		},
		{
			MethodName: "GetRoomSession",
			Handler:    _Room_GetRoomSession_Handler,
		},
		{
			MethodName: "RemoveRoomSession",
			Handler:    _Room_RemoveRoomSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "room/service/v1/room.proto",
}
