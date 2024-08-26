package nes

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/stellarisJAY/nesgo/nes/apu"
	"github.com/stellarisJAY/nesgo/nes/bus"
	"github.com/stellarisJAY/nesgo/nes/cartridge"
	"github.com/stellarisJAY/nesgo/nes/config"
	"github.com/stellarisJAY/nesgo/nes/cpu"
	"github.com/stellarisJAY/nesgo/nes/ppu"
	"github.com/stellarisJAY/nesgo/nes/trace"
)

type RawEmulator struct {
	processor *cpu.Processor
	cartridge cartridge.Cartridge
	bus       *bus.Bus
	ppu       *ppu.PPU
	joyPad1   *bus.JoyPad
	joyPad2   *bus.JoyPad
	apu       *apu.BasicAPU
	config    config.Config

	lastSnapshotTime time.Time
	m                *sync.Mutex
	snapshots        []Snapshot
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

func (e *RawEmulator) SetJoyPadButtonPressed(id int, button bus.JoyPadButton, pressed bool) {
	if id == 1 {
		e.joyPad1.SetButtonPressed(button, pressed)
	} else {
		e.joyPad2.SetButtonPressed(button, pressed)
	}
}

func (e *RawEmulator) Pause() {
	e.processor.Pause()
}

func (e *RawEmulator) Resume() {
	e.processor.Resume()
}

func (e *RawEmulator) SetCPUBoostRate(rate float64) float64 {
	return e.bus.SetCPUBoostRate(rate)
}

func (e *RawEmulator) BoostCPU(delta float64) float64 {
	return e.bus.BoostCPU(delta)
}

func (e *RawEmulator) CPUBoostRate() float64 {
	return e.bus.CPUBoostRate()
}
