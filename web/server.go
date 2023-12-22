package web

import (
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/emulator"
)

type HttpServer struct {
	r    *gin.Engine
	addr string
	e    *emulator.Emulator
}

func NewServer(config config.Config) *HttpServer {
	router := setupRouter(config)
	return &HttpServer{
		r:    router,
		addr: config.ServerAddr,
	}
}

func (s *HttpServer) Start() error {
	return s.r.Run(s.addr)
}
