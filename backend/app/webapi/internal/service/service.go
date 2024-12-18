package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewWebApiService)

type WebApiService struct {
	v1.UnimplementedWebApiServer
	uc     *biz.UserUseCase
	ac     *biz.AuthUseCase
	rc     *biz.RoomUseCase
	gc     *biz.GamingUseCase
	uk     *biz.UserKeyboardBindingUseCase
	mu     *biz.MacroUseCase
	logger *log.Helper
}

func NewWebApiService(uc *biz.UserUseCase,
	ac *biz.AuthUseCase,
	rc *biz.RoomUseCase,
	gc *biz.GamingUseCase,
	uk *biz.UserKeyboardBindingUseCase,
	mu *biz.MacroUseCase, logger log.Logger) *WebApiService {
	return &WebApiService{
		uc:     uc,
		ac:     ac,
		rc:     rc,
		gc:     gc,
		uk:     uk,
		mu:     mu,
		logger: log.NewHelper(log.With(logger, "module", "service/webapi")),
	}
}
