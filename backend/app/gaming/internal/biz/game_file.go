package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisJAY/nesgo/nes/cartridge"
	"time"
)

type GameFileUseCase struct {
	repo   GameFileRepo
	logger *log.Helper
}

func NewGameFileUseCase(repo GameFileRepo, logger log.Logger) *GameFileUseCase {
	return &GameFileUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/game_file")),
	}
}

func (uc *GameFileUseCase) UploadGameFile(ctx context.Context, game string, data []byte) error {
	cart, err := cartridge.ParseCartridgeInfo(data)
	if err != nil {
		return v1.ErrorInvalidGameFile("parse cartridge failed: %v", err)
	}
	metadata := &GameFileMetadata{
		Name:       game,
		Mapper:     cartridge.MapperToString(cart.Mapper),
		Mirroring:  cartridge.MirroringToString(cart.Mirroring),
		Size:       int32(len(data)),
		UploadTime: time.Now().UnixMilli(),
	}
	err = uc.repo.UploadGameData(ctx, game, data, metadata)
	if err != nil {
		return v1.ErrorUploadFileFailed("upload game file failed: %v", err)
	}
	return nil
}

func (uc *GameFileUseCase) ListGames(ctx context.Context, page, pageSize int) ([]*GameFileMetadata, int, error) {
	games, total, err := uc.repo.ListGames(ctx, page, pageSize)
	if err != nil {
		return nil, 0, v1.ErrorListGameFailed("list games failed: %v", err)
	}
	return games, total, nil
}

func (uc *GameFileUseCase) DeleteGames(ctx context.Context, games []string) (int, error) {
	deleted, err := uc.repo.DeleteGames(ctx, games)
	if err != nil {
		return 0, v1.ErrorDeleteGameFailed("delete game failed: %v", err)
	}
	return deleted, nil
}

func (uc *GameFileUseCase) ListSaves(ctx context.Context, roomId int64, page, pageSize int32) ([]*GameSave, int32, error) {
	saves, total, err := uc.repo.ListSaves(ctx, roomId, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return saves, total, nil
}

func (uc *GameFileUseCase) DeleteSave(ctx context.Context, saveId int64) error {
	err := uc.repo.DeleteSave(ctx, saveId)
	if err != nil {
		return v1.ErrorDeleteSaveFailed("database error: %v", err)
	}
	return nil
}
