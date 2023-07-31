package cartridge

import (
	"fmt"
)

// Mapper0 not switchable prg and chr rom
type Mapper0 struct {
	program          []byte // program 程序代码
	chr              []byte // chr 图形数据
	mirroring        byte   // mirroring
	batteryBackedRAM []byte
	trainer          []byte
}

func newMapper0(raw []byte, mirroring byte) *Mapper0 {
	var prgRAM []byte
	if raw[6]&0b10 != 0 {
		prgRAM = make([]byte, 0x1000)
	}
	trainer, prgROM, chrROM := loadPrgAndChrROM(raw)
	return &Mapper0{
		program:          prgROM,
		chr:              chrROM,
		mirroring:        mirroring,
		trainer:          trainer,
		batteryBackedRAM: prgRAM,
	}
}

func (r *Mapper0) Read(addr uint16) byte {
	switch {
	case addr >= 0x6000 && addr <= 0x7fff:
		if len(r.trainer) != 0 && addr >= 0x7000 && addr <= 0x71ff {
			return r.readTrainer(addr)
		} else {
			return r.readPrgRAM(addr)
		}
	case addr >= 0x8000:
		return r.readProgramROM8(addr)
	default:
		return 0
	}
}

func (r *Mapper0) Write(addr uint16, val byte) {
	switch {
	case addr >= 0x6000 && addr <= 0x7fff:
		r.writePrgRAM(addr, val)
	case addr >= 0x8000:
		panic(fmt.Errorf("attempt to write mapper-0 prg ROM at 0x%x", addr))
	default:
		return
	}
}

func (r *Mapper0) readProgramROM8(addr uint16) byte {
	addr = addr - 0x8000
	if len(r.program) == 0x4000 && addr >= 0x4000 {
		addr = addr % 0x4000
	}
	return r.program[addr]
}

func (r *Mapper0) readTrainer(addr uint16) byte {
	addr = addr - 0x7000
	return r.trainer[addr]
}

func (r *Mapper0) readPrgRAM(addr uint16) byte {
	if len(r.batteryBackedRAM) == 0 {
		return 0
	}
	addr = addr - 0x6000
	return r.batteryBackedRAM[addr]
}

func (r *Mapper0) writePrgRAM(addr uint16, val byte) {
	if len(r.batteryBackedRAM) != 0 {
		addr = addr - 0x6000
		r.batteryBackedRAM[addr] = val
	}
}

func (r *Mapper0) GetMirroring() byte {
	return r.mirroring
}

func (r *Mapper0) GetChrROM() []byte {
	return r.chr
}

func (r *Mapper0) GetChrBank(bank byte) []byte {
	offset := uint16(bank) * 0x1000
	return r.chr[offset : offset+4096]
}
