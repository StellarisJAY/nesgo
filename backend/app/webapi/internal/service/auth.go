package service

import (
	"context"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
)

func (ws *WebApiService) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	err := ws.ac.Register(ctx, request.Name, request.Password)
	if err != nil {
		return nil, err
	}
	return &v1.RegisterResponse{}, nil
}

func (ws *WebApiService) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	token, err := ws.ac.Login(ctx, request.Name, request.Password)
	if err != nil {
		return nil, err
	}
	return &v1.LoginResponse{Token: token}, nil
}
