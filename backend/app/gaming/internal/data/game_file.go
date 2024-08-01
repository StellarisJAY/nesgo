package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/biz"
)

type gameFileRepo struct {
	logger *log.Helper
	data   *Data
}

func NewGameFileRepo(data *Data, logger log.Logger) biz.GameFileRepo {
	return &gameFileRepo{
		logger: log.NewHelper(log.With(logger, "module", "data/gameFile")),
		data:   data,
	}
}

func (g *gameFileRepo) GetGameData(ctx context.Context, game string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (g *gameFileRepo) GetSavedGame(ctx context.Context, id int64) (*biz.GameSave, error) {
	//TODO implement me
	panic("implement me")
}

func (g *gameFileRepo) SaveGame(ctx context.Context, save *biz.GameSave) error {
	//TODO implement me
	panic("implement me")
}
