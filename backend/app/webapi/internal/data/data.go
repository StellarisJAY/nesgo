package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	userAPI "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo)

// Data .
type Data struct {
	uc     userAPI.UserClient
	logger *log.Helper
}

func NewData(discovery registry.Discovery, logger log.Logger) (*Data, func(), error) {
	data := new(Data)
	data.logger = log.NewHelper(log.With(logger, "module", "data"))
	userConn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///nesgo.service.user"),
		grpc.WithDiscovery(discovery))
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		data.logger.Info("closing grpc connections")
		_ = userConn.Close()
	}
	data.uc = userAPI.NewUserClient(userConn)
	return data, cleanup, nil
}
