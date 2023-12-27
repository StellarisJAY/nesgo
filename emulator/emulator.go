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
	"time"
)

const (
	snapshotInterval = 30 * time.Second
	maxSnapshots     = 4
)

type RawEmulator struct {
	processor *cpu.Processor
	cartridge cartridge.Cartridge
	bus       *bus.Bus
	ppu       *ppu.PPU
	joyPad    *bus.JoyPad

	config config.Config

	lastSnapshotTime time.Time
	snapshots        []Snapshot
}

type Snapshot struct {
	processor cpu.Snapshot
	ppu       ppu.Snapshot
	bus       bus.Snapshot
	timestamp time.Time
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

func (e *RawEmulator) MakeSnapshot() {
	e.lastSnapshotTime = time.Now()
	s := Snapshot{
		processor: e.processor.MakeSnapshot(),
		ppu:       e.ppu.MakeSnapshot(),
		bus:       e.bus.MakeSnapshot(),
		timestamp: e.lastSnapshotTime,
	}
	e.snapshots = append(e.snapshots, s)
	if len(e.snapshots) > maxSnapshots {
		e.snapshots = e.snapshots[1:]
	}
	log.Println("new snapshot done")
}

func (e *RawEmulator) Pause() {
	e.processor.Pause()
}

func (e *RawEmulator) Resume() {
	e.processor.Resume()
}
