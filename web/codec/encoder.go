package codec

import (
	"errors"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/stellarisJAY/nesgo/web/codec/h264"
)

type IEncoder interface {
	Encode(frame *ppu.Frame) error
	Flush() error
	FlushBuffer() []byte
}

func NewEncoder(codec string) (IEncoder, error) {
	switch codec {
	case "h264":
		return h264.NewEncoder()
	}
	panic(errors.New("codec not available"))
}
