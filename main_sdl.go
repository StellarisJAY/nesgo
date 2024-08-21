//go:build sdl

package main

import (
	"context"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisJAY/nesgo/emulator/config"
)

func main() {
	config.InitConfigs()
	conf := config.GetEmulatorConfig()
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
