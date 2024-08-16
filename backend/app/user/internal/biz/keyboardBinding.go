package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
)

type KeyboardBinding struct {
	KeyboardKey string `json:"keyboardKey" bson:"keyboardKey"`
	EmulatorKey string `json:"emulatorKey" bson:"emulatorKey"`
}

type UserKeyboardBinding struct {
	Id               int64              `json:"id" bson:"id"`
	Name             string             `json:"name" bson:"name"`
	UserId           int64              `json:"userId" bson:"userId"`
	KeyboardBindings []*KeyboardBinding `json:"keyboardBindings" bson:"keyboardBindings"`
}

type UserKeyboardBindingRepo interface {
	CreateKeyboardBinding(ctx context.Context, ub *UserKeyboardBinding) error
	UpdateKeyboardBinding(ctx context.Context, ub *UserKeyboardBinding) error
	DeleteKeyboardBinding(ctx context.Context, id int64) error
	GetKeyboardBinding(ctx context.Context, id int64) (*UserKeyboardBinding, error)
	ListUserKeyboardBinding(ctx context.Context, userId int64, page, pageSize int32) ([]*UserKeyboardBinding, int32, error)
	GetBindingByName(ctx context.Context, userId int64, name string) (*UserKeyboardBinding, error)
}

type UserKeyboardBindingUseCase struct {
	repo   UserKeyboardBindingRepo
	logger *log.Helper
}

func NewUserKeyboardBindingUseCase(repo UserKeyboardBindingRepo, logger log.Logger) *UserKeyboardBindingUseCase {
	return &UserKeyboardBindingUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/UserKeyboardBinding")),
	}
}

func (uc *UserKeyboardBindingUseCase) CreateKeyboardBinding(ctx context.Context, ub *UserKeyboardBinding) error {
	binding, _ := uc.repo.GetBindingByName(ctx, ub.UserId, ub.Name)
	if binding != nil {
		return v1.ErrorCreateKeyboardBindingFailed("name already exists")
	}
	err := uc.repo.CreateKeyboardBinding(ctx, ub)
	if err != nil {
		return v1.ErrorCreateKeyboardBindingFailed("create keyboard binding failed: %v", err)
	}
	return nil
}

func (uc *UserKeyboardBindingUseCase) UpdateKeyboardBinding(ctx context.Context, ub *UserKeyboardBinding) error {
	binding, _ := uc.repo.GetBindingByName(ctx, ub.UserId, ub.Name)
	if binding != nil {
		return v1.ErrorUpdateKeyboardBindingFailed("name already exists")
	}
	err := uc.repo.UpdateKeyboardBinding(ctx, ub)
	if err != nil {
		return v1.ErrorUpdateKeyboardBindingFailed("update keyboard binding failed: %v", err)
	}
	return nil
}

func (uc *UserKeyboardBindingUseCase) DeleteKeyboardBinding(ctx context.Context, id int64) error {
	err := uc.repo.DeleteKeyboardBinding(ctx, id)
	if err != nil {
		return v1.ErrorDeleteKeyboardBindingFailed("delete keyboard binding failed: %v", err)
	}
	return nil
}

func (uc *UserKeyboardBindingUseCase) GetKeyboardBinding(ctx context.Context, id int64) (*UserKeyboardBinding, error) {
	return uc.repo.GetKeyboardBinding(ctx, id)
}

func (uc *UserKeyboardBindingUseCase) ListUserKeyboardBinding(ctx context.Context, userId int64, page, pageSize int32) ([]*UserKeyboardBinding, int32, error) {
	return uc.repo.ListUserKeyboardBinding(ctx, userId, page, pageSize)
}
