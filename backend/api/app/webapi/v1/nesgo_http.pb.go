// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.0
// - protoc             v3.12.4
// source: app/webapi/v1/nesgo.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationWebApiAddICECandidate = "/nesgo.webapi.v1.WebApi/AddICECandidate"
const OperationWebApiCreateRoom = "/nesgo.webapi.v1.WebApi/CreateRoom"
const OperationWebApiDeleteMember = "/nesgo.webapi.v1.WebApi/DeleteMember"
const OperationWebApiDeleteRoom = "/nesgo.webapi.v1.WebApi/DeleteRoom"
const OperationWebApiGetRoom = "/nesgo.webapi.v1.WebApi/GetRoom"
const OperationWebApiGetRoomMember = "/nesgo.webapi.v1.WebApi/GetRoomMember"
const OperationWebApiGetRoomSession = "/nesgo.webapi.v1.WebApi/GetRoomSession"
const OperationWebApiGetUser = "/nesgo.webapi.v1.WebApi/GetUser"
const OperationWebApiJoinRoom = "/nesgo.webapi.v1.WebApi/JoinRoom"
const OperationWebApiListAllRooms = "/nesgo.webapi.v1.WebApi/ListAllRooms"
const OperationWebApiListGames = "/nesgo.webapi.v1.WebApi/ListGames"
const OperationWebApiListMembers = "/nesgo.webapi.v1.WebApi/ListMembers"
const OperationWebApiListMyRooms = "/nesgo.webapi.v1.WebApi/ListMyRooms"
const OperationWebApiLogin = "/nesgo.webapi.v1.WebApi/Login"
const OperationWebApiOpenGameConnection = "/nesgo.webapi.v1.WebApi/OpenGameConnection"
const OperationWebApiRegister = "/nesgo.webapi.v1.WebApi/Register"
const OperationWebApiSDPAnswer = "/nesgo.webapi.v1.WebApi/SDPAnswer"
const OperationWebApiSetController = "/nesgo.webapi.v1.WebApi/SetController"
const OperationWebApiUpdateMemberRole = "/nesgo.webapi.v1.WebApi/UpdateMemberRole"
const OperationWebApiUpdateRoom = "/nesgo.webapi.v1.WebApi/UpdateRoom"

