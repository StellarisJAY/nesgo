package cartridge

type Mapper004 struct {
	raw []byte

	prgStart  uint32
	chrStart  uint32
	mirroring byte

	prgBanks [4][]byte // prgBanks 4个8KiB的prg rom banks
	chrBanks [8][]byte // chrBanks 8个1KiB的chr rom banks 0x0000~0x1fff

	bankSelect mapper004BankSelect

	prgRAM           []byte
	enablePrgRAM     bool
	allowPrgRAMWrite bool

	irqEnabled bool
}

type mapper004BankSelect byte

const (
	mapper4BankSelPrgMode byte = 1 << 6
	mapper4BankSelChrMode byte = 1 << 7

	mapper4RAMWriteProtect byte = 1 << 6
	mapper4RAMEnable       byte = 1 << 7
)

func NewMapper004(raw []byte, mirroring byte) *Mapper004 {
	_, prgStart, chrStart, _ := splitPrgAndChr(raw)
	lastPrg := chrStart - 0x2000
	chrBanks := [8][]byte{}
	offset := chrStart
	for i := 0; i < 8; i++ {
		chrBanks[i] = raw[offset : offset+1024]
		offset += 1024
	}
	return &Mapper004{
		raw:       raw,
		mirroring: mirroring,
		prgStart:  prgStart,
		chrStart:  chrStart,
		prgBanks: [4][]byte{
			nil,
			nil,
			nil,
			raw[lastPrg : lastPrg+0x2000], // prg rom 0xe000~0xffff 固定最后一个bank
		},
		chrBanks: chrBanks,
		prgRAM:   make([]byte, 0x2000),
	}
}

func (m *Mapper004) Read(addr uint16) byte {
	switch {
	case addr >= 0x6000 && addr <= 0x7fff:
		return m.readPrgRAM(addr)
	case addr >= 0x8000:
		addr = addr - 0x8000
		bank, off := addr/0x2000, addr%0x2000
		return m.prgBanks[bank][off]
	}
	return 0
}

func (m *Mapper004) Write(addr uint16, val byte) {
	switch {
	case addr >= 0x6000 && addr <= 0x7fff:
		m.writePrgRAM(addr, val)
	case addr >= 0x8000 && addr <= 0x9fff:
		// even: bank select, odd: bank data
		m.writeEvenOrOdd(addr, val, m.writeBankSelect, m.writeBankData)
	case addr >= 0xa000 && addr <= 0xbfff:
		// even: mirroring, odd: ram protect
		m.writeEvenOrOdd(addr, val, m.writeMirroring, m.writePrgRAMProtect)
	case addr >= 0xc000 && addr <= 0xefff:
		// even: irq latch, odd: irq reload
		m.writeEvenOrOdd(addr, val, m.writeIRQLatch, m.writeIRQReload)
	case addr >= 0xe000:
		// even: disable irq, odd: enable irq
		m.writeEvenOrOdd(addr, val, m.disableIRQ, m.enableIRQ)
	}
}

func (m *Mapper004) GetMirroring() byte {
	return m.mirroring
}

func (m *Mapper004) GetChrBank(bank byte) []byte {
	idx := bank * 4
	data := make([]byte, 0x1000)
	offset := 0
	var i byte
	for i = 0; i < 4; i++ {
		copy(data[offset:], m.chrBanks[idx+i])
		offset += 1024
	}
	return data
}

func (m *Mapper004) WriteCHR(addr uint16, val byte) {
	//panic("can't write chr")
}

func (m *Mapper004) writePrgRAM(addr uint16, val byte) {
	if m.enablePrgRAM && m.allowPrgRAMWrite {
		m.prgRAM[addr-0x6000] = val
	}
}

func (m *Mapper004) readPrgRAM(addr uint16) byte {
	if m.enablePrgRAM {
		return m.prgRAM[addr-0x6000]
	}
	return 0
}

func (m *Mapper004) writeEvenOrOdd(addr uint16, val byte, even, odd func(byte)) {
	switch addr & 1 {
	case 0:
		even(val)
	case 1:
		odd(val)
	}
}

