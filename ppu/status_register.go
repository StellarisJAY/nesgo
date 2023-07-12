package ppu

type StatusRegister struct {
	val byte
}

const (
	SpriteOverflow byte = 1 << 5
	Sprite0Hit     byte = 1 << 6
	VBlankStarted  byte = 1 << 7
)

func NewStatusRegister() StatusRegister {
	return StatusRegister{}
}

func (s *StatusRegister) resetVBlankStarted() {
	s.val = s.val & (^VBlankStarted)
}

func (s *StatusRegister) resetSprite0Hit() {
	s.val = s.val & (^Sprite0Hit)
}

func (s *StatusRegister) setVBlankStarted() {
	s.val = s.val | VBlankStarted
}
