package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	roomAPI "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	userAPI "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewRoomRepo, NewGamingRepo)

// Data .
type Data struct {
	uc     userAPI.UserClient
	rc     roomAPI.RoomClient
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

	roomConn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///nesgo.service.room"),
		grpc.WithDiscovery(discovery))
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		data.logger.Info("closing grpc connections")
		_ = userConn.Close()
		_ = roomConn.Close()
	}
	data.uc = userAPI.NewUserClient(userConn)
	data.rc = roomAPI.NewRoomClient(roomConn)
	return data, cleanup, nil
}
