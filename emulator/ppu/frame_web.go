//go:build web

package ppu

import (
	"image"
	"image/color"
)

type Frame struct {
	ycbcr *image.YCbCr
}

const (
	WIDTH  = 256
	HEIGHT = 240
)

func NewFrame() *Frame {
	return &Frame{
		image.NewYCbCr(image.Rect(0, 0, WIDTH, HEIGHT), image.YCbCrSubsampleRatio420),
	}
}

func (f *Frame) setPixel(x, y uint32, c Color) {
	Y, cb, cr := color.RGBToYCbCr(c.R, c.G, c.B)
	yOff := f.ycbcr.YOffset(int(x), int(y))
	cOff := f.ycbcr.COffset(int(x), int(y))
	if yOff < len(f.ycbcr.Y) && cOff < len(f.ycbcr.Cb) && cOff < len(f.ycbcr.Cr) {
		f.ycbcr.Y[yOff] = Y
		f.ycbcr.Cb[cOff] = cb
		f.ycbcr.Cr[cOff] = cr
	}
}

func (f *Frame) YCbCr() *image.YCbCr {
	return f.ycbcr
}

func (f *Frame) Data() []byte {
	return nil
}

func (f *Frame) Read() (img image.Image, release func(), err error) {
	return f.ycbcr, func() {}, nil
}
