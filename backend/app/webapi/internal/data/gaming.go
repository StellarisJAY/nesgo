package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	gamingAPI "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
)

type gamingRepo struct {
	data   *Data
	logger *log.Helper
}

func NewGamingRepo(data *Data, logger log.Logger) biz.GamingRepo {
	return &gamingRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/gaming")),
	}
}

func (r *gamingRepo) ListGames(ctx context.Context) ([]*biz.GameMetadata, error) {
	response, err := r.data.gc.ListGames(ctx, &gamingAPI.ListGamesRequest{})
	if err != nil {
		return nil, err
	}
	result := make([]*biz.GameMetadata, 0, len(response.Games))
	for _, game := range response.Games {
		result = append(result, &biz.GameMetadata{
			Name:      game.Name,
			Mapper:    game.Mapper,
			Mirroring: game.Mirroring,
		})
	}
	return result, nil
}
