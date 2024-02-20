package bus

import (
	"github.com/stellarisJAY/nesgo/apu"
	"github.com/stellarisJAY/nesgo/cartridge"
	"github.com/stellarisJAY/nesgo/ppu"
	"time"
)

const (
	RAMSize                = 2048
	CpuRAMEnd              = 0x1FFF
	CpuRAMMask     uint16  = 0x7FF
	PPURegisterEnd         = 0x3FFF
	CPUFrequency           = 1790000 // Approximate CPUFrequency 1.79MHz see:https://www.nesdev.org/wiki/CPU
	CPUMaxBoost    float64 = 5.0
)

type RenderCallback func(*ppu.PPU)

// Bus 虚拟总线，CPU通过总线地址访问RAM, PPU, Registers
type Bus struct {
	cpuRAM         [RAMSize]byte       // cpuRAM cpu RAM内存区域
	resetPC        uint16              // resetPC 重启pc
	cartridge      cartridge.Cartridge // cartridge PrgROM和ChrROM
	ppu            *ppu.PPU            // ppu 图形处理器
	cycles         uint64              // cycles 总线时钟周期数，用来同步CPU和PPU的周期
	renderCallback RenderCallback
	joyPad         *JoyPad

	lastRenderCycles uint64
	lastRenderTime   time.Time
	cpuBoost         float64

	apu *apu.BasicAPU
}

type Snapshot struct {
	CpuRAM          [RAMSize]byte
	Cycles          uint64
	LastRenderCycle uint64
	CpuBoost        float64
}

// NewBus 创建总线，并将PPU和ROM接入总线
func NewBus(cartridge cartridge.Cartridge, ppu *ppu.PPU, callback RenderCallback, joyPad *JoyPad, apu *apu.BasicAPU) *Bus {
	return &Bus{
		cpuRAM:         [2048]byte{},
		cartridge:      cartridge,
		ppu:            ppu,
		renderCallback: callback,
		joyPad:         joyPad,
		cpuBoost:       1.0,
		apu:            apu,
	}
}

func (b *Bus) MakeSnapshot() Snapshot {
	var ram [RAMSize]byte
	copy(ram[:], b.cpuRAM[:])
	return Snapshot{
		CpuRAM:          ram,
		Cycles:          b.cycles,
		LastRenderCycle: b.lastRenderCycles,
		CpuBoost:        b.cpuBoost,
	}
}

func (b *Bus) Reverse(s Snapshot) {
	b.cpuRAM = s.CpuRAM
	b.cycles = s.Cycles
	b.lastRenderCycles = s.LastRenderCycle
	b.cpuBoost = s.CpuBoost
}

func (b *Bus) Tick(cycles uint64) {
	b.cycles += cycles
	before := b.ppu.PeekInterrupt()
	// ppu的cycles是CPU的三倍
	b.ppu.Tick(cycles * 3)
	if !before && b.ppu.PeekInterrupt() {
		start := time.Now()
		b.renderCallback(b.ppu)
		renderTime := time.Now().Sub(start)
		processorTime := start.Sub(b.lastRenderTime)
		// 两次渲染之间的cpu cycles除以CPU频率等于帧之间的间隔时间, nanoseconds
		freq := uint64(CPUFrequency * b.cpuBoost)
		frameTime := (b.cycles - b.lastRenderCycles) * 1000_000_000 / freq
		sleepTime := time.Duration(frameTime - uint64(renderTime.Nanoseconds()) - uint64(processorTime.Nanoseconds()))
		time.Sleep(sleepTime)
		b.lastRenderCycles = b.cycles
		b.lastRenderTime = time.Now()
	}
}

func (b *Bus) BoostCPU(delta float64) float64 {
	return b.SetCPUBoostRate(b.cpuBoost + delta)
}

func (b *Bus) SetCPUBoostRate(rate float64) float64 {
	b.cpuBoost = min(CPUMaxBoost, max(1.0, rate))
	return b.cpuBoost
}

