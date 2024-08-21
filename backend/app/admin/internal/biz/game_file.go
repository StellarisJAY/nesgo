package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/admin/v1"
)

type GameFileRepo interface {
	UploadGame(ctx context.Context, name string, data []byte) error
	ListGames(ctx context.Context, page, pageSize int32) ([]*v1.GameFileMetadata, int32, error)
	DeleteGames(ctx context.Context, games []string) (int32, error)
}

type GameFileUseCase struct {
	repo   GameFileRepo
	logger *log.Helper
}

func NewGameFileUseCase(repo GameFileRepo, logger log.Logger) *GameFileUseCase {
	return &GameFileUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/GameFile")),
	}
}

func (uc *GameFileUseCase) UploadGame(ctx context.Context, name string, data []byte) error {
	err := uc.repo.UploadGame(ctx, name, data)
	return err
}

func (uc *GameFileUseCase) ListGames(ctx context.Context, page, pageSize int32) ([]*v1.GameFileMetadata, int32, error) {
	games, total, err := uc.repo.ListGames(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return games, total, nil
}

func (uc *GameFileUseCase) DeleteGames(ctx context.Context, games []string) (int32, error) {
	deleted, err := uc.repo.DeleteGames(ctx, games)
	return deleted, err
}
