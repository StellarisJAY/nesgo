package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	userAPI "github.com/stellarisjay/nesgo/backend/api/user/service/v1"
	"github.com/stellarisjay/nesgo/backend/app/webapi/internal/biz"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/user")),
	}
}

func (ur *userRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	resp, err := ur.data.uc.CreateUser(ctx, &userAPI.CreateUserRequest{
		Name:     u.Name,
		Password: u.Password,
	})
	if err != nil {
		return nil, err
	}
	u.ID = resp.Id
	return u, nil
}

func (ur *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	resp, err := ur.data.uc.GetUser(ctx, &userAPI.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		ID:   resp.Id,
		Name: resp.Name,
	}, nil
}

func (ur *userRepo) VerifyPassword(ctx context.Context, name string, password string) error {
	_, err := ur.data.uc.VerifyPassword(ctx, &userAPI.VerifyPasswordRequest{Name: name, Password: password})
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepo) GetUserByName(ctx context.Context, name string) (*biz.User, error) {
	user, err := ur.data.uc.GetUserByName(ctx, &userAPI.GetUserByNameRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		ID:   user.Id,
		Name: user.Name,
	}, nil
}
