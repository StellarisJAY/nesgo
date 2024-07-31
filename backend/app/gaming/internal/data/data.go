package data

import (
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/stellarisjay/nesgo/backend/app/gaming/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRoomSessionRepo)

// Data .
type Data struct {
	consul *consulAPI.Client
}

// NewData .
func NewData(c *conf.Registry, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data"))
	client, err := consulAPI.NewClient(&consulAPI.Config{
		Address: c.Consul.Address,
		Scheme:  c.Consul.Scheme,
	})
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		logHelper.Info("closing the data resources")
	}
	return &Data{
		consul: client,
	}, cleanup, nil
}
