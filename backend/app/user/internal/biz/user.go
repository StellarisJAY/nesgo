package biz

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
	"gorm.io/gorm"
)

type User struct {
	ID       int64  `json:"ID"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserRepo interface {
	GetUser(ctx context.Context, id int64) (*User, error)
	GetUserByName(ctx context.Context, name string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
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
	u, err := uc.ur.GetUser(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, v1.ErrorUserNotFound("user not found")
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UserUseCase) GetUserByName(ctx context.Context, name string) (*User, error) {
	u, err := uc.ur.GetUserByName(ctx, name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, v1.ErrorUserNotFound("user not found")
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user *User) error {
	sum := md5.New().Sum([]byte(user.Password))
	user.Password = hex.EncodeToString(sum[:])
	err := uc.ur.CreateUser(ctx, user)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return v1.ErrorUsernameConflict("username already exists")
	}
	return err
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user *User) error {
	panic("implement me")
}

func (uc *UserUseCase) VerifyPassword(ctx context.Context, name string, password string) error {
	user, err := uc.ur.GetUserByName(ctx, name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return v1.ErrorUserNotFound("user not found")
	}
	if err != nil {
		return err
	}
	sum := md5.New().Sum([]byte(password))
	if user.Password != hex.EncodeToString(sum[:]) {
		return v1.ErrorVerifyPasswordFailed("password error")
	}
	return nil
}
