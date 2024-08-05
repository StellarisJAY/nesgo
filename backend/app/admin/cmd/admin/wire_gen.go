// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/biz"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/conf"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/data"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/server"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, registry *conf.Registry, logger log.Logger) (*kratos.App, func(), error) {
	discovery := server.NewDiscovery(registry)
	dataData, cleanup, err := data.NewData(registry, discovery, logger)
	if err != nil {
		return nil, nil, err
	}
	gameFileRepo := data.NewGameFileRepo(dataData, logger)
	gameFileUseCase := biz.NewGameFileUseCase(gameFileRepo, logger)
	adminService := service.NewAdminService(gameFileUseCase, logger)
	grpcServer := server.NewGRPCServer(confServer, adminService, logger)
	registrar := server.NewRegistrar(registry)
	httpServer := server.NewHTTPServer(confServer, adminService, logger)
	app := newApp(logger, grpcServer, registrar, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
