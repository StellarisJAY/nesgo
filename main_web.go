//go:build web

package main

import (
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/web"
)

func main() {
	conf := config.ParseConfig()
	server := web.NewServer(conf)
	_ = server.Start()
}
