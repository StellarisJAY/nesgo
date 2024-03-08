package emulator

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/emulator/apu"
	"github.com/stellarisJAY/nesgo/emulator/bus"
	"github.com/stellarisJAY/nesgo/emulator/cartridge"
	"github.com/stellarisJAY/nesgo/emulator/config"
	"github.com/stellarisJAY/nesgo/emulator/cpu"
	"github.com/stellarisJAY/nesgo/emulator/ppu"
	"github.com/stellarisJAY/nesgo/emulator/trace"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	snapshotInterval = 5 * time.Second // 快照间隔5s
	maxSnapshots     = 12
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

// ReverseOnce 回退游戏进度，返回回退后的画面帧
func (e *RawEmulator) ReverseOnce() []byte {
	// snapshot队列可能被多线程访问
	e.m.Lock()
	defer e.m.Unlock()
	if len(e.snapshots) > 0 {
		// 暂停cpu循环
		e.Pause()
		defer e.Resume()
		last := e.snapshots[len(e.snapshots)-1]
		e.processor.Reverse(last.Processor)
		revFrame := e.ppu.Reverse(last.PPU)
		e.bus.Reverse(last.Bus)
		e.snapshots = e.snapshots[:len(e.snapshots)-1]
		return revFrame
	}
	return nil
}

// ReverseOnceNoBlock 不阻塞的回溯，只能在单线程模式下使用
func (e *RawEmulator) ReverseOnceNoBlock() []byte {
	if len(e.snapshots) > 0 {
		last := e.snapshots[len(e.snapshots)-1]
		e.processor.Reverse(last.Processor)
		revFrame := e.ppu.Reverse(last.PPU)
		e.bus.Reverse(last.Bus)
		e.snapshots = e.snapshots[:len(e.snapshots)-1]
		return revFrame
	}
	return nil
}

func (e *RawEmulator) SetCPUBoostRate(rate float64) float64 {
	return e.bus.SetCPUBoostRate(rate)
}

func (e *RawEmulator) BoostCPU(delta float64) float64 {
	return e.bus.BoostCPU(delta)
}