type WebApiHTTPServer interface {
	AddICECandidate(context.Context, *AddICECandidateRequest) (*AddICECandidateResponse, error)
	CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error)
	DeleteMember(context.Context, *DeleteMemberRequest) (*DeleteMemberResponse, error)
	DeleteRoom(context.Context, *DeleteRoomRequest) (*DeleteRoomResponse, error)
	GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error)
	GetRoomMember(context.Context, *GetRoomMemberRequest) (*GetRoomMemberResponse, error)
	GetRoomSession(context.Context, *GetRoomSessionRequest) (*GetRoomSessionResponse, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	JoinRoom(context.Context, *JoinRoomRequest) (*JoinRoomResponse, error)
	ListAllRooms(context.Context, *ListRoomRequest) (*ListRoomResponse, error)
	ListGames(context.Context, *ListGamesRequest) (*ListGamesResponse, error)
	ListMembers(context.Context, *ListMemberRequest) (*ListMemberResponse, error)
	ListMyRooms(context.Context, *ListRoomRequest) (*ListRoomResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	OpenGameConnection(context.Context, *OpenGameConnectionRequest) (*OpenGameConnectionResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	SDPAnswer(context.Context, *SDPAnswerRequest) (*SDPAnswerResponse, error)
	SetController(context.Context, *SetControllerRequest) (*SetControllerResponse, error)
	UpdateMemberRole(context.Context, *UpdateMemberRoleRequest) (*UpdateMemberRoleResponse, error)
	UpdateRoom(context.Context, *UpdateRoomRequest) (*UpdateRoomResponse, error)
}

func RegisterWebApiHTTPServer(s *http.Server, srv WebApiHTTPServer) {
	r := s.Route("/")
	r.POST("/api/v1/register", _WebApi_Register0_HTTP_Handler(srv))
	r.POST("/api/v1/login", _WebApi_Login1_HTTP_Handler(srv))
	r.GET("/api/v1/rooms/joined", _WebApi_ListMyRooms0_HTTP_Handler(srv))
	r.GET("/api/v1/rooms", _WebApi_ListAllRooms0_HTTP_Handler(srv))
	r.POST("/api/v1/room", _WebApi_CreateRoom0_HTTP_Handler(srv))
	r.GET("/api/v1/room/{id}", _WebApi_GetRoom0_HTTP_Handler(srv))
	r.GET("/api/v1/user/{id}", _WebApi_GetUser0_HTTP_Handler(srv))
	r.GET("/api/v1/room/session", _WebApi_GetRoomSession0_HTTP_Handler(srv))
	r.POST("/api/v1/game/connection", _WebApi_OpenGameConnection0_HTTP_Handler(srv))
	r.POST("/api/v1/game/sdp", _WebApi_SDPAnswer0_HTTP_Handler(srv))
	r.POST("/api/v1/game/ice", _WebApi_AddICECandidate0_HTTP_Handler(srv))
	r.GET("/api/v1/members", _WebApi_ListMembers0_HTTP_Handler(srv))
	r.POST("/api/v1/room/{roomId}/join", _WebApi_JoinRoom0_HTTP_Handler(srv))
	r.DELETE("/api/v1/room/{roomId}", _WebApi_DeleteRoom0_HTTP_Handler(srv))
	r.PUT("/api/v1/room/{roomId}", _WebApi_UpdateRoom0_HTTP_Handler(srv))
	r.GET("/api/v1/member/{roomId}", _WebApi_GetRoomMember0_HTTP_Handler(srv))
	r.GET("/api/v1/games", _WebApi_ListGames1_HTTP_Handler(srv))
	r.POST("/api/v1/game/controller", _WebApi_SetController0_HTTP_Handler(srv))
	r.PUT("/api/v1/member/role", _WebApi_UpdateMemberRole0_HTTP_Handler(srv))
	r.DELETE("/api/v1/member", _WebApi_DeleteMember0_HTTP_Handler(srv))
}

func _WebApi_Register0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RegisterRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiRegister)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Register(ctx, req.(*RegisterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RegisterResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_Login1_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in LoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*LoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_ListMyRooms0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListRoomRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiListMyRooms)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListMyRooms(ctx, req.(*ListRoomRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListRoomResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_ListAllRooms0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListRoomRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiListAllRooms)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListAllRooms(ctx, req.(*ListRoomRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListRoomResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_CreateRoom0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateRoomRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiCreateRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateRoom(ctx, req.(*CreateRoomRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateRoomResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_GetRoom0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRoomRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiGetRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRoom(ctx, req.(*GetRoomRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetRoomResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_GetUser0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetUserRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiGetUser)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetUser(ctx, req.(*GetUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetUserResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_GetRoomSession0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRoomSessionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiGetRoomSession)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRoomSession(ctx, req.(*GetRoomSessionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetRoomSessionResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_OpenGameConnection0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in OpenGameConnectionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiOpenGameConnection)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.OpenGameConnection(ctx, req.(*OpenGameConnectionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*OpenGameConnectionResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_SDPAnswer0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SDPAnswerRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiSDPAnswer)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SDPAnswer(ctx, req.(*SDPAnswerRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SDPAnswerResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_AddICECandidate0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AddICECandidateRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiAddICECandidate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AddICECandidate(ctx, req.(*AddICECandidateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AddICECandidateResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_ListMembers0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListMemberRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiListMembers)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListMembers(ctx, req.(*ListMemberRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListMemberResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_JoinRoom0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in JoinRoomRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiJoinRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.JoinRoom(ctx, req.(*JoinRoomRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*JoinRoomResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_DeleteRoom0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteRoomRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiDeleteRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteRoom(ctx, req.(*DeleteRoomRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteRoomResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_UpdateRoom0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateRoomRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiUpdateRoom)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateRoom(ctx, req.(*UpdateRoomRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateRoomResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_GetRoomMember0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRoomMemberRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiGetRoomMember)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRoomMember(ctx, req.(*GetRoomMemberRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetRoomMemberResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_ListGames1_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListGamesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiListGames)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListGames(ctx, req.(*ListGamesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListGamesResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_SetController0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SetControllerRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiSetController)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SetController(ctx, req.(*SetControllerRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SetControllerResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_UpdateMemberRole0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateMemberRoleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiUpdateMemberRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateMemberRole(ctx, req.(*UpdateMemberRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateMemberRoleResponse)
		return ctx.Result(200, reply)
	}
}

func _WebApi_DeleteMember0_HTTP_Handler(srv WebApiHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteMemberRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationWebApiDeleteMember)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteMember(ctx, req.(*DeleteMemberRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteMemberResponse)
		return ctx.Result(200, reply)
	}
}

type WebApiHTTPClient interface {
	AddICECandidate(ctx context.Context, req *AddICECandidateRequest, opts ...http.CallOption) (rsp *AddICECandidateResponse, err error)
	CreateRoom(ctx context.Context, req *CreateRoomRequest, opts ...http.CallOption) (rsp *CreateRoomResponse, err error)
	DeleteMember(ctx context.Context, req *DeleteMemberRequest, opts ...http.CallOption) (rsp *DeleteMemberResponse, err error)
	DeleteRoom(ctx context.Context, req *DeleteRoomRequest, opts ...http.CallOption) (rsp *DeleteRoomResponse, err error)
	GetRoom(ctx context.Context, req *GetRoomRequest, opts ...http.CallOption) (rsp *GetRoomResponse, err error)
	GetRoomMember(ctx context.Context, req *GetRoomMemberRequest, opts ...http.CallOption) (rsp *GetRoomMemberResponse, err error)
	GetRoomSession(ctx context.Context, req *GetRoomSessionRequest, opts ...http.CallOption) (rsp *GetRoomSessionResponse, err error)
	GetUser(ctx context.Context, req *GetUserRequest, opts ...http.CallOption) (rsp *GetUserResponse, err error)
	JoinRoom(ctx context.Context, req *JoinRoomRequest, opts ...http.CallOption) (rsp *JoinRoomResponse, err error)
	ListAllRooms(ctx context.Context, req *ListRoomRequest, opts ...http.CallOption) (rsp *ListRoomResponse, err error)
	ListGames(ctx context.Context, req *ListGamesRequest, opts ...http.CallOption) (rsp *ListGamesResponse, err error)
	ListMembers(ctx context.Context, req *ListMemberRequest, opts ...http.CallOption) (rsp *ListMemberResponse, err error)
	ListMyRooms(ctx context.Context, req *ListRoomRequest, opts ...http.CallOption) (rsp *ListRoomResponse, err error)
	Login(ctx context.Context, req *LoginRequest, opts ...http.CallOption) (rsp *LoginResponse, err error)
	OpenGameConnection(ctx context.Context, req *OpenGameConnectionRequest, opts ...http.CallOption) (rsp *OpenGameConnectionResponse, err error)
	Register(ctx context.Context, req *RegisterRequest, opts ...http.CallOption) (rsp *RegisterResponse, err error)
	SDPAnswer(ctx context.Context, req *SDPAnswerRequest, opts ...http.CallOption) (rsp *SDPAnswerResponse, err error)
	SetController(ctx context.Context, req *SetControllerRequest, opts ...http.CallOption) (rsp *SetControllerResponse, err error)
	UpdateMemberRole(ctx context.Context, req *UpdateMemberRoleRequest, opts ...http.CallOption) (rsp *UpdateMemberRoleResponse, err error)
	UpdateRoom(ctx context.Context, req *UpdateRoomRequest, opts ...http.CallOption) (rsp *UpdateRoomResponse, err error)
}

type WebApiHTTPClientImpl struct {
	cc *http.Client
}

func NewWebApiHTTPClient(client *http.Client) WebApiHTTPClient {
	return &WebApiHTTPClientImpl{client}
}

func (c *WebApiHTTPClientImpl) AddICECandidate(ctx context.Context, in *AddICECandidateRequest, opts ...http.CallOption) (*AddICECandidateResponse, error) {
	var out AddICECandidateResponse
	pattern := "/api/v1/game/ice"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiAddICECandidate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...http.CallOption) (*CreateRoomResponse, error) {
	var out CreateRoomResponse
	pattern := "/api/v1/room"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiCreateRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) DeleteMember(ctx context.Context, in *DeleteMemberRequest, opts ...http.CallOption) (*DeleteMemberResponse, error) {
	var out DeleteMemberResponse
	pattern := "/api/v1/member"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationWebApiDeleteMember))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) DeleteRoom(ctx context.Context, in *DeleteRoomRequest, opts ...http.CallOption) (*DeleteRoomResponse, error) {
	var out DeleteRoomResponse
	pattern := "/api/v1/room/{roomId}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationWebApiDeleteRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) GetRoom(ctx context.Context, in *GetRoomRequest, opts ...http.CallOption) (*GetRoomResponse, error) {
	var out GetRoomResponse
	pattern := "/api/v1/room/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationWebApiGetRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) GetRoomMember(ctx context.Context, in *GetRoomMemberRequest, opts ...http.CallOption) (*GetRoomMemberResponse, error) {
	var out GetRoomMemberResponse
	pattern := "/api/v1/member/{roomId}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationWebApiGetRoomMember))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) GetRoomSession(ctx context.Context, in *GetRoomSessionRequest, opts ...http.CallOption) (*GetRoomSessionResponse, error) {
	var out GetRoomSessionResponse
	pattern := "/api/v1/room/session"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiGetRoomSession))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) GetUser(ctx context.Context, in *GetUserRequest, opts ...http.CallOption) (*GetUserResponse, error) {
	var out GetUserResponse
	pattern := "/api/v1/user/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationWebApiGetUser))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...http.CallOption) (*JoinRoomResponse, error) {
	var out JoinRoomResponse
	pattern := "/api/v1/room/{roomId}/join"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiJoinRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) ListAllRooms(ctx context.Context, in *ListRoomRequest, opts ...http.CallOption) (*ListRoomResponse, error) {
	var out ListRoomResponse
	pattern := "/api/v1/rooms"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationWebApiListAllRooms))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) ListGames(ctx context.Context, in *ListGamesRequest, opts ...http.CallOption) (*ListGamesResponse, error) {
	var out ListGamesResponse
	pattern := "/api/v1/games"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationWebApiListGames))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) ListMembers(ctx context.Context, in *ListMemberRequest, opts ...http.CallOption) (*ListMemberResponse, error) {
	var out ListMemberResponse
	pattern := "/api/v1/members"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationWebApiListMembers))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) ListMyRooms(ctx context.Context, in *ListRoomRequest, opts ...http.CallOption) (*ListRoomResponse, error) {
	var out ListRoomResponse
	pattern := "/api/v1/rooms/joined"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationWebApiListMyRooms))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) Login(ctx context.Context, in *LoginRequest, opts ...http.CallOption) (*LoginResponse, error) {
	var out LoginResponse
	pattern := "/api/v1/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) OpenGameConnection(ctx context.Context, in *OpenGameConnectionRequest, opts ...http.CallOption) (*OpenGameConnectionResponse, error) {
	var out OpenGameConnectionResponse
	pattern := "/api/v1/game/connection"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiOpenGameConnection))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) Register(ctx context.Context, in *RegisterRequest, opts ...http.CallOption) (*RegisterResponse, error) {
	var out RegisterResponse
	pattern := "/api/v1/register"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiRegister))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) SDPAnswer(ctx context.Context, in *SDPAnswerRequest, opts ...http.CallOption) (*SDPAnswerResponse, error) {
	var out SDPAnswerResponse
	pattern := "/api/v1/game/sdp"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiSDPAnswer))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) SetController(ctx context.Context, in *SetControllerRequest, opts ...http.CallOption) (*SetControllerResponse, error) {
	var out SetControllerResponse
	pattern := "/api/v1/game/controller"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiSetController))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) UpdateMemberRole(ctx context.Context, in *UpdateMemberRoleRequest, opts ...http.CallOption) (*UpdateMemberRoleResponse, error) {
	var out UpdateMemberRoleResponse
	pattern := "/api/v1/member/role"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiUpdateMemberRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *WebApiHTTPClientImpl) UpdateRoom(ctx context.Context, in *UpdateRoomRequest, opts ...http.CallOption) (*UpdateRoomResponse, error) {
	var out UpdateRoomResponse
	pattern := "/api/v1/room/{roomId}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationWebApiUpdateRoom))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
