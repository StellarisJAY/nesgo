//go:build sdl

package main

import (
	"context"
	"github.com/stellarisJAY/nesgo/nes"
	"github.com/stellarisJAY/nesgo/nes/config"
)

func main() {
	config.InitConfigs()
	conf := config.GetEmulatorConfig()
	var program []byte
	if p, err := nes.ReadGameFile(conf.Game); err != nil {
		panic(err)
	} else {
		program = p
	}
	e := nes.NewEmulator(program, conf)
	if conf.Disassemble {
		e.Disassemble()
		return
	} else {
		e.LoadAndRun(context.Background(), conf.EnableTrace)
	}
}
