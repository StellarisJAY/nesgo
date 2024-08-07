package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/biz"
)

type GamingService struct {
	v1.UnimplementedGamingServer
	logger *log.Helper
	gi     *biz.GameInstanceUseCase
	gf     *biz.GameFileUseCase
}

func NewGamingService(gi *biz.GameInstanceUseCase, gf *biz.GameFileUseCase, logger log.Logger) *GamingService {
	return &GamingService{
		gi:     gi,
		gf:     gf,
		logger: log.NewHelper(log.With(logger, "module", "service/gaming_service")),
	}
}

func (g *GamingService) CreateGameInstance(ctx context.Context, req *v1.CreateGameInstanceRequest) (*v1.CreateGameInstanceResponse, error) {
	instance, err := g.gi.CreateGameInstance(ctx, req.RoomId, req.Game)
	if err != nil {
		return nil, err
	}
	return &v1.CreateGameInstanceResponse{
		RoomId:  instance.RoomId,
		LeaseId: instance.LeaseID,
	}, nil
}

func (g *GamingService) OpenGameConnection(ctx context.Context, request *v1.OpenGameConnectionRequest) (*v1.OpenGameConnectionResponse, error) {
	sdp, err := g.gi.OpenGameConnection(ctx, request.RoomId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.OpenGameConnectionResponse{
		RoomId:   request.RoomId,
		UserId:   request.UserId,
		SdpOffer: sdp,
	}, nil
}

func (g *GamingService) SDPAnswer(ctx context.Context, request *v1.SDPAnswerRequest) (*v1.SDPAnswerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GamingService) ICECandidate(ctx context.Context, request *v1.ICECandidateRequest) (*v1.ICECandidateResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GamingService) PauseEmulator(ctx context.Context, request *v1.PauseEmulatorRequest) (*v1.PauseEmulatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GamingService) RestartEmulator(ctx context.Context, request *v1.RestartEmulatorRequest) (*v1.RestartEmulatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GamingService) UploadGame(ctx context.Context, request *v1.UploadGameRequest) (*v1.UploadGameResponse, error) {
	err := g.gf.UploadGameFile(ctx, request.Name, request.Data)
	if err != nil {
		return nil, err
	}
	return &v1.UploadGameResponse{}, nil
}

func (g *GamingService) ListGames(ctx context.Context, _ *v1.ListGamesRequest) (*v1.ListGamesResponse, error) {
	games, err := g.gf.ListGames(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.GameFileMetadata, 0, len(games))
	for _, game := range games {
		result = append(result, &v1.GameFileMetadata{
			Name:      game.Name,
			Mapper:    game.Mapper,
			Mirroring: game.Mirroring,
		})
	}
	return &v1.ListGamesResponse{Games: result}, nil
}

func (g *GamingService) DeleteGameFile(ctx context.Context, request *v1.DeleteGameFileRequest) (*v1.DeleteGameFileResponse, error) {
	deleted, err := g.gf.DeleteGames(ctx, request.Games)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteGameFileResponse{Deleted: int32(deleted)}, nil
}
