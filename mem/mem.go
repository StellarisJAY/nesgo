package mem

import (
	"fmt"
)

type Memory interface {
	ReadMemUint8(addr uint16) byte
	ReadMemUint16(addr uint16) uint16

	WriteMemUint8(addr uint16, val byte)
	WriteMemUint16(addr uint16, val uint16)
}

const (
	RAMSize           = 2048
	CpuRAMEnd         = 0x1FFF
	CpuRAMMask uint16 = 0x7FF

	PPURegisterEnd  = 0x3FFF
	PPURegisterMask = 0x2007

	ROMStart = 0x8000
)

// Bus 虚拟总线，CPU通过总线地址访问RAM, PPU, Registers
type Bus struct {
	cpuRAM  [RAMSize]byte // cpuRAM cpu RAM内存区域
	resetPC uint16        // resetPC 重启pc
	rom     *ROM
}

func NewBus(rom *ROM) *Bus {
	return &Bus{[2048]byte{}, 0, rom}
}

func NewBusWithNoROM() *Bus {
	return NewBus(EmptyROM())
}

func (b *Bus) ReadMemUint8(addr uint16) byte {
	switch {
	case addr <= CpuRAMEnd:
		addr = addr & CpuRAMMask
		return b.readRAM8(addr)
	case addr <= PPURegisterEnd:
		addr = addr & PPURegisterMask
		panic("ppu not available")
	case addr >= ROMStart:
		return b.rom.readProgramROM8(addr)
	default:
		panic(fmt.Errorf("invalid memory addr: 0x%x", addr))
	}
	return 0
}
func (b *Bus) ReadMemUint16(addr uint16) uint16 {
	switch {
	case addr <= CpuRAMEnd:
		addr = addr & CpuRAMMask
		return b.readRAM16(addr)
	case addr <= PPURegisterEnd:
		addr = addr & PPURegisterMask
		panic("ppu not available")
	case addr == 0xFFFC:
		return b.resetPC
	case addr >= ROMStart:
		return b.rom.readProgramROM16(addr)
	default:
		panic(fmt.Errorf("invalid memory addr: 0x%x", addr))
	}
	return 0
}

func (b *Bus) WriteMemUint8(addr uint16, val byte) {
	switch {
	case addr <= CpuRAMEnd:
		addr = addr & CpuRAMMask
		b.writeRAM8(addr, val)
	case addr <= PPURegisterEnd:
		addr = addr & PPURegisterMask
		panic("ppu not available")
	case addr >= ROMStart:
		panic(fmt.Errorf("can't write ROM addr: 0x%x", addr))
	default:
		panic(fmt.Errorf("invalid memory addr: 0x%x", addr))
	}
}
func (b *Bus) WriteMemUint16(addr uint16, val uint16) {
	switch {
	case addr <= CpuRAMEnd:
		addr = addr & CpuRAMMask
		b.writeRAM16(addr, val)
	case addr <= PPURegisterEnd:
		addr = addr & PPURegisterMask
		panic("ppu not available")
	case addr == 0xFFFC:
		b.resetPC = val
	case addr >= ROMStart:
		panic(fmt.Errorf("can't write ROM addr: 0x%x", addr))
	default:
		panic(fmt.Errorf("invalid memory addr: 0x%x", addr))
	}
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
