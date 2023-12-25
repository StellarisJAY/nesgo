package emulator

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/cartridge"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/cpu"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/stellarisJAY/nesgo/trace"
	"io"
	"log"
	"os"
)

type RawEmulator struct {
	processor *cpu.Processor
	cartridge cartridge.Cartridge
	bus       *bus.Bus
	ppu       *ppu.PPU
	joyPad    *bus.JoyPad

	config config.Config
}

func ReadGameFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("can't open game file %s,  %w", fileName, err)
	}
	program, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read game file error %w", err)
	}
	log.Printf("loaded program file: %s, size: %d", fileName, len(program))
	return program, nil
}

func (e *RawEmulator) Disassemble() {
	e.processor.Disassemble(trace.PrintDisassemble)
}

func (e *RawEmulator) SetJoyPadButtonPressed(button bus.JoyPadButton, pressed bool) {
	e.joyPad.SetButtonPressed(button, pressed)
}
