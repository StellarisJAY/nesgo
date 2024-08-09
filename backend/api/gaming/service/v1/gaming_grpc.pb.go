// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: gaming/service/v1/gaming.proto

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
	Gaming_CreateGameInstance_FullMethodName = "/gaming.v1.Gaming/CreateGameInstance"
	Gaming_OpenGameConnection_FullMethodName = "/gaming.v1.Gaming/OpenGameConnection"
	Gaming_SDPAnswer_FullMethodName          = "/gaming.v1.Gaming/SDPAnswer"
	Gaming_ICECandidate_FullMethodName       = "/gaming.v1.Gaming/ICECandidate"
	Gaming_PauseEmulator_FullMethodName      = "/gaming.v1.Gaming/PauseEmulator"
	Gaming_RestartEmulator_FullMethodName    = "/gaming.v1.Gaming/RestartEmulator"
	Gaming_DeleteGameInstance_FullMethodName = "/gaming.v1.Gaming/DeleteGameInstance"
	Gaming_UploadGame_FullMethodName         = "/gaming.v1.Gaming/UploadGame"
	Gaming_ListGames_FullMethodName          = "/gaming.v1.Gaming/ListGames"
	Gaming_DeleteGameFile_FullMethodName     = "/gaming.v1.Gaming/DeleteGameFile"
)

// GamingClient is the client API for Gaming service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GamingClient interface {
	CreateGameInstance(ctx context.Context, in *CreateGameInstanceRequest, opts ...grpc.CallOption) (*CreateGameInstanceResponse, error)
	OpenGameConnection(ctx context.Context, in *OpenGameConnectionRequest, opts ...grpc.CallOption) (*OpenGameConnectionResponse, error)
	SDPAnswer(ctx context.Context, in *SDPAnswerRequest, opts ...grpc.CallOption) (*SDPAnswerResponse, error)
	ICECandidate(ctx context.Context, in *ICECandidateRequest, opts ...grpc.CallOption) (*ICECandidateResponse, error)
	PauseEmulator(ctx context.Context, in *PauseEmulatorRequest, opts ...grpc.CallOption) (*PauseEmulatorResponse, error)
	RestartEmulator(ctx context.Context, in *RestartEmulatorRequest, opts ...grpc.CallOption) (*RestartEmulatorResponse, error)
	DeleteGameInstance(ctx context.Context, in *DeleteGameInstanceRequest, opts ...grpc.CallOption) (*DeleteGameInstanceResponse, error)
	UploadGame(ctx context.Context, in *UploadGameRequest, opts ...grpc.CallOption) (*UploadGameResponse, error)
	ListGames(ctx context.Context, in *ListGamesRequest, opts ...grpc.CallOption) (*ListGamesResponse, error)
	DeleteGameFile(ctx context.Context, in *DeleteGameFileRequest, opts ...grpc.CallOption) (*DeleteGameFileResponse, error)
}

type gamingClient struct {
	cc grpc.ClientConnInterface
}

func NewGamingClient(cc grpc.ClientConnInterface) GamingClient {
	return &gamingClient{cc}
}

