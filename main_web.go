//go:build web

package main

import (
	"github.com/stellarisJAY/nesgo/web"
)

func main() {
	server := web.NewServer()
	_ = server.Start()
}
