package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	userAPI "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/biz"
)

type userRepo struct {
	data   *Data
	logger *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/user")),
	}
}

func (u *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	user, err := u.data.userCli.GetUser(ctx, &userAPI.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &biz.User{Id: id, Name: user.Name}, nil
}
