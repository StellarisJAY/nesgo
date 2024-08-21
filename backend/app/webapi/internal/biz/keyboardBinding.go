package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
)

type KeyboardBinding struct {
	KeyboardKey string
	EmulatorKey string
}

type UserKeyboardBinding struct {
	Id               int64
	Name             string
	UserId           int64
	KeyboardBindings []*KeyboardBinding
}

type UserKeyboardBindingRepo interface {
	CreateKeyboardBinding(ctx context.Context, ub *UserKeyboardBinding) error
	UpdateKeyboardBinding(ctx context.Context, ub *UserKeyboardBinding) error
	DeleteKeyboardBinding(ctx context.Context, id int64) error
	GetKeyboardBinding(ctx context.Context, id int64) (*UserKeyboardBinding, error)
	ListUserKeyboardBinding(ctx context.Context, userId int64, page, pageSize int32) ([]*UserKeyboardBinding, int32, error)
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

func TranslateKeyboardKey(key string) string {
	// TODO 键盘按键编号 -> 可读格式
	return key
}

func TranslateEmulatorKey(key string) string {
	// TODO 模拟器按键编号 -> 可读格式
	return key
}

func (uc *UserKeyboardBindingUseCase) Create(ctx context.Context, ub *UserKeyboardBinding) error {
	return uc.repo.CreateKeyboardBinding(ctx, ub)
}

func (uc *UserKeyboardBindingUseCase) List(ctx context.Context, userId int64, page, pageSize int32) ([]*UserKeyboardBinding, int32, error) {
	return uc.repo.ListUserKeyboardBinding(ctx, userId, page, pageSize)
}

func (uc *UserKeyboardBindingUseCase) Get(ctx context.Context, id int64) (*UserKeyboardBinding, error) {
	return uc.repo.GetKeyboardBinding(ctx, id)
}

func (uc *UserKeyboardBindingUseCase) Update(ctx context.Context, ub *UserKeyboardBinding) error {
	return uc.repo.UpdateKeyboardBinding(ctx, ub)
}

func (uc *UserKeyboardBindingUseCase) Delete(ctx context.Context, id, userId int64) error {
	binding, err := uc.repo.GetKeyboardBinding(ctx, id)
	if err != nil {
		return err
	}
	if binding.UserId != userId {
		return v1.ErrorOperationFailed("can not delete other user's keyboard binding")
	}
	return uc.repo.DeleteKeyboardBinding(ctx, id)
}
