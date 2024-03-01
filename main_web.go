//go:build web

package main

import (
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/web"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	server := web.NewServer()
	if config.GetEmulatorConfig().Pprof {
		go func() {
			_ = http.ListenAndServe(":9999", nil)
		}()
	}
	_ = server.Start()
}
