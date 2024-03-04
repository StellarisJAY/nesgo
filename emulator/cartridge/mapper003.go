package cartridge

type Mapper003 struct {
	raw       []byte
	mirroring byte

	prgStart uint32
	chrStart uint32

	chrBanks [2][]byte
}

func NewMapper003(raw []byte, mirroring byte) *Mapper003 {
	_, prgStart, chrStart, _ := splitPrgAndChr(raw)
	return &Mapper003{
		raw:       raw,
		mirroring: mirroring,
		prgStart:  prgStart,
		chrStart:  chrStart,
		chrBanks: [2][]byte{
			raw[chrStart : chrStart+0x1000],
			raw[chrStart+0x1000 : chrStart+0x2000],
		},
	}
}

func (m *Mapper003) Read(addr uint16) byte {
	switch {
	case addr >= 0x6000 && addr <= 0x7fff:
	case addr >= 0x8000:
		addr = addr - 0x8000
		return m.raw[m.prgStart+uint32(addr)]
	}
	return 0
}

func (m *Mapper003) Write(addr uint16, val byte) {
	switch {
	case addr >= 0x6000 && addr <= 0x7fff:
	case addr >= 0x8000:
		start := uint32(val&0b11)*0x2000 + m.chrStart
		m.chrBanks[0] = m.raw[start : start+0x1000]
		m.chrBanks[1] = m.raw[start+0x1000 : start+0x2000]
	}
}

func (m *Mapper003) GetMirroring() byte {
	return m.mirroring
}

func (m *Mapper003) GetChrBank(bank byte) []byte {
	return m.chrBanks[bank]
}

func (m *Mapper003) WriteCHR(addr uint16, val byte) {

}
