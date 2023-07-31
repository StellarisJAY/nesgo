package cartridge

import "log"

var NES = [4]byte{0x4E, 0x45, 0x53, 0x1A}

const (
	Vertical byte = iota
	Horizontal
	FourScreen
	OneScreen
	ProgramPageSize = 16 * 1024
	CHRPageSize     = 8 * 1024
	HeaderSize      = 16
)

type Cartridge interface {
	Read(addr uint16) byte
	Write(addr uint16, val byte)
	GetMirroring() byte
	// GetChrBank 获取bank编号对应的chr patternTable
	GetChrBank(bank byte) []byte
}

func EmptyMapper0() Cartridge {
	return &Mapper0{}
}

// MakeCartridge 从iNES文件读取cartridge
func MakeCartridge(raw []byte) Cartridge {
	if raw[0] != NES[0] || raw[1] != NES[1] || raw[2] != NES[2] || raw[3] != NES[3] {
		panic("invalid NES Header")
		return nil
	}
	// read mapperNumber
	mapperNumber := (raw[7] & 0xf0) | (raw[6] >> 4)
	log.Printf("mapper: %d\n", mapperNumber)
	// check version
	if version := (raw[7] >> 2) & 0b11; version != 0 {
		panic("NES version 2.0 not supported")
		return nil
	}
	// read Mirroring from raw[6] raw[7]
	var mirroring byte
	if raw[6]&1 != 0 {
		mirroring = Vertical
	} else {
		mirroring = Horizontal
	}
	if raw[6]&0b1000 != 0 {
		mirroring = FourScreen
	}
	// 创建mapper
	switch mapperNumber {
	case 0:
		return newMapper0(raw, mirroring)
	case 1:
		return NewMapper1(raw)
	default:
		panic("unsupported mapper")
	}
}

// loadPrgAndChrROM 分割prg和chr rom，返回trainer，prg，chr
func loadPrgAndChrROM(raw []byte) ([]byte, []byte, []byte) {
	// program and chr offsets
	hasTrainer := raw[6]&0b100 != 0
	programSize := uint16(raw[4]) * ProgramPageSize
	chrSize := uint16(raw[5]) * CHRPageSize
	var programStart uint16 = HeaderSize
	if hasTrainer {
		programStart += 512
	}
	chrStart := programStart + programSize
	log.Printf("prg rom: %d KiB", programSize>>10)
	log.Printf("chr rom: %d KiB", chrSize>>10)
	return raw[HeaderSize : HeaderSize+512], raw[programStart : programStart+programSize], raw[chrStart : chrStart+chrSize]
}

func splitPrgAndChr(raw []byte) (uint16, uint16, uint16) {
	// program and chr offsets
	hasTrainer := raw[6]&0b100 != 0
	programSize := uint16(raw[4]) * ProgramPageSize
	chrSize := uint16(raw[5]) * CHRPageSize
	var programStart uint16 = HeaderSize
	if hasTrainer {
		programStart += 512
	}
	chrStart := programStart + programSize
	log.Printf("prg rom: %d KiB", programSize>>10)
	log.Printf("chr rom: %d KiB", chrSize>>10)
	return HeaderSize, programStart, chrStart
}
