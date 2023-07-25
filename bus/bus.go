package bus

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/ppu"
)

const (
	RAMSize           = 2048
	CpuRAMEnd         = 0x1FFF
	CpuRAMMask uint16 = 0x7FF

	PPURegisterEnd  = 0x3FFF
	PPURegisterMask = 0x2007

	ROMStart = 0x8000
)

type RenderCallback func(*ppu.PPU)

// Bus 虚拟总线，CPU通过总线地址访问RAM, PPU, Registers
type Bus struct {
	cpuRAM         [RAMSize]byte // cpuRAM cpu RAM内存区域
	resetPC        uint16        // resetPC 重启pc
	rom            *ROM          // rom PrgROM和ChrROM
	ppu            *ppu.PPU      // ppu 图形处理器
	cycles         uint64        // cycles 总线时钟周期数，用来同步CPU和PPU的周期
	renderCallback RenderCallback
	joyPad         *JoyPad
}

// NewBus 创建总线，并将PPU和ROM接入总线
func NewBus(rom *ROM, ppu *ppu.PPU, callback RenderCallback, joyPad *JoyPad) *Bus {
	return &Bus{[2048]byte{}, 0, rom, ppu, 0, callback, joyPad}
}

func NewBusWithNoROM() *Bus {
	return NewBus(EmptyROM(), nil, nil, nil)
}

func (b *Bus) Tick(cycles uint64) {
	b.cycles += cycles
	nmiBefore := b.ppu.PeekInterrupt()
	// ppu的cycles是CPU的三倍
	b.ppu.Tick(cycles * 3)
	nmiAfter := b.ppu.PeekInterrupt()
	if !nmiBefore && nmiAfter {
		b.renderCallback(b.ppu)
	}
}

func (b *Bus) ReadMemUint8(addr uint16) byte {
	switch {
	case addr <= CpuRAMEnd: // cpu内存
		addr = addr & CpuRAMMask
		return b.readRAM8(addr)
	case addr == 0x2000 || addr == 0x2001 || addr == 0x2003 || addr == 0x2005 || addr == 0x2006: // 禁止读取ppu寄存器
		panic("can't read write-only registers")
	case addr == 0x2002: // ppu 状态寄存器
		return b.ppu.ReadStatus()
	case addr == 0x2007: // ppu读请求
		return b.ppu.ReadData()
	case addr <= PPURegisterEnd: // ppu寄存器mirroring
		addr = addr & 0b00100000_00000111
		return b.ReadMemUint8(addr)
	case addr >= 0x4000 && addr <= 0x4015:
	case addr == 0x4016:
		return b.joyPad.read()
	case addr == 0x4017:
		return b.joyPad.read()
	case addr >= ROMStart: // ROM
		return b.rom.readProgramROM8(addr)
	default:
		panic(fmt.Errorf("invalid memory addr: 0x%x", addr))
	}
	return 0
}

func (b *Bus) WriteMemUint8(addr uint16, val byte) {
	switch {
	case addr <= CpuRAMEnd: // CPU RAM
		addr = addr & CpuRAMMask
		b.writeRAM8(addr, val)
	case addr == 0x2000: // PPU 状态寄存器
		b.ppu.WriteControl(val)
	case addr == 0x2001:
		b.ppu.WriteMask(val)
	case addr == 0x2002:
		panic("can't write ppu status")
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
	case addr <= PPURegisterEnd:
		addr = addr & PPURegisterMask
		b.WriteMemUint8(addr, val)
	case (addr >= 0x4000 && addr <= 0x4013) || addr == 0x4015:
	case addr == 0x4014:
		buffer := make([]byte, 256)
		base := uint16(val) << 8
		for i := 0; i < 256; i++ {
			buffer[i] = b.ReadMemUint8(base + uint16(i))
		}
		b.ppu.WriteOamDMA(buffer)
	case addr == 0x4016:
		b.joyPad.write(val)
	case addr == 0x4017:
		b.joyPad.write(val)
	case addr >= ROMStart:
		panic(fmt.Errorf("can't write ROM addr: 0x%x", addr))
	default:
		panic(fmt.Errorf("invalid memory addr: 0x%x", addr))
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
