package web

import (
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/emulator/config"
)

type HttpServer struct {
	r    *gin.Engine
	addr string
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
