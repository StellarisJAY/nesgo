package main

import (
	"context"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisJAY/nesgo/web"
)

func main() {
	conf := config.ParseConfig()
	if conf.Web {
		server := web.NewServer(conf)
		_ = server.Start()
	} else {
		var program []byte
		if p, err := emulator.ReadGameFile(conf.Game); err != nil {
			panic(err)
		} else {
			program = p
		}
		e := emulator.NewEmulator(program, conf)
		if conf.Disassemble {
			e.Disassemble()
			return
		} else {
			e.LoadAndRun(context.Background(), conf.EnableTrace)
		}
	}
}
