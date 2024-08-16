// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: user/service/v1/user.proto

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
	User_CreateUser_FullMethodName                = "/user.v1.User/CreateUser"
	User_GetUser_FullMethodName                   = "/user.v1.User/GetUser"
	User_GetUserByName_FullMethodName             = "/user.v1.User/GetUserByName"
	User_UpdateUser_FullMethodName                = "/user.v1.User/UpdateUser"
	User_VerifyPassword_FullMethodName            = "/user.v1.User/VerifyPassword"
	User_CreateUserKeyboardBinding_FullMethodName = "/user.v1.User/CreateUserKeyboardBinding"
	User_ListUserKeyboardBinding_FullMethodName   = "/user.v1.User/ListUserKeyboardBinding"
	User_GetUserKeyboardBinding_FullMethodName    = "/user.v1.User/GetUserKeyboardBinding"
	User_UpdateUserKeyboardBinding_FullMethodName = "/user.v1.User/UpdateUserKeyboardBinding"
	User_DeleteUserKeyboardBinding_FullMethodName = "/user.v1.User/DeleteUserKeyboardBinding"
)

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	GetUserByName(ctx context.Context, in *GetUserByNameRequest, opts ...grpc.CallOption) (*GetUserByNameResponse, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error)
	VerifyPassword(ctx context.Context, in *VerifyPasswordRequest, opts ...grpc.CallOption) (*VerifyPasswordResponse, error)
	CreateUserKeyboardBinding(ctx context.Context, in *CreateUserKeyboardBindingRequest, opts ...grpc.CallOption) (*CreateUserKeyboardBindingResponse, error)
	ListUserKeyboardBinding(ctx context.Context, in *ListUserKeyboardBindingRequest, opts ...grpc.CallOption) (*ListUserKeyboardBindingResponse, error)
	GetUserKeyboardBinding(ctx context.Context, in *GetUserKeyboardBindingRequest, opts ...grpc.CallOption) (*GetUserKeyboardBindingResponse, error)
	UpdateUserKeyboardBinding(ctx context.Context, in *UpdateUserKeyboardBindingRequest, opts ...grpc.CallOption) (*UpdateUserKeyboardBindingResponse, error)
	DeleteUserKeyboardBinding(ctx context.Context, in *DeleteUserKeyboardBindingRequest, opts ...grpc.CallOption) (*DeleteUserKeyboardBindingResponse, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, User_CreateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, User_GetUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetUserByName(ctx context.Context, in *GetUserByNameRequest, opts ...grpc.CallOption) (*GetUserByNameResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserByNameResponse)
	err := c.cc.Invoke(ctx, User_GetUserByName_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateUserResponse)
	err := c.cc.Invoke(ctx, User_UpdateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) VerifyPassword(ctx context.Context, in *VerifyPasswordRequest, opts ...grpc.CallOption) (*VerifyPasswordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VerifyPasswordResponse)
	err := c.cc.Invoke(ctx, User_VerifyPassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) CreateUserKeyboardBinding(ctx context.Context, in *CreateUserKeyboardBindingRequest, opts ...grpc.CallOption) (*CreateUserKeyboardBindingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateUserKeyboardBindingResponse)
	err := c.cc.Invoke(ctx, User_CreateUserKeyboardBinding_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) ListUserKeyboardBinding(ctx context.Context, in *ListUserKeyboardBindingRequest, opts ...grpc.CallOption) (*ListUserKeyboardBindingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListUserKeyboardBindingResponse)
	err := c.cc.Invoke(ctx, User_ListUserKeyboardBinding_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetUserKeyboardBinding(ctx context.Context, in *GetUserKeyboardBindingRequest, opts ...grpc.CallOption) (*GetUserKeyboardBindingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserKeyboardBindingResponse)
	err := c.cc.Invoke(ctx, User_GetUserKeyboardBinding_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) UpdateUserKeyboardBinding(ctx context.Context, in *UpdateUserKeyboardBindingRequest, opts ...grpc.CallOption) (*UpdateUserKeyboardBindingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateUserKeyboardBindingResponse)
	err := c.cc.Invoke(ctx, User_UpdateUserKeyboardBinding_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) DeleteUserKeyboardBinding(ctx context.Context, in *DeleteUserKeyboardBindingRequest, opts ...grpc.CallOption) (*DeleteUserKeyboardBindingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteUserKeyboardBindingResponse)
	err := c.cc.Invoke(ctx, User_DeleteUserKeyboardBinding_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility.
type UserServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	GetUserByName(context.Context, *GetUserByNameRequest) (*GetUserByNameResponse, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserResponse, error)
	VerifyPassword(context.Context, *VerifyPasswordRequest) (*VerifyPasswordResponse, error)
	CreateUserKeyboardBinding(context.Context, *CreateUserKeyboardBindingRequest) (*CreateUserKeyboardBindingResponse, error)
	ListUserKeyboardBinding(context.Context, *ListUserKeyboardBindingRequest) (*ListUserKeyboardBindingResponse, error)
	GetUserKeyboardBinding(context.Context, *GetUserKeyboardBindingRequest) (*GetUserKeyboardBindingResponse, error)
	UpdateUserKeyboardBinding(context.Context, *UpdateUserKeyboardBindingRequest) (*UpdateUserKeyboardBindingResponse, error)
	DeleteUserKeyboardBinding(context.Context, *DeleteUserKeyboardBindingRequest) (*DeleteUserKeyboardBindingResponse, error)
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserServer struct{}

func (UnimplementedUserServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserServer) GetUserByName(context.Context, *GetUserByNameRequest) (*GetUserByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByName not implemented")
}
func (UnimplementedUserServer) UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUserServer) VerifyPassword(context.Context, *VerifyPasswordRequest) (*VerifyPasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyPassword not implemented")
}
func (UnimplementedUserServer) CreateUserKeyboardBinding(context.Context, *CreateUserKeyboardBindingRequest) (*CreateUserKeyboardBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserKeyboardBinding not implemented")
}
func (UnimplementedUserServer) ListUserKeyboardBinding(context.Context, *ListUserKeyboardBindingRequest) (*ListUserKeyboardBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUserKeyboardBinding not implemented")
}
func (UnimplementedUserServer) GetUserKeyboardBinding(context.Context, *GetUserKeyboardBindingRequest) (*GetUserKeyboardBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserKeyboardBinding not implemented")
}
func (UnimplementedUserServer) UpdateUserKeyboardBinding(context.Context, *UpdateUserKeyboardBindingRequest) (*UpdateUserKeyboardBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserKeyboardBinding not implemented")
}
func (UnimplementedUserServer) DeleteUserKeyboardBinding(context.Context, *DeleteUserKeyboardBindingRequest) (*DeleteUserKeyboardBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserKeyboardBinding not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}
func (UnimplementedUserServer) testEmbeddedByValue()              {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	// If the following call pancis, it indicates UnimplementedUserServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetUserByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUserByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetUserByName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUserByName(ctx, req.(*GetUserByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_UpdateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_VerifyPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).VerifyPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_VerifyPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).VerifyPassword(ctx, req.(*VerifyPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_CreateUserKeyboardBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserKeyboardBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).CreateUserKeyboardBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_CreateUserKeyboardBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).CreateUserKeyboardBinding(ctx, req.(*CreateUserKeyboardBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_ListUserKeyboardBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserKeyboardBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).ListUserKeyboardBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_ListUserKeyboardBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).ListUserKeyboardBinding(ctx, req.(*ListUserKeyboardBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetUserKeyboardBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserKeyboardBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUserKeyboardBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetUserKeyboardBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUserKeyboardBinding(ctx, req.(*GetUserKeyboardBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_UpdateUserKeyboardBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserKeyboardBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).UpdateUserKeyboardBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_UpdateUserKeyboardBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).UpdateUserKeyboardBinding(ctx, req.(*UpdateUserKeyboardBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_DeleteUserKeyboardBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserKeyboardBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).DeleteUserKeyboardBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_DeleteUserKeyboardBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).DeleteUserKeyboardBinding(ctx, req.(*DeleteUserKeyboardBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _User_CreateUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _User_GetUser_Handler,
		},
		{
			MethodName: "GetUserByName",
			Handler:    _User_GetUserByName_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _User_UpdateUser_Handler,
		},
		{
			MethodName: "VerifyPassword",
			Handler:    _User_VerifyPassword_Handler,
		},
		{
			MethodName: "CreateUserKeyboardBinding",
			Handler:    _User_CreateUserKeyboardBinding_Handler,
		},
		{
			MethodName: "ListUserKeyboardBinding",
			Handler:    _User_ListUserKeyboardBinding_Handler,
		},
		{
			MethodName: "GetUserKeyboardBinding",
			Handler:    _User_GetUserKeyboardBinding_Handler,
		},
		{
			MethodName: "UpdateUserKeyboardBinding",
			Handler:    _User_UpdateUserKeyboardBinding_Handler,
		},
		{
			MethodName: "DeleteUserKeyboardBinding",
			Handler:    _User_DeleteUserKeyboardBinding_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user/service/v1/user.proto",
}
