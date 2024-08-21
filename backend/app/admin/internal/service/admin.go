package service

import (
	"context"
	adminAPI "github.com/stellarisJAY/nesgo/backend/api/app/admin/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/biz"
)

func (s *AdminService) Login(ctx context.Context, request *adminAPI.LoginRequest) (*adminAPI.LoginResponse, error) {
	token, err := s.a.Login(ctx, request.Name, request.Password)
	if err != nil {
		return nil, err
	}
	return &adminAPI.LoginResponse{Token: token}, nil
}

func (s *AdminService) CreateAdmin(ctx context.Context, request *adminAPI.CreateAdminRequest) (*adminAPI.CreateAdminResponse, error) {
	err := s.a.CreateAdmin(ctx, &biz.Admin{Name: request.Name, Password: request.Password})
	if err != nil {
		return nil, err
	}
	return &adminAPI.CreateAdminResponse{}, nil
}
