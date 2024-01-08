package ppu

type MaskRegister struct {
	Val byte
}

const (
	GreyScale byte = 1 << iota
	ShowBackground8
	ShowSprite8
	ShowBackground
	ShowSprite
	EmphasizeRed
	EmphasizeGreen
	EmphasizeBlue
)

func NewMaskRegister() MaskRegister {
	return MaskRegister{0}
}

func (m *MaskRegister) set(val byte) {
	m.Val = val
}

func (m *MaskRegister) getBit(offset byte) bool {
	return m.Val&offset != 0
}
