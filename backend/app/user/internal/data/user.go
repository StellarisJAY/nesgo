package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/stellarisJAY/nesgo/backend/app/user/internal/biz"
	"gorm.io/gorm"
)

type userRepo struct {
	data   *Data
	logger *log.Helper
}

type User struct {
	gorm.Model
	Id       int64  `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name     string `gorm:"size:255" json:"name"`
	Password string `gorm:"size:255" json:"password"`
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
		Model:    gorm.Model{},
		Id:       0,
		Name:     user.Name,
		Password: user.Password,
	}
	err := u.data.db.Model(model).Create(model).Error
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

func (u *userRepo) cacheGetUser(_ context.Context, id int64) *User {
	result, err := u.data.rdb.Get(userCacheKey(id)).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		u.logger.Errorf("cacheGetUser error: %v", err)
		return nil
	}
	user := &User{}
	err = json.Unmarshal([]byte(result), user)
	if err != nil {
		u.logger.Errorf("cacheGetUser error: %v", err)
		return nil
	}
	return user
}

func (u *userRepo) cacheSetUser(_ context.Context, user *User) {
	data, err := json.Marshal(user)
	if err != nil {
		u.logger.Errorf("cacheSetUser error: %v", err)
		return
	}
	_, err = u.data.rdb.Set(userCacheKey(user.Id), data, 0).Result()
	if err != nil {
		u.logger.Errorf("cacheSetUser error: %v", err)
	}
}

func (u *userRepo) cacheDeleteUser(_ context.Context, id int64) error {
	_, err := u.data.rdb.Del(userCacheKey(id)).Result()
	if err != nil {
		u.logger.Errorf("cacheDeleteUser error: %v", err)
		return err
	}
	return nil
}
