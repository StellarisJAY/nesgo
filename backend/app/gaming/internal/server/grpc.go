package server

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/conf"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/service"
	"net/url"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, gamingService *service.GamingService) *grpc.Server {
	endpoint, _ := url.Parse(c.Grpc.Addr)
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
		grpc.Endpoint(endpoint),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterGamingServer(srv, gamingService)
	return srv
}
