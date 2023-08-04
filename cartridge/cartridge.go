package cartridge

import "log"

var NES = [4]byte{0x4E, 0x45, 0x53, 0x1A}

const (
	Vertical byte = iota
	Horizontal
	FourScreen
	OneScreenLow
	OneScreenHigh

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
	WriteCHR(addr uint16, val byte)
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
	case 2:
		return NewMapper002(raw, mirroring)
	case 3:
		return NewMapper003(raw, mirroring)
	case 4:
		return NewMapper004(raw, mirroring)
	default:
		panic("unsupported mapper")
	}
}

// loadPrgAndChrROM 分割prg和chr rom，返回trainer，prg，chr
func loadPrgAndChrROM(raw []byte) (trainer, prg, chr []byte, chrRAM bool) {
	// program and chr offsets
	hasTrainer := raw[6]&0b100 != 0
	programSize := uint16(raw[4]) * ProgramPageSize
	chrSize := uint16(raw[5]) * CHRPageSize
	var programStart uint16 = HeaderSize
	if hasTrainer {
		programStart += 512
		trainer = raw[HeaderSize : HeaderSize+512]
	}
	chrStart := programStart + programSize
	prg = raw[programStart : programStart+programSize]
	log.Printf("prg rom: %d KiB", programSize>>10)
	if chrSize == 0 {
		log.Println("using chr RAM")
		chrRAM = true
		return
	} else {
		log.Printf("chr rom: %d KiB", chrSize>>10)
		chr = raw[chrStart : chrStart+chrSize]
		return
	}
}

func splitPrgAndChr(raw []byte) (trainerStart, prgStart, chrStart uint32, chrRAM bool) {
	// program and chr offsets
	hasTrainer := raw[6]&0b100 != 0
	programSize := uint32(raw[4]) * ProgramPageSize
	chrSize := uint32(raw[5]) * CHRPageSize
	prgStart = HeaderSize
	if hasTrainer {
		prgStart += 512
	}
	chrStart = prgStart + programSize
	chrRAM = chrSize == 0
	log.Printf("prg rom: %d KiB", programSize>>10)
	if chrRAM {
		log.Println("using chr RAM")
	} else {
		log.Printf("chr rom: %d KiB", chrSize>>10)
	}

	return
}
