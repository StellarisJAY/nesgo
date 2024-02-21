package h264

import (
	"bytes"
	"github.com/gen2brain/x264-go"
	"github.com/stellarisJAY/nesgo/ppu"
)

type Encoder struct {
	enc    *x264.Encoder
	opts   *x264.Options
	buffer *bytes.Buffer
}

func (e *Encoder) FlushBuffer() []byte {
	data := e.buffer.Bytes()
	e.buffer.Reset()
	return data
}

func NewEncoder() (*Encoder, error) {
	buffer := bytes.NewBuffer([]byte{})
	opts := &x264.Options{
		Width:     ppu.WIDTH,
		Height:    ppu.HEIGHT,
		Tune:      "zerolatency",
		Preset:    "ultrafast",
		FrameRate: 60,
		Profile:   "baseline",
		LogLevel:  x264.LogError,
	}
	enc, err := x264.NewEncoder(buffer, opts)
	if err != nil {
		return nil, err
	}
	e := &Encoder{
		enc:    enc,
		opts:   opts,
		buffer: buffer,
	}
	return e, nil
}

func (e *Encoder) Encode(frame *ppu.Frame) error {
	return e.enc.Encode(frame.YCbCr())
}

func (e *Encoder) Flush() error {
	return e.enc.Flush()
}
