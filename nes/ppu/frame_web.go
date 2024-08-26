//go:build !sdl

package ppu

import (
	"cmp"
	"image"
	"image/color"
	"slices"
)

type Frame struct {
	ycbcr              *image.YCbCr
	pixelPreprocessors []PixelPreprocessor
}

const (
	WIDTH  = 256
	HEIGHT = 240
)

type PixelPreprocessorFunc func(c Color) Color
type PixelPreprocessor struct {
	f        PixelPreprocessorFunc
	priority int
}

const (
	ReverseColorPriority = 1
	GrayscalePriority    = 2
)

func NewFrame() *Frame {
	return &Frame{
		image.NewYCbCr(image.Rect(0, 0, WIDTH, HEIGHT), image.YCbCrSubsampleRatio420),
		make([]PixelPreprocessor, 0),
	}
}

func NewCustomSizeFrame(ycbcr *image.YCbCr) *Frame {
	return &Frame{
		ycbcr,
		make([]PixelPreprocessor, 0),
	}
}

func (f *Frame) setPixel(x, y uint32, c Color) {
	for _, p := range f.pixelPreprocessors {
		c = p.f(c)
	}
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

// UseReverseColorPreprocessor 反色滤镜
func (f *Frame) UseReverseColorPreprocessor() {
	f.pixelPreprocessors = append(f.pixelPreprocessors, PixelPreprocessor{
		f: func(c Color) Color {
			return Color{^c.R, ^c.G, ^c.B}
		},
		priority: ReverseColorPriority,
	})
	slices.SortFunc(f.pixelPreprocessors, func(a, b PixelPreprocessor) int {
		return cmp.Compare(a.priority, b.priority)
	})
}

// UseGrayscalePreprocessor 灰度滤镜
func (f *Frame) UseGrayscalePreprocessor() {
	f.pixelPreprocessors = append(f.pixelPreprocessors, PixelPreprocessor{
		f: func(c Color) Color {
			gray := (int(c.R)*299 + int(c.G)*587 + int(c.B)*114) / 1000
			c.R, c.G, c.B = uint8(gray), uint8(gray), uint8(gray)
			return c
		},
		priority: GrayscalePriority,
	})
	slices.SortFunc(f.pixelPreprocessors, func(a, b PixelPreprocessor) int {
		return cmp.Compare(a.priority, b.priority)
	})
}

func (f *Frame) RemoveGrayscalePreprocessor() {
	f.pixelPreprocessors = slices.DeleteFunc(f.pixelPreprocessors, func(p PixelPreprocessor) bool {
		return p.priority == GrayscalePriority
	})
}

func (f *Frame) RemoveReverseColorPreprocessor() {
	f.pixelPreprocessors = slices.DeleteFunc(f.pixelPreprocessors, func(p PixelPreprocessor) bool {
		return p.priority == ReverseColorPriority
	})
}

func (f *Frame) ResetPixelPreprocessor() {
	f.pixelPreprocessors = make([]PixelPreprocessor, 0)
}
