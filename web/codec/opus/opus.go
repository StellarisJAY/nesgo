package opus

import (
	"gopkg.in/hraban/opus.v2"
	"log"
	"math"
)

type Encoder struct {
	enc *opus.Encoder

	SampleRate int
	Channels   int
}

func NewEncoder(sampleRate int) (*Encoder, error) {
	enc, err := opus.NewEncoder(sampleRate, 1, opus.AppVoIP)
	if err != nil {
		log.Println("new encoder error: ", err)
		return nil, err
	}
	return &Encoder{
		enc:        enc,
		SampleRate: sampleRate,
		Channels:   1,
	}, nil
}

func (e *Encoder) Encode(pcm []float32) ([]byte, error) {
	pcm16 := make([]int16, len(pcm))
	for i, s32 := range pcm {
		pcm16[i] = pcmFloat32ToInt16(s32)
	}
	buffer := make([]byte, 1024)
	n, err := e.enc.Encode(pcm16, buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func pcmFloat32ToInt16(s float32) int16 {
	if s < -0.9999 {
		return math.MinInt16
	} else if s > 0.9999 {
		return math.MaxInt16
	} else {
		return int16(s * 32767)
	}
}
