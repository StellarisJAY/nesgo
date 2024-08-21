package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/user/internal/biz"
	"github.com/stellarisJAY/nesgo/backend/pkg/cache"
	"gorm.io/gorm"
	"time"
)

type userRepo struct {
	data   *Data
	logger *log.Helper
}

type User struct {
	Id        int64  `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name      string `gorm:"size:255" json:"name"`
	Password  string `gorm:"size:255" json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/user")),
	}
}

func userCacheKey(id int64) string {
	return fmt.Sprintf("/nesgo/user/%d", id)
}

func (u *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	target := u.cacheGetUser(ctx, id)
	if target == nil {
		user := &User{}
		err := u.data.db.Model(user).Where("id=?", id).First(user).Error
		if err != nil {
			return nil, err
		}
		u.cacheSetUser(ctx, user)
		return &biz.User{
			ID:       user.Id,
			Name:     user.Name,
			Password: user.Password,
		}, nil
	}
	return &biz.User{
		ID:       target.Id,
		Name:     target.Name,
		Password: target.Password,
	}, nil
}

func (u *userRepo) CreateUser(ctx context.Context, user *biz.User) error {
	model := &User{
		Id:        u.data.snowflake.Generate().Int64(),
		Name:      user.Name,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := u.data.db.Model(model).WithContext(ctx).Create(model).Error
	return err
}

func (u *userRepo) UpdateUser(ctx context.Context, user *biz.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) GetUserByName(ctx context.Context, name string) (*biz.User, error) {
	user := &User{}
	err := u.data.db.Model(user).Where("name=?", name).First(user).Error
	if err != nil {
		return nil, err
	}
	return &biz.User{
		ID:       user.Id,
		Name:     user.Name,
		Password: user.Password,
	}, nil
}

func (u *userRepo) cacheGetUser(ctx context.Context, id int64) *User {
	user, _ := cache.Get[User](ctx, u.data.rdb, userCacheKey(id))
	return user
}

func (u *userRepo) cacheSetUser(ctx context.Context, user *User) {
	err := cache.Set(ctx, u.data.rdb, userCacheKey(user.Id), user)
	if err != nil {
		u.logger.Errorf("cacheSetUser error: %v", err)
	}
}

func (u *userRepo) cacheDeleteUser(ctx context.Context, id int64) error {
	err := cache.Del(ctx, u.data.rdb, userCacheKey(id))
	if err != nil {
		u.logger.Errorf("cacheDeleteUser error: %v", err)
		return err
	}
	return nil
}
