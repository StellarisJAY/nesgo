//go:build sdl

package main

import (
	"context"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/emulator"
)

func main() {
	conf := config.ParseConfig()
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
