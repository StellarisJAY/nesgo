package cartridge

import (
	"log"
)

// Mapper0 not switchable prg and chr rom
type Mapper0 struct {
	program   []byte // program 程序代码
	Chr       []byte // Chr 图形数据
	mirroring byte   // mirroring
	PrgRAM    []byte
	trainer   []byte
	ChrRAM    bool
}

func newMapper0(raw []byte, mirroring byte) *Mapper0 {
	var prgRAM []byte
	if raw[6]&0b10 != 0 {
		prgRAM = make([]byte, 0x1000)
	}
	trainer, prgROM, chrROM, chrRAM := loadPrgAndChrROM(raw)
	var chr []byte
	if chrRAM {
		chr = make([]byte, 0x2000)
	} else {
		chr = chrROM
	}
	return &Mapper0{
		program:   prgROM,
		Chr:       chr,
		mirroring: mirroring,
		trainer:   trainer,
		PrgRAM:    prgRAM,
		ChrRAM:    chrRAM,
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
		log.Printf("attempt to write mapper-0 prg ROM at 0x%x\n", addr)
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
	if len(r.PrgRAM) == 0 {
		return 0
	}
	addr = addr - 0x6000
	return r.PrgRAM[addr]
}

func (r *Mapper0) writePrgRAM(addr uint16, val byte) {
	if len(r.PrgRAM) != 0 {
		addr = addr - 0x6000
		r.PrgRAM[addr] = val
	}
}

func (r *Mapper0) GetMirroring() byte {
	return r.mirroring
}

func (r *Mapper0) GetChrROM() []byte {
	return r.Chr
}

func (r *Mapper0) GetChrBank(bank byte) []byte {
	offset := uint16(bank) * 0x1000
	return r.Chr[offset : offset+4096]
}

func (r *Mapper0) WriteCHR(addr uint16, val byte) {
	if r.ChrRAM {
		r.Chr[addr] = val
	}
}
