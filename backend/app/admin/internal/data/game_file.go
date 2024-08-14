package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	adminAPI "github.com/stellarisJAY/nesgo/backend/api/app/admin/v1"
	gamingAPI "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/biz"
)

type gameFileRepo struct {
	data   *Data
	logger *log.Helper
}

func NewGameFileRepo(data *Data, logger log.Logger) biz.GameFileRepo {
	return &gameFileRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/game_file")),
	}
}

func (g *gameFileRepo) UploadGame(ctx context.Context, name string, data []byte) error {
	_, err := g.data.gamingCli.UploadGame(ctx, &gamingAPI.UploadGameRequest{Name: name, Data: data})
	if err != nil {
		return err
	}
	return nil
}

func (g *gameFileRepo) ListGames(ctx context.Context, page, pageSize int32) ([]*adminAPI.GameFileMetadata, int32, error) {
	response, err := g.data.gamingCli.ListGames(ctx, &gamingAPI.ListGamesRequest{Page: page, PageSize: pageSize})
	if err != nil {
		return nil, 0, err
	}
	result := make([]*adminAPI.GameFileMetadata, 0, len(response.Games))
	for _, game := range response.Games {
		result = append(result, &adminAPI.GameFileMetadata{
			Name:      game.Name,
			Mapper:    game.Mapper,
			Mirroring: game.Mirroring,
		})
	}
	return result, response.Total, nil
}

func (g *gameFileRepo) DeleteGames(ctx context.Context, games []string) (int32, error) {
	response, err := g.data.gamingCli.DeleteGameFile(ctx, &gamingAPI.DeleteGameFileRequest{Games: games})
	if err != nil {
		return 0, err
	}
	return response.Deleted, nil
}
