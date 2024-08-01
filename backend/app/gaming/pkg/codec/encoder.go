package codec

import (
	"errors"
	"github.com/pion/mediadevices/pkg/codec"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/pkg/codec/opus"
	"github.com/stellarisJAY/nesgo/emulator/ppu"
	"image"
	"log"
)

type IVideoEncoder interface {
	Encode(frame *ppu.Frame) ([]byte, func(), error)
	Close()
}

type FrameReader struct {
	frame *ppu.Frame
}

func (f *FrameReader) Read() (img image.Image, release func(), err error) {
	return f.frame.Read()
}

func (f *FrameReader) SetFrame(frame *ppu.Frame) {
	f.frame = frame
}

type VideoEncoder struct {
	enc         codec.ReadCloser
	params      any
	frameReader *FrameReader
}

func NewVideoEncoder(codec string) (IVideoEncoder, error) {
	media := prop.Media{
		DeviceID: "nesgo-video",
		Video: prop.Video{
			Width:       ppu.WIDTH,
			Height:      ppu.HEIGHT,
			FrameRate:   60,
			FrameFormat: frame.FormatI444,
		},
	}
	switch codec {
	case "h264":
		return NewX264Encoder(media)
	default:
		panic(errors.New("codec not available"))
	}
}

func (v *VideoEncoder) Encode(frame *ppu.Frame) ([]byte, func(), error) {
	v.frameReader.SetFrame(frame)
	return v.enc.Read()
}

func (v *VideoEncoder) Close() {
	if err := v.enc.Close(); err != nil {
		log.Println("close video encoder error:", err)
	}
}

type IAudioEncoder interface {
	// Encode PCM to opus packet, Emulator outputs float32 PCM
	Encode(pcm []float32) ([]byte, error)
}

func NewAudioEncoder(sampleRate int) (IAudioEncoder, error) {
	return opus.NewEncoder(sampleRate)
}
