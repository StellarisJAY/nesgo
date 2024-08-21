package codec

import (
	"fmt"
	"github.com/pion/mediadevices/pkg/codec"
	"github.com/pion/mediadevices/pkg/codec/vpx"
	"github.com/pion/mediadevices/pkg/prop"
)

func NewVpxEncoder(media prop.Media, version int) (*VideoEncoder, error) {
	var enc codec.ReadCloser = nil
	frameReader := &FrameReader{}
	if version == 8 {
		params, err := vpx.NewVP8Params()
		if err != nil {
			return nil, err
		}
		params.BitRate = 8_000_000
		params.ErrorResilient = vpx.ErrorResilientDefault
		params.RateControlEndUsage = vpx.RateControlVBR
		params.RateControlMaxQuantizer = 4
		params.RateControlMinQuantizer = 4
		enc, err = params.BuildVideoEncoder(frameReader, media)
		if err != nil {
			return nil, err
		}
	} else if version == 9 {
		params, err := vpx.NewVP9Params()
		if err != nil {
			return nil, err
		}
		params.BitRate = 8_000_000
		params.ErrorResilient = vpx.ErrorResilientDefault
		params.RateControlEndUsage = vpx.RateControlVBR
		params.RateControlMaxQuantizer = 4
		params.RateControlMinQuantizer = 4
		params.LagInFrames = 15
		enc, err = params.BuildVideoEncoder(frameReader, media)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unsupported version %d", version)
	}

	return &VideoEncoder{
		enc:         enc,
		frameReader: frameReader,
	}, nil
}
