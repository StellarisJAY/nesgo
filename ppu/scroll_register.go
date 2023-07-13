package ppu

type ScrollRegister struct {
	x     byte
	y     byte
	latch bool
}

func NewScrollRegister() ScrollRegister {
	return ScrollRegister{0, 0, false}
}

func (s *ScrollRegister) write(val byte) {
	if !s.latch {
		s.x = val
	} else {
		s.y = val
	}
	s.latch = !s.latch
}

func (s *ScrollRegister) resetLatch() {
	s.latch = false
}
