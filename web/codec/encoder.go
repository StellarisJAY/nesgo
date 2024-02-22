package codec

import (
	"errors"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/stellarisJAY/nesgo/web/codec/h264"
	"github.com/stellarisJAY/nesgo/web/codec/opus"
)

type IVideoEncoder interface {
	Encode(frame *ppu.Frame) error
	Flush() error
	FlushBuffer() []byte
}

func NewVideoEncoder(codec string) (IVideoEncoder, error) {
	switch codec {
	case "h264":
		return h264.NewEncoder()
	}
	panic(errors.New("codec not available"))
}

type IAudioEncoder interface {
	// Encode PCM to opus packet, Emulator outputs float32 PCM
	Encode(pcm []float32) ([]byte, error)
}

func NewAudioEncoder(sampleRate int) (IAudioEncoder, error) {
	return opus.NewEncoder(sampleRate)
}
