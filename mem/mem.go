package mem

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/ppu"
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
	ppu     *ppu.PPU
}

func NewBus(rom *ROM) *Bus {
	return &Bus{[2048]byte{}, 0, rom, ppu.NewPPU(rom.chr, byte(rom.mirroring))}
}

func NewBusWithNoROM() *Bus {
	return NewBus(EmptyROM())
}

func (b *Bus) ReadMemUint8(addr uint16) byte {
	switch {
	case addr <= CpuRAMEnd: // cpu内存
		addr = addr & CpuRAMMask
		return b.readRAM8(addr)
	case addr == 0x2000 || addr == 0x2001 || addr == 0x2003 || addr == 0x2005 || addr == 0x2006: // 禁止读取ppu寄存器
		panic("can't read write-only registers")
	case addr == 0x2007: // ppu读请求
		return b.ppu.ReadData()
	case addr <= PPURegisterEnd: // ppu寄存器mirroring
		addr = addr & PPURegisterMask
		return b.ReadMemUint8(addr)
	case addr >= ROMStart: // ROM
		return b.rom.readProgramROM8(addr)
	default:
		panic(fmt.Errorf("invalid memory addr: 0x%x", addr))
	}
	return 0
}
func (b *Bus) ReadMemUint16(addr uint16) uint16 {
	switch {
	case addr <= CpuRAMEnd: // CPU RAM
		addr = addr & CpuRAMMask
		return b.readRAM16(addr)
	case addr == 0x2000 || addr == 0x2001 || addr == 0x2003 || addr == 0x2005 || addr == 0x2006: // 禁止读取ppu寄存器
		panic("can't read write-only registers")
	case addr == 0x2007: // ppu读请求
		return uint16(b.ppu.ReadData())
	case addr <= PPURegisterEnd: // ppu寄存器mirroring
		addr = addr & PPURegisterMask
		return b.ReadMemUint16(addr)
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
	case addr <= CpuRAMEnd: // CPU RAM
		addr = addr & CpuRAMMask
		b.writeRAM8(addr, val)
	case addr == 0x2000: // PPU 状态寄存器
		b.ppu.WriteControl(val)
	case addr == 0x2006: // PPU请求地址
		b.ppu.WriteAddrReg(val)
	case addr == 0x2007: // 发起写请求
		b.ppu.WriteData(val)
	case addr <= PPURegisterEnd:
		addr = addr & PPURegisterMask
		b.WriteMemUint8(addr, val)
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
	case addr == 0x2000: // PPU 状态寄存器
		b.ppu.WriteControl(byte(val))
	case addr == 0x2006: // PPU请求地址
		b.ppu.WriteAddrReg(byte(val))
	case addr == 0x2007: // 发起写请求
		b.ppu.WriteData(byte(val))
	case addr <= PPURegisterEnd:
		addr = addr & PPURegisterMask
		b.WriteMemUint8(addr, byte(val))
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
