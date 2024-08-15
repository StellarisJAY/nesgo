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

func (r *gamingRepo) RestartEmulator(ctx context.Context, roomId int64, game string, endpoint string) error {
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(endpoint))
	if err != nil {
		return err
	}
	defer conn.Close()
	gamingCli := gamingAPI.NewGamingClient(conn)
	_, err = gamingCli.RestartEmulator(ctx, &gamingAPI.RestartEmulatorRequest{
		RoomId: roomId,
		Game:   game,
	})
	return err
}

func (r *gamingRepo) ListSaves(ctx context.Context, roomId int64, page, pageSize int32) ([]*biz.SaveMetadata, int32, error) {
	response, err := r.data.gc.ListSaves(ctx, &gamingAPI.ListSavesRequest{
		RoomId:   roomId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, 0, err
	}
	saves := make([]*biz.SaveMetadata, 0, len(response.Saves))
	for _, save := range response.Saves {
		saves = append(saves, &biz.SaveMetadata{
			RoomId:     roomId,
			Id:         save.Id,
			Game:       save.Game,
			CreateTime: save.CreateTime,
		})
	}
	return saves, response.Total, nil
}

func (r *gamingRepo) SaveGame(ctx context.Context, roomId int64, endpoint string) error {
	conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(endpoint))
	if err != nil {
		return err
	}
	defer conn.Close()
	gamingCli := gamingAPI.NewGamingClient(conn)
	_, err = gamingCli.SaveGame(ctx, &gamingAPI.SaveGameRequest{
		RoomId: roomId,
	})
	return err
}
