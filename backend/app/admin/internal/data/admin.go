package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/biz"
	"time"
)

type Admin struct {
	Id        int64  `gorm:"primary_key;auto_increment"`
	Name      string `gorm:"size:255"`
	Password  string `gorm:"size:255"`
	CreatedAt time.Time
}

type adminRepo struct {
	data   *Data
	logger *log.Helper
}

func NewAdminRepo(data *Data, logger log.Logger) biz.AdminRepo {
	return &adminRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/admin")),
	}
}

func (r *adminRepo) GetAdmin(ctx context.Context, id int64) (*biz.Admin, error) {
	admin := &Admin{}
	err := r.data.db.Model(admin).
		Where("id = ?", id).
		First(admin).
		Error
	if err != nil {
		return nil, err
	}
	return admin.ToBizAdmin(), nil
}

func (r *adminRepo) GetAdminByName(ctx context.Context, name string) (*biz.Admin, error) {
	admin := &Admin{}
	err := r.data.db.Model(admin).
		Where("name = ?", name).
		First(admin).
		Error
	if err != nil {
		return nil, err
	}
	return admin.ToBizAdmin(), nil
}

func (r *adminRepo) CreateAdmin(ctx context.Context, admin *biz.Admin) error {
	model := &Admin{
		Name:      admin.Name,
		Password:  admin.Password,
		CreatedAt: time.Now(),
	}
	return r.data.db.Model(model).
		Create(model).
		Error
}

func (a *Admin) ToBizAdmin() *biz.Admin {
	return &biz.Admin{
		Id:        a.Id,
		Name:      a.Name,
		Password:  a.Password,
		CreatedAt: a.CreatedAt,
	}
}
