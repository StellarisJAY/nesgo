package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/user/internal/biz"
)

type UserService struct {
	v1.UnimplementedUserServer
	uc     *biz.UserUseCase
	kb     *biz.UserKeyboardBindingUseCase
	mu     *biz.MacroUseCase
	logger *log.Helper
}

func NewUserService(uc *biz.UserUseCase, kb *biz.UserKeyboardBindingUseCase, mu *biz.MacroUseCase, logger log.Logger) *UserService {
	return &UserService{
		uc:     uc,
		kb:     kb,
		mu:     mu,
		logger: log.NewHelper(log.With(logger, "module", "service/user")),
	}
}

func (u *UserService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	user := &biz.User{
		ID:       0,
		Name:     request.Name,
		Password: request.Password,
	}
	err := u.uc.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &v1.CreateUserResponse{
		Id:   user.ID,
		Name: user.Name,
	}, nil
}

func (u *UserService) GetUser(ctx context.Context, request *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	user, err := u.uc.GetUser(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserResponse{
		Id:   user.ID,
		Name: user.Name,
	}, nil
}

func (u *UserService) GetUserByName(ctx context.Context, request *v1.GetUserByNameRequest) (*v1.GetUserByNameResponse, error) {
	user, err := u.uc.GetUserByName(ctx, request.Name)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserByNameResponse{
		Id:   user.ID,
		Name: user.Name,
	}, nil
}

func (u *UserService) UpdateUser(ctx context.Context, request *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) VerifyPassword(ctx context.Context, request *v1.VerifyPasswordRequest) (*v1.VerifyPasswordResponse, error) {
	err := u.uc.VerifyPassword(ctx, request.Name, request.Password)
	if err != nil {
		return nil, err
	}
	return &v1.VerifyPasswordResponse{}, nil
}
