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

func NewServer() *HttpServer {
	router := setupRouter()
	return &HttpServer{
		r:    router,
		addr: config.GetEmulatorConfig().ServerAddr,
	}
}

func (s *HttpServer) Start() error {
	return s.r.Run(s.addr)
}
