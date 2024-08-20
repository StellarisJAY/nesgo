package server

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/conf"
	etcdAPI "go.etcd.io/etcd/client/v3"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewEtcdClient, NewRegistrar, NewGRPCServer, NewDiscovery)

func NewEtcdClient(c *conf.Registry) *etcdAPI.Client {
	cli, err := etcdAPI.New(etcdAPI.Config{
		Endpoints: c.Etcd.Endpoints,
	})
	if err != nil {
		panic(err)
	}
	return cli
}

func NewRegistrar(cli *etcdAPI.Client) registry.Registrar {
	return etcd.New(cli)
}

func NewDiscovery(cli *etcdAPI.Client) registry.Discovery {
	return etcd.New(cli)
}
