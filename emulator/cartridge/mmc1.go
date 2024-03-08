package cartridge

// Mapper1 支持切换chr和prg
type Mapper1 struct {
	Raw      []byte
	PrgStart uint32
	ChrStart uint32

	PrgRAM      []byte
	ShiftReg    Mapper0Shift
	CtrlReg     byte
	PrgBankReg  byte
	PrgBanks    [2][]byte
	ChrBankRegs [2]byte
	ChrBanks    [2][]byte
	ChrRAM      bool
}

type Mapper0Shift struct {
	Val   byte
	Index byte
}

const (
	ctrlRegPrgBankMode byte = 2
	ctrlRegChrBankMode byte = 4
)

func NewMapper1(raw []byte) *Mapper1 {
	var prgRAM []byte
	if raw[6]&0b10 != 0 {
		prgRAM = make([]byte, 0x1000)
	}
	_, prgStart, chrStart, chrRAM := splitPrgAndChr(raw)
	var chrBanks [2][]byte
	if chrRAM {
		chrBanks = [2][]byte{
			make([]byte, 0x1000),
			make([]byte, 0x1000),
		}
	} else {
		chrBanks = [2][]byte{
			raw[chrStart : chrStart+0x1000],
			raw[chrStart+0x1000 : chrStart+0x2000],
		}
	}

	return &Mapper1{
		PrgRAM:   prgRAM,
		Raw:      raw,
		PrgStart: prgStart,
		ChrStart: chrStart,
		PrgBanks: [2][]byte{
			raw[prgStart : prgStart+0x4000],
			raw[chrStart-0x4000 : chrStart],
		},
		ChrBanks: chrBanks,
		ChrRAM:   chrRAM,
	}
}

func (m *Mapper1) Read(addr uint16) byte {
	switch {
	case addr >= 0x6000 && addr <= 0x7fff:
		return m.readPrgRAM(addr)
	case addr >= 0x8000:
		return m.readPrgROM(addr)
	default:
		return 0
	}
}

// Write cpu地址映射的寄存器
func (m *Mapper1) Write(addr uint16, val byte) {
	switch {
	case addr >= 0x6000 && addr <= 0x7fff: // prg RAM
		m.writePrgRAM(addr, val)
	case addr >= 0x8000: // shift register
		m.writeShiftRegister(addr, val)
	}
}

func (m *Mapper1) GetMirroring() byte {
	switch m.CtrlReg & 0b11 {
	case 0:
		return OneScreenLow
	case 1:
		return OneScreenHigh
	case 2:
		return Vertical
	case 3:
		return Horizontal
	default:
		panic("not possible")
	}
}

func (m *Mapper1) GetChrBank(bank byte) []byte {
	return m.ChrBanks[bank]
}

func (m *Mapper1) WriteCHR(addr uint16, val byte) {
	if m.ChrRAM {
		bank, off := addr/0x1000, addr%0x1000
		m.ChrBanks[bank][off] = val
	}
}

func (m *Mapper1) writePrgRAM(addr uint16, val byte) {
	if len(m.PrgRAM) > 0 {
		m.PrgRAM[(addr-0x6000)%0x1000] = val
	}
}

func (m *Mapper1) readPrgRAM(addr uint16) byte {
	if len(m.PrgRAM) > 0 {
		return m.PrgRAM[(addr-0x6000)%0x1000]
	}
	return 0
}

func (m *Mapper1) readPrgROM(addr uint16) byte {
	addr = addr - 0x8000
	bank, offset := addr/0x4000, addr%0x4000
	return m.PrgBanks[bank][offset]
}

// writeInternal 覆盖写内部地址映射的寄存器
func (m *Mapper1) writeInternal(addr uint16, val byte) {
	switch {
	case addr >= 0x8000 && addr <= 0x9fff: // control register
		m.CtrlReg = val
	case addr >= 0xa000 && addr <= 0xbfff: // switch chr bank0
		m.writeChrBankReg(0, val)
	case addr >= 0xc000 && addr <= 0xdfff: // switch chr bank1
		m.writeChrBankReg(1, val)
	case addr >= 0xe000: // prg bank
		m.writePrgBankReg(val)
	}
}

func (m *Mapper1) writeShiftRegister(addr uint16, val byte) {
	// bit 7 set, clear shift register
	if val&0x80 != 0 {
		m.ShiftReg.clear()
		m.CtrlReg = m.CtrlReg | 0x0c
	} else {
		// 第五次写shiftReg，shiftReg值写到addr对应的内部地址，并重置shiftReg
		if res, last := m.ShiftReg.write(val); last {
			m.writeInternal(addr, res)
		}
	}
}

func (sr *Mapper0Shift) clear() {
	sr.Index = 0
	sr.Val = 1
}

// write 写shiftReg，返回写入后的shiftReg值 和 是否是最后一次写入
func (sr *Mapper0Shift) write(val byte) (byte, bool) {
	if sr.Index <= 3 {
		sr.Val = (sr.Val >> 1) | ((val & 1) << 4)
		sr.Index++
		return sr.Val, false
	} else {
		result := (sr.Val >> 1) | ((val & 1) << 4)
		sr.clear()
		return result, true
	}
}

func (m *Mapper1) writeChrBankReg(bank byte, val byte) {
	m.ChrBankRegs[bank] = val
	mode := m.CtrlReg & (1 << ctrlRegChrBankMode)
	// 非8KiB模式，切换4KiB的单个chr bank
	if mode != 0 {
		offset := m.ChrStart + uint32(val&0b11111)*4096
		m.ChrBanks[bank] = m.Raw[offset : offset+4096]
	} else if bank == 0 {
		// 8KiB模式，切换两个chr banks, 忽略bit0
		offset := m.ChrStart + uint32(val&0xfe)*4096
		m.ChrBanks[0] = m.Raw[offset : offset+4096]
		m.ChrBanks[1] = m.Raw[offset+4096 : offset+8192]
	}
}

func (m *Mapper1) writePrgBankReg(val byte) {
	m.PrgBankReg = val
	mode := (m.CtrlReg & (0b11 << ctrlRegPrgBankMode)) >> ctrlRegPrgBankMode
	offset := m.PrgStart + uint32(val&0b1111)*0x4000
	switch mode {
	case 0:
		fallthrough
	case 1:
		start := m.PrgStart + uint32(val&0xfe)*0x4000
		m.PrgBanks[0] = m.Raw[start : start+0x4000]
		m.PrgBanks[1] = m.Raw[start+0x4000 : start+0x8000]
	case 2:
		m.PrgBanks[1] = m.Raw[offset : offset+0x4000]
		m.PrgBanks[0] = m.Raw[m.PrgStart : m.PrgStart+0x4000]
	case 3:
		m.PrgBanks[0] = m.Raw[offset : offset+0x4000]
		lastPage := m.ChrStart - 0x4000
		m.PrgBanks[1] = m.Raw[lastPage : lastPage+0x4000]
	default:
	}
}
