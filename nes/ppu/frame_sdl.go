//go:build sdl

package ppu

import (
	"image/color"
	"slices"
)

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

func (f *Frame) setPixel(x, y uint32, c Color) {
	first := y*3*WIDTH + x*3
	if first+2 < uint32(len(f.data)) {
		f.data[first] = c.R
		f.data[first+1] = c.G
		f.data[first+2] = c.B
	}
}

func (f *Frame) getPixel(x, y uint32) color.Color {
	first := y*3*WIDTH + x*3
	return color.RGBA{
		R: f.data[first],
		G: f.data[first+1],
		B: f.data[first+2],
		A: 255,
	}
}

func (f *Frame) Clone() *Frame {
	return &Frame{
		data: slices.Clone(f.data),
	}
}

func (f *Frame) Data() []byte {
	return f.data
}
