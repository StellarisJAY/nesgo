package server

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/stellarisjay/nesgo/backend/app/user/internal/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewRegistrar, NewGRPCServer)

func NewRegistrar(c *conf.Registry) registry.Registrar {
	config := consulAPI.DefaultConfig()
	config.Address = c.Consul.Address
	config.Scheme = c.Consul.Scheme
	client, err := consulAPI.NewClient(config)
	if err != nil {
		panic(err)
	}
	r := consul.New(client, consul.WithHealthCheck(true))
	return r
}