func (m *Mapper004) writeBankSelect(val byte) {
	m.bankSelect = mapper004BankSelect(val)
	mode0 := byte(m.bankSelect)&mapper4BankSelPrgMode == 0
	secondLast := m.chrStart - 0x4000
	secondLastData := m.raw[secondLast : secondLast+0x2000]
	if mode0 {
		// 固定 0xc000 -> (-2)
		m.prgBanks[2] = secondLastData
	} else {
		// 固定 0x8000 -> (-2)
		m.prgBanks[0] = secondLastData
	}
}

func (m *Mapper004) writeBankData(val byte) {
	s := m.bankSelect & 0b111
	switch {
	case s == 0 || s == 1: // 切换2KiB chr
		m.switchChr2K(byte(s), val)
	case s >= 2 && s <= 5: // 切换1KiB chr
		m.switchChr1K(byte(s), val)
	case s == 6 || s == 7: // 切换 prg
		m.switchPrgBank(byte(s), val)
	default:
		panic("impossible")
	}
}

func (m *Mapper004) switchPrgBank(sel byte, bank byte) {
	mode0 := byte(m.bankSelect)&mapper4BankSelPrgMode == 0
	secondLast := m.chrStart - 0x4000
	// 忽略最高2位
	bank = bank & 0b11111
	start := uint32(bank)*0x2000 + m.prgStart
	secondLastData := m.raw[secondLast : secondLast+0x2000]
	switchedData := m.raw[start : start+0x2000]
	if mode0 {
		// 固定 0xc000 -> (-2)
		m.prgBanks[2] = secondLastData
		if sel == 6 {
			// 切换 0x8000
			m.prgBanks[0] = switchedData
		}
	} else {
		// 固定 0x8000 -> (-2)
		m.prgBanks[0] = secondLastData
		if sel == 6 {
			// 切换 0xc000
			m.prgBanks[2] = switchedData
		}
	}
	// 切换 0xa000
	if sel == 7 {
		m.prgBanks[1] = switchedData
	}
}

func (m *Mapper004) switchChr2K(sel byte, bank byte) {
	mode0 := byte(m.bankSelect)&mapper4BankSelChrMode == 0
	// ignore bottom bit
	start := uint32(bank&0b11111110)*1024 + m.chrStart
	first1K := m.raw[start : start+1024]
	second1K := m.raw[start+1024 : start+2048]
	if mode0 {
		switch sel {
		case 0:
			// 切换0x000~0x07ff
			m.chrBanks[0] = first1K
			m.chrBanks[1] = second1K
		case 1:
			// 切换0x0800~0x0fff
			m.chrBanks[2] = first1K
			m.chrBanks[3] = second1K
		}
	} else {
		switch sel {
		case 0:
			// 切换0x1000~0x17ff
			m.chrBanks[4] = first1K
			m.chrBanks[5] = second1K
		case 1:
			// 切换0x1800~0x1fff
			m.chrBanks[6] = first1K
			m.chrBanks[7] = second1K
		}
	}
}

func (m *Mapper004) switchChr1K(sel byte, bank byte) {
	mode0 := byte(m.bankSelect)&mapper4BankSelChrMode == 0
	start := uint32(bank)*1024 + m.chrStart
	data := m.raw[start : start+1024]
	idx := 0
	if mode0 {
		idx += 4
	}
	switch sel {
	case 2:
		m.chrBanks[idx] = data
	case 3:
		m.chrBanks[idx+1] = data
	case 4:
		m.chrBanks[idx+2] = data
	case 5:
		m.chrBanks[idx+3] = data
	}
}

func (m *Mapper004) writeMirroring(val byte) {
	m.mirroring = val & 1
}

func (m *Mapper004) writePrgRAMProtect(val byte) {
	m.enablePrgRAM = val&mapper4RAMEnable != 0
	m.allowPrgRAMWrite = val&mapper4RAMWriteProtect == 0
}

func (m *Mapper004) writeIRQLatch(val byte) {

}

func (m *Mapper004) writeIRQReload(val byte) {

}

func (m *Mapper004) enableIRQ(_ byte) {
	m.irqEnabled = true
}

func (m *Mapper004) disableIRQ(_ byte) {
	m.irqEnabled = false
}
