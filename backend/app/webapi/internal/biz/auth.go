package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
)

type AuthUseCase struct {
	ur     UserRepo
	logger *log.Helper
}

func NewAuthUseCase(ur UserRepo, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{
		ur:     ur,
		logger: log.NewHelper(log.With(logger, "module", "biz/auth")),
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, username, password string) (*User, error) {
	panic("implement me")
}

func (uc *AuthUseCase) Register(ctx context.Context, username, password string) error {
	user, err := uc.ur.GetUserByName(ctx, username)
	if user != nil {
		return v1.ErrorUsernameConflict("username already exists")
	}
	_, err = uc.ur.CreateUser(ctx, &User{
		ID:       0,
		Name:     username,
		Password: password,
	})
	if err != nil {
		return v1.ErrorRegisterFailed("create user failed", err)
	}
	return nil
}
