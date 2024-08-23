package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type MacroAction struct {
	EmulatorKey  string `json:"emulatorKey" bson:"emulatorKey"`
	ReleaseDelay int64  `json:"releaseDelay" bson:"releaseDelay"`
}

type Macro struct {
	Id          int64         `json:"id" bson:"id"`
	UserId      int64         `json:"userId" bson:"userId"`
	Name        string        `json:"name" bson:"name"`
	KeyboardKey string        `json:"keyboardKey" bson:"keyboardKey"`
	Actions     []MacroAction `json:"actions" bson:"actions"`
}

type MacroRepo interface {
	CreateMacro(ctx context.Context, macro *Macro) error
	GetMacro(ctx context.Context, id int64) (*Macro, error)
	ListMacro(ctx context.Context, userId int64, page, pageSize int32) ([]*Macro, int32, error)
	DeleteMacro(ctx context.Context, id int64) error
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
	return uc.repo.CreateMacro(ctx, macro)
}

func (uc *MacroUseCase) GetMacro(ctx context.Context, id int64) (*Macro, error) {
	return uc.repo.GetMacro(ctx, id)
}

func (uc *MacroUseCase) ListMacro(ctx context.Context, userId int64, page, pageSize int32) ([]*Macro, int32, error) {
	return uc.repo.ListMacro(ctx, userId, page, pageSize)
}

func (uc *MacroUseCase) DeleteMacro(ctx context.Context, id int64) error {
	return uc.repo.DeleteMacro(ctx, id)
}
