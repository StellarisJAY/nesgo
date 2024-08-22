package cartridge

import (
	"bytes"
	"encoding/gob"
	"errors"
)

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

var ErrUnknownCartridgeFormat = errors.New("unknown cartridge format")
var ErrUnsupportedMapper = errors.New("unsupported mapper type")

var mapperNames = []string{
	"mapper0", "mapper1", "mapper2", "mapper3", "mapper4",
}

type Cartridge interface {
	Read(addr uint16) byte
	Write(addr uint16, val byte)
	GetMirroring() byte
	// GetChrBank 获取bank编号对应的chr patternTable
	GetChrBank(bank byte) []byte
	WriteCHR(addr uint16, val byte)
}

type Info struct {
	Name      string
	Mapper    byte
	Mirroring byte
}

// MakeCartridge 从iNES文件读取cartridge
func MakeCartridge(raw []byte) (Cartridge, error) {
	c, _, err := makeCartridge(raw)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func ParseCartridgeInfo(raw []byte) (*Info, error) {
	_, info, err := makeCartridge(raw)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func makeCartridge(raw []byte) (Cartridge, *Info, error) {
	if len(raw) <= 8 {
		return nil, nil, ErrUnknownCartridgeFormat
	}
	if raw[0] != NES[0] || raw[1] != NES[1] || raw[2] != NES[2] || raw[3] != NES[3] {
		return nil, nil, ErrUnknownCartridgeFormat
	}
	// read mapperNumber
	mapperNumber := (raw[7] & 0xf0) | (raw[6] >> 4)
	// check version
	if version := (raw[7] >> 2) & 0b11; version != 0 {
		return nil, nil, ErrUnknownCartridgeFormat
	}
	info := new(Info)
	info.Mapper = mapperNumber
	// read Mirroring from raw[6] raw[7]
	if raw[6]&1 != 0 {
		info.Mirroring = Vertical
	} else {
		info.Mirroring = Horizontal
	}
	if raw[6]&0b1000 != 0 {
		info.Mirroring = FourScreen
	}
	// 创建mapper
	switch info.Mapper {
	case 0:
		return newMapper0(raw, info.Mirroring), info, nil
	case 1:
		return NewMapper1(raw), info, nil
	case 2:
		return NewMapper002(raw, info.Mirroring), info, nil
	case 3:
		return NewMapper003(raw, info.Mirroring), info, nil
	case 4:
		return NewMapper004(raw, info.Mirroring), info, nil
	default:
		return nil, info, ErrUnsupportedMapper
	}
}

func Save(c Cartridge) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(c); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Load(c Cartridge, data []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(c)
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
	if chrSize == 0 {
		chrRAM = true
		return
	} else {
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
	return
}

func MirroringToString(m byte) string {
	switch m {
	case Horizontal:
		return "Horizontal"
	case Vertical:
		return "Vertical"
	case FourScreen:
		return "FourScreen"
	default:
		return "Unknown"
	}
}

func MapperToString(m byte) string {
	if int(m) >= len(mapperNames) {
		return "Unsupported"
	}
	return mapperNames[m]
}
