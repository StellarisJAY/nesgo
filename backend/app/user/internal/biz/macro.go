package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	userAPI "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
	"go.mongodb.org/mongo-driver/mongo"
)

// MacroAction 宏动作，宏指令的组成部分
type MacroAction struct {
	EmulatorKey  string `json:"emulatorKey" bson:"emulatorKey"`   // 模拟器按键名称
	ReleaseDelay int64  `json:"releaseDelay" bson:"releaseDelay"` // 按键释放延迟，毫秒
}

// Macro 宏指令，用户可按键组合
type Macro struct {
	Id          int64         `json:"id" bson:"id"`
	UserId      int64         `json:"userId" bson:"userId"`
	Name        string        `json:"name" bson:"name"`
	KeyboardKey string        `json:"keyboardKey" bson:"keyboardKey"`
	Actions     []MacroAction `json:"actions" bson:"actions"` // 宏指令的动作列表
}

type MacroRepo interface {
	CreateMacro(ctx context.Context, macro *Macro) error
	GetMacro(ctx context.Context, id int64) (*Macro, error)
	ListMacro(ctx context.Context, userId int64, page, pageSize int32) ([]*Macro, int32, error)
	DeleteMacro(ctx context.Context, id int64) error
	GetMacroByName(ctx context.Context, userId int64, name string) (*Macro, error)
}

type MacroUseCase struct {
	repo   MacroRepo
	logger *log.Helper
}

func NewMacroUseCase(repo MacroRepo, logger log.Logger) *MacroUseCase {
	return &MacroUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/macro")),
	}
}

func (uc *MacroUseCase) CreateMacro(ctx context.Context, macro *Macro) error {
	existing, _ := uc.repo.GetMacroByName(ctx, macro.UserId, macro.Name)
	if existing != nil {
		return userAPI.ErrorCreateMacroFailed("macro already exists")
	}
	err := uc.repo.CreateMacro(ctx, macro)
	if err != nil {
		return userAPI.ErrorCreateMacroFailed("database error: %v", err)
	}
	return nil
}

func (uc *MacroUseCase) GetMacro(ctx context.Context, id int64) (*Macro, error) {
	macro, err := uc.repo.GetMacro(ctx, id)
	if err != nil {
		return nil, userAPI.ErrorGetMacroFailed("database error: %v", err)
	}
	if macro == nil {
		return nil, userAPI.ErrorGetMacroFailed("macro not found")
	}
	return macro, nil
}

func (uc *MacroUseCase) ListMacro(ctx context.Context, userId int64, page, pageSize int32) ([]*Macro, int32, error) {
	macros, total, err := uc.repo.ListMacro(ctx, userId, page, pageSize)
	if err != nil {
		return nil, 0, userAPI.ErrorGetMacroFailed("database error: %v", err)
	}
	return macros, total, nil
}

func (uc *MacroUseCase) DeleteMacro(ctx context.Context, id int64) error {
	err := uc.repo.DeleteMacro(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return userAPI.ErrorDeleteMacroFailed("macro not found")
	}
	if err != nil {
		return userAPI.ErrorDeleteMacroFailed("database error: %v", err)
	}
	return nil
}
