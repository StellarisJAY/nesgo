package server

import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewRegistrar, NewDiscovery, NewGRPCServer, NewHTTPServer)

func NewRegistrar(c *conf.Registry) registry.Registrar {
	config := consulAPI.DefaultConfig()
	config.Address = c.Consul.Address
	config.Scheme = c.Consul.Scheme
	client, err := consulAPI.NewClient(config)
	if err != nil {
		panic(err)
	}
	return consul.New(client, consul.WithHealthCheck(true))
}

func NewDiscovery(c *conf.Registry) registry.Discovery {
	config := consulAPI.DefaultConfig()
	config.Address = c.Consul.Address
	config.Scheme = c.Consul.Scheme
	client, err := consulAPI.NewClient(config)
	if err != nil {
		panic(err)
	}
	return consul.New(client, consul.WithHealthCheck(true))
}
