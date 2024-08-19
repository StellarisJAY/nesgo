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
	err := g.gi.SDPAnswer(ctx, request.RoomId, request.UserId, request.SdpAnswer)
	if err != nil {
		return nil, err
	}
	return &v1.SDPAnswerResponse{
		RoomId: request.RoomId,
		UserId: request.UserId,
	}, nil
}

func (g *GamingService) ICECandidate(ctx context.Context, request *v1.ICECandidateRequest) (*v1.ICECandidateResponse, error) {
	err := g.gi.ICECandidate(ctx, request.RoomId, request.UserId, request.Candidate)
	if err != nil {
		return nil, err
	}
	return &v1.ICECandidateResponse{
		RoomId: request.RoomId,
		UserId: request.UserId,
	}, nil
}

func (g *GamingService) PauseEmulator(ctx context.Context, request *v1.PauseEmulatorRequest) (*v1.PauseEmulatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GamingService) RestartEmulator(ctx context.Context, request *v1.RestartEmulatorRequest) (*v1.RestartEmulatorResponse, error) {
	err := g.gi.RestartEmulator(ctx, request.RoomId, request.Game)
	if err != nil {
		return nil, err
	}
	return &v1.RestartEmulatorResponse{}, nil
}

func (g *GamingService) UploadGame(ctx context.Context, request *v1.UploadGameRequest) (*v1.UploadGameResponse, error) {
	err := g.gf.UploadGameFile(ctx, request.Name, request.Data)
	if err != nil {
		return nil, err
	}
	return &v1.UploadGameResponse{}, nil
}

func (g *GamingService) ListGames(ctx context.Context, request *v1.ListGamesRequest) (*v1.ListGamesResponse, error) {
	games, total, err := g.gf.ListGames(ctx, int(request.Page), int(request.PageSize))
	if err != nil {
		return nil, err
	}
	result := make([]*v1.GameFileMetadata, 0, len(games))
	for _, game := range games {
		result = append(result, &v1.GameFileMetadata{
			Name:       game.Name,
			Mapper:     game.Mapper,
			Mirroring:  game.Mirroring,
			Size:       game.Size,
			UploadTime: game.UploadTime,
		})
	}
	return &v1.ListGamesResponse{Games: result, Total: int32(total)}, nil
}

func (g *GamingService) DeleteGameFile(ctx context.Context, request *v1.DeleteGameFileRequest) (*v1.DeleteGameFileResponse, error) {
	deleted, err := g.gf.DeleteGames(ctx, request.Games)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteGameFileResponse{Deleted: int32(deleted)}, nil
}

func (g *GamingService) DeleteGameInstance(ctx context.Context, request *v1.DeleteGameInstanceRequest) (*v1.DeleteGameInstanceResponse, error) {
	err := g.gi.ReleaseGameInstance(ctx, request.RoomId, request.Force)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteGameInstanceResponse{
		RoomId: request.RoomId,
	}, nil
}

func (g *GamingService) SetController(ctx context.Context, request *v1.SetControllerRequest) (*v1.SetControllerResponse, error) {
	err := g.gi.SetController(ctx, request.RoomId, request.UserId, int(request.Controller))
	if err != nil {
		return nil, err
	}
	return &v1.SetControllerResponse{}, nil
}

func (g *GamingService) GetControllers(ctx context.Context, request *v1.GetControllersRequest) (*v1.GetControllersResponse, error) {
	p1, p2, err := g.gi.GetController(ctx, request.RoomId)
	if err != nil {
		return nil, err
	}
	return &v1.GetControllersResponse{
		Controller1: p1,
		Controller2: p2,
	}, nil
}

func (g *GamingService) GetGameInstanceStats(ctx context.Context, request *v1.GetGameInstanceStatsRequest) (*v1.GetGameInstanceStatsResponse, error) {
	stats, err := g.gi.GetGameInstanceStats(ctx, request.RoomId)
	if err != nil {
		return nil, err
	}
	return &v1.GetGameInstanceStatsResponse{
		Stats: &v1.GameInstanceStats{
			RoomId:            stats.RoomId,
			Connections:       int32(stats.Connections),
			ActiveConnections: int32(stats.ActiveConnections),
			Game:              stats.Game,
			Uptime:            stats.Uptime.Milliseconds(),
		},
	}, nil
}

func (g *GamingService) ListGameInstances(ctx context.Context, _ *v1.ListGameInstancesRequest) (*v1.ListGameInstancesResponse, error) {
	instances, err := g.gi.ListGameInstances(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.GameInstanceStats, 0, len(instances))
	for _, instance := range instances {
		result = append(result, &v1.GameInstanceStats{
			RoomId:            instance.RoomId,
			Connections:       int32(instance.Connections),
			ActiveConnections: int32(instance.ActiveConnections),
			Game:              instance.Game,
			Uptime:            instance.Uptime.Milliseconds(),
		})
	}
	return &v1.ListGameInstancesResponse{Instances: result}, nil
}

func (g *GamingService) DeleteMemberConnection(ctx context.Context, request *v1.DeleteMemberConnectionRequest) (*v1.DeleteMemberConnectionResponse, error) {
	err := g.gi.DeleteMemberConnection(ctx, request.RoomId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteMemberConnectionResponse{}, nil
}

func (g *GamingService) GetEndpointStats(ctx context.Context, request *v1.GetEndpointStatsRequest) (*v1.GetEndpointStatsResponse, error) {
	stats, err := g.gi.GetEndpointStats(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetEndpointStatsResponse{
		EmulatorCount: stats.EmulatorCount,
		CpuUsage:      stats.CpuUsage,
		MemoryUsed:    stats.MemoryUsed,
		MemoryTotal:   stats.MemoryTotal,
		Uptime:        stats.Uptime,
	}, nil
}

func (g *GamingService) SaveGame(ctx context.Context, request *v1.SaveGameRequest) (*v1.SaveGameResponse, error) {
	err := g.gi.SaveGame(ctx, request.RoomId)
	if err != nil {
		return nil, err
	}
	return &v1.SaveGameResponse{}, nil
}

func (g *GamingService) LoadSave(ctx context.Context, request *v1.LoadSaveRequest) (*v1.LoadSaveResponse, error) {
	err := g.gi.LoadSave(ctx, request.SaveId)
	if err != nil {
		return nil, err
	}
	return &v1.LoadSaveResponse{}, nil
}

func (g *GamingService) ListSaves(ctx context.Context, request *v1.ListSavesRequest) (*v1.ListSavesResponse, error) {
	saves, total, err := g.gf.ListSaves(ctx, request.RoomId, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.Save, 0, len(saves))
	for _, save := range saves {
		result = append(result, &v1.Save{
			RoomId:     save.RoomId,
			Id:         save.Id,
			CreateTime: save.CreateTime,
			Game:       save.Game,
		})
	}
	return &v1.ListSavesResponse{Saves: result, Total: total}, nil
}

func (g *GamingService) DeleteSave(ctx context.Context, request *v1.DeleteSaveRequest) (*v1.DeleteSaveResponse, error) {
	err := g.gf.DeleteSave(ctx, request.SaveId)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteSaveResponse{}, nil
}

func (g *GamingService) GetServerICECandidate(ctx context.Context, request *v1.GetServerICECandidateRequest) (*v1.GetServerICECandidateResponse, error) {
	candidates, err := g.gi.GetICECandidates(ctx, request.RoomId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetServerICECandidateResponse{Candidates: candidates}, nil
}
