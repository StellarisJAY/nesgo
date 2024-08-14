package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/admin/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAdminService)

type AdminService struct {
	v1.UnimplementedAdminServer
	gf     *biz.GameFileUseCase
	a      *biz.AdminUseCase
	r      *biz.RoomUseCase
	g      *biz.GamingUseCase
	logger *log.Helper
}

func NewAdminService(gf *biz.GameFileUseCase, a *biz.AdminUseCase, r *biz.RoomUseCase, g *biz.GamingUseCase, logger log.Logger) *AdminService {
	return &AdminService{
		gf:     gf,
		a:      a,
		r:      r,
		g:      g,
		logger: log.NewHelper(log.With(logger, "module", "service/admin")),
	}
}

func (s *AdminService) UploadGame(ctx context.Context, request *v1.UploadFileRequest) (*v1.UploadFileResponse, error) {
	err := s.gf.UploadGame(ctx, request.Name, request.Data)
	if err != nil {
		return nil, err
	}
	return &v1.UploadFileResponse{}, nil
}

func (s *AdminService) ListGames(ctx context.Context, _ *v1.ListGamesRequest) (*v1.ListGamesResponse, error) {
	games, err := s.gf.ListGames(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.ListGamesResponse{Games: games}, nil
}

func (s *AdminService) DeleteGameFiles(ctx context.Context, request *v1.DeleteGameFileRequest) (*v1.DeleteGameFileResponse, error) {
	deleted, err := s.gf.DeleteGames(ctx, request.Games)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteGameFileResponse{Deleted: deleted}, nil
}
