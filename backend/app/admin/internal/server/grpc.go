package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/admin/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/conf"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, as *service.AdminService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
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
	v1.RegisterAdminServer(srv, as)
	return srv
}
