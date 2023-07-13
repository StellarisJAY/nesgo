package ppu

type MaskRegister struct {
	val byte
}

func NewMaskRegister() MaskRegister {
	return MaskRegister{0}
}

func (m *MaskRegister) set(val byte) {
	m.val = val
}

func (m *MaskRegister) getBit(offset byte) bool {
	return m.val&offset != 0
}
