package mem

type Mirroring byte

const (
	Vertical Mirroring = iota
	Horizontal
	FourScreen
	ProgramPageSize = 16 * 1024
	CHRPageSize     = 8 * 1024
	HeaderSize      = 16
)

var NES = [4]byte{0x4E, 0x45, 0x53, 0x1A}

// ROM 只读内存
type ROM struct {
	program   []byte    // program 程序代码
	chr       []byte    // chr 图形数据
	mapper    byte      // mapper
	mirroring Mirroring // mirroring
}

func EmptyROM() *ROM {
	return &ROM{}
}

// NewROM 从NES数据创建ROM内存
func NewROM(raw []byte) *ROM {
	if raw[0] != NES[0] || raw[1] != NES[1] || raw[2] != NES[2] || raw[3] != NES[3] {
		panic("invalid NES Header")
		return nil
	}
	// read mapper
	mapper := (raw[7] & 0xF) | (raw[6] >> 4)
	// check version
	if version := (raw[7] >> 2) & 0b11; version != 0 {
		panic("NES version 2.0 not supported")
		return nil
	}
	// read Mirroring from raw[6] raw[7]
	var mirroring Mirroring
	if raw[6]&1 != 0 {
		mirroring = Vertical
	} else {
		mirroring = Horizontal
	}
	if raw[6]&0b1000 != 0 {
		mirroring = FourScreen
	}

	// program and chr offsets
	skipTrainer := raw[6]&0b10 == 0
	programSize := uint16(raw[4]) * ProgramPageSize
	chrSize := uint16(raw[5]) * CHRPageSize
	var programStart uint16 = HeaderSize
	if !skipTrainer {
		programStart += 512
	}
	chrStart := programStart + programSize
	program := make([]byte, programSize)
	chr := make([]byte, chrSize)
	copy(program, raw[programSize:programStart+programSize])
	copy(chr, raw[chrStart:chrStart+chrSize])

	return &ROM{
		program:   program,
		chr:       chr,
		mapper:    mapper,
		mirroring: mirroring,
	}
}

func (r *ROM) readProgramROM8(addr uint16) byte {
	addr = addr - 0x8000
	if len(r.program) == 0x4000 && addr >= 0x4000 {
		addr = addr % 0x4000
	}
	return r.program[addr]
}

func (r *ROM) readProgramROM16(addr uint16) uint16 {
	addr = addr - 0x8000
	if len(r.program) == 0x4000 && addr >= 0x4000 {
		addr = addr % 0x4000
	}
	low := r.program[addr]
	high := r.program[addr+1]
	return uint16(high)<<8 + uint16(low)
}
