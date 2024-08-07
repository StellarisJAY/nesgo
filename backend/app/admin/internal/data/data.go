package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	gaming "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGameFileRepo)

// Data .
type Data struct {
	gamingCli gaming.GamingClient
	logger    *log.Helper
}

// NewData .
func NewData(c *conf.Registry, r registry.Discovery, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data"))
	gamingConn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///nesgo.service.gaming"),
		grpc.WithDiscovery(r))
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		logHelper.Info("closing data resources")
		_ = gamingConn.Close()
	}
	return &Data{
		logger:    logHelper,
		gamingCli: gaming.NewGamingClient(gamingConn),
	}, cleanup, nil
}
