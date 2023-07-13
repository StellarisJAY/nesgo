package ppu

type Frame struct {
	data []byte
}

const (
	WIDTH  = 256
	HEIGHT = 240
)

func NewFrame() *Frame {
	return &Frame{
		make([]byte, WIDTH*HEIGHT*3),
	}
}

func (f *Frame) setPixel(x, y uint16, color Color) {
	first := y*3*WIDTH + x*3
	if first+2 < uint16(len(f.data)) {
		f.data[first] = color.R
		f.data[first+1] = color.G
		f.data[first+2] = color.B
	}
}

func (f *Frame) Data() []byte {
	return f.data
}
