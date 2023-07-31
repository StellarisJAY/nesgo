package cartridge

// Mapper1 支持切换chr和prg
type Mapper1 struct {
	raw      []byte
	prgStart uint16
	chrStart uint16

	prgRAM      []byte
	shiftReg    mapper1ShiftRegister
	ctrlReg     byte
	prgBankReg  byte
	prgROM      [2][]byte
	chrBankRegs [2]byte
	chrBanks    [2][]byte
}

type mapper1ShiftRegister struct {
	val   byte
	index byte
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
	_, prgStart, chrStart := splitPrgAndChr(raw)
	return &Mapper1{
		prgRAM:   prgRAM,
		raw:      raw,
		prgStart: prgStart,
		chrStart: chrStart,
		prgROM: [2][]byte{
			raw[prgStart : prgStart+0x4000],
			raw[prgStart+0x4000 : prgStart+0x8000],
		},
		chrBanks: [2][]byte{
			raw[chrStart : chrStart+0x1000],
			raw[chrStart+0x1000 : chrStart+0x2000],
		},
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
	switch m.ctrlReg & 0b11 {
	case 0:
		fallthrough
	case 1:
		return OneScreen
	case 2:
		return Vertical
	case 3:
		return Horizontal
	default:
		panic("not possible")
	}
}

func (m *Mapper1) GetChrBank(bank byte) []byte {
	return m.chrBanks[bank]
}

func (m *Mapper1) writePrgRAM(addr uint16, val byte) {
	if len(m.prgRAM) > 0 {
		m.prgRAM[addr-0x6000] = val
	}
}

func (m *Mapper1) readPrgRAM(addr uint16) byte {
	if len(m.prgRAM) > 0 {
		return m.prgRAM[addr-0x6000]
	}
	return 0
}

func (m *Mapper1) readPrgROM(addr uint16) byte {
	addr = addr - 0x8000
	bank, offset := addr/0x4000, addr%0x4000
	return m.prgROM[bank][offset]
}

// writeInternal 覆盖写内部地址映射的寄存器
func (m *Mapper1) writeInternal(addr uint16, val byte) {
	switch {
	case addr >= 0x8000 && addr <= 0x9fff: // control register
		m.ctrlReg = val
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
		m.shiftReg.clear()
	} else {
		// 第五次写shiftReg，shiftReg值写到addr对应的内部地址，并重置shiftReg
		if res, last := m.shiftReg.write(val); last {
			m.writeInternal(addr, res)
		}
	}
}

func (sr *mapper1ShiftRegister) clear() {
	sr.index = 0
	sr.val = 1
}

// write 写shiftReg，返回写入后的shiftReg值 和 是否是最后一次写入
func (sr *mapper1ShiftRegister) write(val byte) (byte, bool) {
	if sr.index <= 3 {
		sr.val = (sr.val >> 1) | ((val & 1) << 4)
		sr.index++
		return sr.val, false
	} else {
		result := (sr.val >> 1) | ((val & 1) << 4)
		sr.clear()
		return result, true
	}
}

func (m *Mapper1) writeCtrl(val byte) {
	m.ctrlReg = val
}

func (m *Mapper1) writeChrBankReg(bank byte, val byte) {
	m.chrBankRegs[bank] = val
	mode := m.ctrlReg & (1 << ctrlRegChrBankMode)
	// 非8KiB模式，切换4KiB的单个chr bank
	if mode != 0 {
		offset := m.chrStart + uint16(val&0b11111)*4096
		m.chrBanks[bank] = m.raw[offset : offset+4096]
	} else {
		// 8KiB模式，切换两个chr banks, 忽略bit0
		offset := m.chrStart + uint16(val>>1)*8192
		m.chrBanks[0] = m.raw[offset : offset+4096]
		m.chrBanks[1] = m.raw[offset+4096 : offset+8192]
	}
}

func (m *Mapper1) writePrgBankReg(val byte) {
	m.prgBankReg = val
	mode := m.ctrlReg & (0b11 << ctrlRegPrgBankMode)
	if mode == 0 || mode == 1 {

		return
	}
	var switchIdx byte
	if mode == 2 {
		switchIdx = 1
	} else {
		switchIdx = 0
	}
	// switch 16KiB bank
	offset := m.prgStart + uint16(val&0xff)*0x4000
	m.prgROM[switchIdx] = m.raw[offset : offset+0x4000]
}
