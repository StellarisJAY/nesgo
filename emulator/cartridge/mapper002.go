package cartridge

type Mapper002 struct {
	raw       []byte
	mirroring byte
	prgStart  uint32
	prgBanks  [2][]byte
	chrStart  uint32
	chrRAM    bool
}

func NewMapper002(raw []byte, mirroring byte) *Mapper002 {
	_, prgStart, chrStart, chrRAM := splitPrgAndChr(raw)
	lastBank := chrStart - 0x4000
	return &Mapper002{
		raw:       raw,
		mirroring: mirroring,
		prgStart:  prgStart,
		prgBanks: [2][]byte{
			raw[prgStart : prgStart+0x4000], // switchable 16KiB bank
			raw[lastBank : lastBank+0x4000], // fix to last bank
		},
		chrStart: chrStart,
		chrRAM:   chrRAM,
	}
}

func (m *Mapper002) Read(addr uint16) byte {
	switch {
	case addr >= 0x6000 && addr <= 0x7fff:
	case addr >= 0x8000:
		addr = addr - 0x8000
		bank := addr / 0x4000
		return m.prgBanks[bank][addr%0x4000]
	}
	return 0
}

func (m *Mapper002) Write(addr uint16, val byte) {
	switch {
	case addr >= 0x6000 && addr <= 0x7fff:
	case addr >= 0x8000:
		// ignore top 4 bits
		bank := uint32(val & 0x0f)
		// switch first 16KiB bank
		start := bank*0x4000 + m.prgStart
		m.prgBanks[0] = m.raw[start : start+0x4000]
	}
}

func (m *Mapper002) GetMirroring() byte {
	return m.mirroring
}

func (m *Mapper002) GetChrBank(bank byte) []byte {
	offset := uint32(bank)*0x1000 + m.chrStart
	return m.raw[offset : offset+0x1000]
}

func (m *Mapper002) WriteCHR(addr uint16, val byte) {
	if m.chrRAM {
		bank, off := addr/0x1000, addr%0x1000
		chr := m.GetChrBank(byte(bank))
		chr[off] = val
	}
}
