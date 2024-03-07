package codec

import (
	"github.com/pion/mediadevices/pkg/codec/x264"
	"github.com/pion/mediadevices/pkg/prop"
)

func NewX264Encoder(media prop.Media) (*VideoEncoder, error) {
	params, err := x264.NewParams()
	if err != nil {
		return nil, err
	}
	params.KeyFrameInterval = 60
	params.Preset = x264.PresetUltrafast
	params.BitRate = 3000_000
	r := &FrameReader{}
	enc, err := params.BuildVideoEncoder(r, media)
	if err != nil {
		return nil, err
	}
	return &VideoEncoder{
		enc:         enc,
		params:      params,
		frameReader: r,
	}, nil
}