func (c *gamingClient) CreateGameInstance(ctx context.Context, in *CreateGameInstanceRequest, opts ...grpc.CallOption) (*CreateGameInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateGameInstanceResponse)
	err := c.cc.Invoke(ctx, Gaming_CreateGameInstance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamingClient) OpenGameConnection(ctx context.Context, in *OpenGameConnectionRequest, opts ...grpc.CallOption) (*OpenGameConnectionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpenGameConnectionResponse)
	err := c.cc.Invoke(ctx, Gaming_OpenGameConnection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamingClient) SDPAnswer(ctx context.Context, in *SDPAnswerRequest, opts ...grpc.CallOption) (*SDPAnswerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SDPAnswerResponse)
	err := c.cc.Invoke(ctx, Gaming_SDPAnswer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamingClient) ICECandidate(ctx context.Context, in *ICECandidateRequest, opts ...grpc.CallOption) (*ICECandidateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ICECandidateResponse)
	err := c.cc.Invoke(ctx, Gaming_ICECandidate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamingClient) PauseEmulator(ctx context.Context, in *PauseEmulatorRequest, opts ...grpc.CallOption) (*PauseEmulatorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PauseEmulatorResponse)
	err := c.cc.Invoke(ctx, Gaming_PauseEmulator_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamingClient) RestartEmulator(ctx context.Context, in *RestartEmulatorRequest, opts ...grpc.CallOption) (*RestartEmulatorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RestartEmulatorResponse)
	err := c.cc.Invoke(ctx, Gaming_RestartEmulator_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamingClient) DeleteGameInstance(ctx context.Context, in *DeleteGameInstanceRequest, opts ...grpc.CallOption) (*DeleteGameInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteGameInstanceResponse)
	err := c.cc.Invoke(ctx, Gaming_DeleteGameInstance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamingClient) UploadGame(ctx context.Context, in *UploadGameRequest, opts ...grpc.CallOption) (*UploadGameResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadGameResponse)
	err := c.cc.Invoke(ctx, Gaming_UploadGame_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamingClient) ListGames(ctx context.Context, in *ListGamesRequest, opts ...grpc.CallOption) (*ListGamesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListGamesResponse)
	err := c.cc.Invoke(ctx, Gaming_ListGames_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamingClient) DeleteGameFile(ctx context.Context, in *DeleteGameFileRequest, opts ...grpc.CallOption) (*DeleteGameFileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteGameFileResponse)
	err := c.cc.Invoke(ctx, Gaming_DeleteGameFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GamingServer is the server API for Gaming service.
// All implementations must embed UnimplementedGamingServer
// for forward compatibility.
type GamingServer interface {
	CreateGameInstance(context.Context, *CreateGameInstanceRequest) (*CreateGameInstanceResponse, error)
	OpenGameConnection(context.Context, *OpenGameConnectionRequest) (*OpenGameConnectionResponse, error)
	SDPAnswer(context.Context, *SDPAnswerRequest) (*SDPAnswerResponse, error)
	ICECandidate(context.Context, *ICECandidateRequest) (*ICECandidateResponse, error)
	PauseEmulator(context.Context, *PauseEmulatorRequest) (*PauseEmulatorResponse, error)
	RestartEmulator(context.Context, *RestartEmulatorRequest) (*RestartEmulatorResponse, error)
	DeleteGameInstance(context.Context, *DeleteGameInstanceRequest) (*DeleteGameInstanceResponse, error)
	UploadGame(context.Context, *UploadGameRequest) (*UploadGameResponse, error)
	ListGames(context.Context, *ListGamesRequest) (*ListGamesResponse, error)
	DeleteGameFile(context.Context, *DeleteGameFileRequest) (*DeleteGameFileResponse, error)
	mustEmbedUnimplementedGamingServer()
}

// UnimplementedGamingServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGamingServer struct{}

func (UnimplementedGamingServer) CreateGameInstance(context.Context, *CreateGameInstanceRequest) (*CreateGameInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGameInstance not implemented")
}
func (UnimplementedGamingServer) OpenGameConnection(context.Context, *OpenGameConnectionRequest) (*OpenGameConnectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenGameConnection not implemented")
}
func (UnimplementedGamingServer) SDPAnswer(context.Context, *SDPAnswerRequest) (*SDPAnswerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SDPAnswer not implemented")
}
func (UnimplementedGamingServer) ICECandidate(context.Context, *ICECandidateRequest) (*ICECandidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ICECandidate not implemented")
}
func (UnimplementedGamingServer) PauseEmulator(context.Context, *PauseEmulatorRequest) (*PauseEmulatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PauseEmulator not implemented")
}
func (UnimplementedGamingServer) RestartEmulator(context.Context, *RestartEmulatorRequest) (*RestartEmulatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestartEmulator not implemented")
}
func (UnimplementedGamingServer) DeleteGameInstance(context.Context, *DeleteGameInstanceRequest) (*DeleteGameInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGameInstance not implemented")
}
func (UnimplementedGamingServer) UploadGame(context.Context, *UploadGameRequest) (*UploadGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadGame not implemented")
}
func (UnimplementedGamingServer) ListGames(context.Context, *ListGamesRequest) (*ListGamesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGames not implemented")
}
func (UnimplementedGamingServer) DeleteGameFile(context.Context, *DeleteGameFileRequest) (*DeleteGameFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGameFile not implemented")
}
func (UnimplementedGamingServer) mustEmbedUnimplementedGamingServer() {}
func (UnimplementedGamingServer) testEmbeddedByValue()                {}

// UnsafeGamingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GamingServer will
// result in compilation errors.
type UnsafeGamingServer interface {
	mustEmbedUnimplementedGamingServer()
}

func RegisterGamingServer(s grpc.ServiceRegistrar, srv GamingServer) {
	// If the following call pancis, it indicates UnimplementedGamingServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Gaming_ServiceDesc, srv)
}

func _Gaming_CreateGameInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGameInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamingServer).CreateGameInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gaming_CreateGameInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamingServer).CreateGameInstance(ctx, req.(*CreateGameInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gaming_OpenGameConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenGameConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamingServer).OpenGameConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gaming_OpenGameConnection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamingServer).OpenGameConnection(ctx, req.(*OpenGameConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gaming_SDPAnswer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SDPAnswerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamingServer).SDPAnswer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gaming_SDPAnswer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamingServer).SDPAnswer(ctx, req.(*SDPAnswerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gaming_ICECandidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ICECandidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamingServer).ICECandidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gaming_ICECandidate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamingServer).ICECandidate(ctx, req.(*ICECandidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gaming_PauseEmulator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PauseEmulatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamingServer).PauseEmulator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gaming_PauseEmulator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamingServer).PauseEmulator(ctx, req.(*PauseEmulatorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gaming_RestartEmulator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestartEmulatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamingServer).RestartEmulator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gaming_RestartEmulator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamingServer).RestartEmulator(ctx, req.(*RestartEmulatorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gaming_DeleteGameInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteGameInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamingServer).DeleteGameInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gaming_DeleteGameInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamingServer).DeleteGameInstance(ctx, req.(*DeleteGameInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gaming_UploadGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamingServer).UploadGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gaming_UploadGame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamingServer).UploadGame(ctx, req.(*UploadGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gaming_ListGames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListGamesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamingServer).ListGames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gaming_ListGames_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamingServer).ListGames(ctx, req.(*ListGamesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gaming_DeleteGameFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteGameFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamingServer).DeleteGameFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gaming_DeleteGameFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamingServer).DeleteGameFile(ctx, req.(*DeleteGameFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Gaming_ServiceDesc is the grpc.ServiceDesc for Gaming service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gaming_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gaming.v1.Gaming",
	HandlerType: (*GamingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGameInstance",
			Handler:    _Gaming_CreateGameInstance_Handler,
		},
		{
			MethodName: "OpenGameConnection",
			Handler:    _Gaming_OpenGameConnection_Handler,
		},
		{
			MethodName: "SDPAnswer",
			Handler:    _Gaming_SDPAnswer_Handler,
		},
		{
			MethodName: "ICECandidate",
			Handler:    _Gaming_ICECandidate_Handler,
		},
		{
			MethodName: "PauseEmulator",
			Handler:    _Gaming_PauseEmulator_Handler,
		},
		{
			MethodName: "RestartEmulator",
			Handler:    _Gaming_RestartEmulator_Handler,
		},
		{
			MethodName: "DeleteGameInstance",
			Handler:    _Gaming_DeleteGameInstance_Handler,
		},
		{
			MethodName: "UploadGame",
			Handler:    _Gaming_UploadGame_Handler,
		},
		{
			MethodName: "ListGames",
			Handler:    _Gaming_ListGames_Handler,
		},
		{
			MethodName: "DeleteGameFile",
			Handler:    _Gaming_DeleteGameFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gaming/service/v1/gaming.proto",
}
