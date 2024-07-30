//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/stellarisjay/nesgo/backend/app/room/internal/biz"
	"github.com/stellarisjay/nesgo/backend/app/room/internal/conf"
	"github.com/stellarisjay/nesgo/backend/app/room/internal/data"
	"github.com/stellarisjay/nesgo/backend/app/room/internal/server"
	"github.com/stellarisjay/nesgo/backend/app/room/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
