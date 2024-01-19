package recorder

import (
	"github.com/gen2brain/x264-go"
	"github.com/stellarisJAY/nesgo/ppu"
	"os"
)

type Recorder struct {
	file *os.File
	enc  *x264.Encoder
}

func NewRecorder(saveFile string) (*Recorder, error) {
	file, err := os.OpenFile(saveFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	opts := &x264.Options{
		Width:        ppu.WIDTH,
		Height:       ppu.HEIGHT,
		Tune:         "zerolatency",
		Preset:       "ultrafast",
		RateControl:  "crf",
		RateConstant: 23,
		LogLevel:     x264.LogDebug,
	}
	enc, err := x264.NewEncoder(file, opts)
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	r := &Recorder{
		file: file,
		enc:  enc,
	}
	return r, nil
}

func (r *Recorder) Encode(frame *ppu.Frame) error {
	return r.enc.Encode(frame)
}

func (r *Recorder) Close() {
	_ = r.enc.Flush()
	_ = r.enc.Close()
	_ = r.file.Close()
}
