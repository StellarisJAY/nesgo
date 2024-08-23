// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/conf"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/data"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/server"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, registry *conf.Registry, auth *conf.Auth, logger log.Logger) (*kratos.App, func(), error) {
	client := server.NewEtcdClient(registry)
	discovery := server.NewDiscovery(client)
	dataData, cleanup, err := data.NewData(discovery, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUseCase := biz.NewUserUseCase(userRepo, logger)
	authUseCase := biz.NewAuthUseCase(userRepo, auth, logger)
	roomRepo := data.NewRoomRepo(dataData, logger)
	gamingRepo := data.NewGamingRepo(dataData, logger)
	roomUseCase := biz.NewRoomUseCase(roomRepo, userRepo, gamingRepo, logger)
	gamingUseCase := biz.NewGamingUseCase(roomRepo, gamingRepo, logger)
	userKeyboardBindingRepo := data.NewUserKeyboardBindingRepo(dataData, logger)
	userKeyboardBindingUseCase := biz.NewUserKeyboardBindingUseCase(userKeyboardBindingRepo, logger)
	macroRepo := data.NewMacroRepo(dataData, logger)
	macroUseCase := biz.NewMacroUseCase(macroRepo, logger)
	webApiService := service.NewWebApiService(userUseCase, authUseCase, roomUseCase, gamingUseCase, userKeyboardBindingUseCase, macroUseCase, logger)
	grpcServer := server.NewGRPCServer(confServer, webApiService, logger)
	httpServer := server.NewHTTPServer(confServer, auth, webApiService, logger)
	registrar := server.NewRegistrar(client)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
