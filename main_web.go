//go:build !sdl

package main

import (
	"github.com/stellarisJAY/nesgo/emulator/config"
	"github.com/stellarisJAY/nesgo/web"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	config.InitConfigs()
	server := web.NewServer()
	if config.GetEmulatorConfig().Debug {
		go func() {
			_ = http.ListenAndServe(":9999", nil)
		}()
	}
	_ = server.Start()
}