func (b *Bus) ReadMemUint8(addr uint16) byte {
	switch {
	case addr <= CpuRAMEnd: // cpu内存
		addr = addr % 2048
		return b.readRAM8(addr)
	case addr == 0x2000 || addr == 0x2001 || addr == 0x2003 || addr == 0x2005 || addr == 0x2006 || addr == 0x4014: // 禁止读取ppu寄存器
		return 0
	case addr == 0x2002: // ppu 状态寄存器
		return b.ppu.ReadStatus()
	case addr == 0x2004: // oam data
		return b.ppu.ReadOam()
	case addr == 0x2007: // ppu读请求
		return b.ppu.ReadData()
	case addr <= PPURegisterEnd: // ppu寄存器mirroring
		addr = addr & 0x2007
		return b.ReadMemUint8(addr)
	case addr >= 0x4000 && addr <= 0x4015:
	case addr == 0x4016:
		return b.joyPad.read()
	case addr == 0x4017:
		return 0
	case addr >= 0x6000:
		return b.cartridge.Read(addr)
	default:
	}
	return 0
}

func (b *Bus) WriteMemUint8(addr uint16, val byte) {
	switch {
	case addr <= CpuRAMEnd: // CPU RAM
		addr = addr % 2048
		b.writeRAM8(addr, val)
	case addr == 0x2000: // PPU 状态寄存器
		b.ppu.WriteControl(val)
	case addr == 0x2001:
		b.ppu.WriteMask(val)
	case addr == 0x2002:
	case addr == 0x2003:
		b.ppu.WriteOamAddr(val)
	case addr == 0x2004:
		b.ppu.WriteOam(val)
	case addr == 0x2005:
		b.ppu.WriteScroll(val)
	case addr == 0x2006: // PPU请求地址
		b.ppu.WriteAddrReg(val)
	case addr == 0x2007: // 发起写请求
		b.ppu.WriteData(val)
	case addr <= PPURegisterEnd: // mirror ppu registers
		addr = addr & 0x2007
		b.WriteMemUint8(addr, val)
	case (addr >= 0x4000 && addr <= 0x4013) || addr == 0x4015: // apu
		b.apu.Write(addr, val)
	case addr == 0x4014: // oam dma
		buffer := make([]byte, 256)
		base := uint16(val) << 8
		for i := 0; i < 256; i++ {
			buffer[i] = b.ReadMemUint8(base + uint16(i))
		}
		b.ppu.WriteOamDMA(buffer)
	case addr == 0x4016: // joyPad 1
		b.joyPad.write(val)
	case addr == 0x4017: // joyPad 2
		//skip joyPad 2
	case addr >= 0x6000:
		b.cartridge.Write(addr, val)
	default:
	}
}

func (b *Bus) PollNMIInterrupt() bool {
	return b.ppu.PollInterrupt()
}

func (b *Bus) GetRAMRange(start, end uint16) []byte {
	startAddr, endAddr := start&CpuRAMMask, end&CpuRAMMask
	return b.cpuRAM[startAddr:endAddr]
}

func (b *Bus) WriteRAM(addr uint16, data []byte) {
	copy(b.cpuRAM[addr&CpuRAMMask:], data)
}

func (b *Bus) readRAM8(addr uint16) byte {
	return b.cpuRAM[addr]
}
func (b *Bus) readRAM16(addr uint16) uint16 {
	low := b.cpuRAM[addr]
	high := b.cpuRAM[addr+1]
	return uint16(high)<<8 + uint16(low)
}

func (b *Bus) writeRAM8(addr uint16, val byte) {
	b.cpuRAM[addr] = val
}

func (b *Bus) writeRAM16(addr uint16, val uint16) {
	low := byte(val & 0xFF)
	high := byte(val >> 8)
	b.cpuRAM[addr] = low
	b.cpuRAM[addr+1] = high
}

func (b *Bus) Cycles() uint64 {
	return b.cycles
}
