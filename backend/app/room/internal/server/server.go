package server

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/stellarisjay/nesgo/backend/app/room/internal/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewRegistrar, NewGRPCServer)

func NewRegistrar(registry *conf.Registry) registry.Registrar {
	client, err := consulAPI.NewClient(&consulAPI.Config{
		Address: registry.Consul.Address,
		Scheme:  registry.Consul.Scheme,
	})
	if err != nil {
		panic(err)
	}
	return consul.New(client, consul.WithHealthCheck(true))
}
