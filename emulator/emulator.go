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
	joyPad    *bus.JoyPad

	config config.Config

	lastSnapshotTime time.Time
	m                *sync.Mutex
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
	if time.Now().After(e.lastSnapshotTime.Add(snapshotInterval)) {
		e.lastSnapshotTime = time.Now()
		// 必须保证每个组件MakeSnapshot时没有线程安全问题
		s := Snapshot{
			processor: e.processor.MakeSnapshot(),
			ppu:       e.ppu.MakeSnapshot(),
			bus:       e.bus.MakeSnapshot(),
			timestamp: e.lastSnapshotTime,
		}
		e.m.Lock()
		defer e.m.Unlock()
		e.snapshots = append(e.snapshots, s)
		if len(e.snapshots) > maxSnapshots {
			e.snapshots = e.snapshots[1:]
		}
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
		e.processor.Reverse(last.processor)
		revFrame := e.ppu.Reverse(last.ppu)
		e.bus.Reverse(last.bus)
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
