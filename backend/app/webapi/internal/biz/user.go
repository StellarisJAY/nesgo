package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID       int64
	Name     string
	Password string
}

type UserRepo interface {
	CreateUser(ctx context.Context, u *User) (*User, error)
	GetUser(ctx context.Context, id int64) (*User, error)
	GetUserByName(ctx context.Context, name string) (*User, error)
	VerifyPassword(ctx context.Context, name string, password string) error
}

type UserUseCase struct {
	ur     UserRepo
	logger *log.Helper
}

func NewUserUseCase(ur UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		ur:     ur,
		logger: log.NewHelper(log.With(logger, "module", "biz/user")),
	}
}

func (uc *UserUseCase) GetUser(ctx context.Context, id int64) (*User, error) {
	user, err := uc.ur.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
