package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
)

type GamingService struct {
	v1.UnimplementedGamingServer
	logger *log.Helper
}

func NewGamingService(logger log.Logger) *GamingService {
	return &GamingService{
		logger: log.NewHelper(log.With(logger, "module", "service/gaming_service")),
	}
}

func (g *GamingService) CreateGameInstance(ctx context.Context, req *v1.CreateGameInstanceRequest) (*v1.CreateGameInstanceResponse, error) {
	panic("implement me")
}

func (g *GamingService) OpenGameConnection(ctx context.Context, request *v1.OpenGameConnectionRequest) (*v1.OpenGameConnectionResponse, error) {
	//TODO implement me
	panic("implement me")
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
