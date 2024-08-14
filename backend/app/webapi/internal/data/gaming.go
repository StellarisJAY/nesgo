package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
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

func (r *gamingRepo) DeleteMemberConnection(ctx context.Context, roomId, userId int64, endpoint string) error {
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(endpoint))
	if err != nil {
		return err
	}
	defer conn.Close()
	gamingCli := gamingAPI.NewGamingClient(conn)
	_, err = gamingCli.DeleteMemberConnection(ctx, &gamingAPI.DeleteMemberConnectionRequest{
		RoomId: roomId,
		UserId: userId,
	})
	if err != nil {
		return err
	}
	return nil
}
