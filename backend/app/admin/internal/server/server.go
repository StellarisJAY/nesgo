package server

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewRegistrar, NewDiscovery, NewGRPCServer, NewHTTPServer)

func NewRegistrar(c *conf.Registry) registry.Registrar {
	client, err := consulAPI.NewClient(&consulAPI.Config{
		Address: c.Consul.Address,
	})
	if err != nil {
		panic(err)
	}
	return consul.New(client, consul.WithHealthCheck(true))
}

func NewDiscovery(c *conf.Registry) registry.Discovery {
	client, err := consulAPI.NewClient(&consulAPI.Config{
		Address: c.Consul.Address,
	})
	if err != nil {
		panic(err)
	}
	return consul.New(client, consul.WithHealthCheck(true))
}
