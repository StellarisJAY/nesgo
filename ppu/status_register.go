package ppu

type StatusRegister struct {
	Val byte
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
	s.Val = s.Val & (^VBlankStarted)
}

func (s *StatusRegister) resetSprite0Hit() {
	s.Val = s.Val & (^Sprite0Hit)
}

func (s *StatusRegister) setVBlankStarted() {
	s.Val = s.Val | VBlankStarted
}

func (s *StatusRegister) isVBlank() bool {
	return s.Val&VBlankStarted != 0
}

func (s *StatusRegister) setSprite0Hit() {
	s.Val = s.Val | Sprite0Hit
}
