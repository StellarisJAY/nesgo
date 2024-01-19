//go:build web

package emulator

import (
	"context"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/cartridge"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/cpu"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/stellarisJAY/nesgo/trace"
	"sync"
)

// Emulator browser render emulator
type Emulator struct {
	RawEmulator
}

func NewEmulator(game string, conf config.Config, callback bus.RenderCallback) (*Emulator, error) {
	nesData, err := ReadGameFile(game)
	if err != nil {
		return nil, err
	}
	e := &Emulator{
		RawEmulator{
			cartridge: cartridge.MakeCartridge(nesData),
			config:    conf,
			m:         &sync.Mutex{},
		},
	}
	e.joyPad = bus.NewJoyPad()
	e.ppu = ppu.NewPPU(e.cartridge.GetChrBank, e.cartridge.GetMirroring, e.cartridge.WriteCHR)
	e.bus = bus.NewBus(e.cartridge, e.ppu, callback, e.joyPad)
	e.processor = cpu.NewProcessor(e.bus)
	e.init()
	return e, nil
}

func (e *Emulator) LoadAndRun(ctx context.Context, enableTrace bool) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	if enableTrace {
		e.processor.LoadAndRunWithCallback(ctx, trace.Trace, nil)
	} else {
		e.processor.LoadAndRunWithCallback(ctx, nil,
			func(_ *cpu.Processor) {
				e.PushSnapshot()
			})
	}
}
