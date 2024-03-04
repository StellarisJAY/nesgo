package ppu

type ScrollRegister struct {
	X     byte
	Y     byte
	Latch bool
}

func NewScrollRegister() ScrollRegister {
	return ScrollRegister{0, 0, false}
}

func (s *ScrollRegister) write(val byte) {
	if !s.Latch {
		s.X = val
	} else {
		s.Y = val
	}
	s.Latch = !s.Latch
}

func (s *ScrollRegister) resetLatch() {
	s.Latch = false
}

func (s *ScrollRegister) reset() {
	s.Latch = false
	s.X = 0
	s.Y = 0
}
