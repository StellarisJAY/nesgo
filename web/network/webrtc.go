package network

import (
	"errors"
	"github.com/pion/webrtc/v3"
)

type WebRTCFactory struct {
	api      *webrtc.API
	peerConf webrtc.Configuration
}

func NewWebRTCFactory(peerConf webrtc.Configuration) (*WebRTCFactory, error) {
	m := webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		return nil, err
	}
	w := &WebRTCFactory{
		api:      webrtc.NewAPI(webrtc.WithMediaEngine(&m)),
		peerConf: peerConf,
	}
	return w, nil
}

func (w *WebRTCFactory) CreatePeerConnection() (*webrtc.PeerConnection, error) {
	conn, err := webrtc.NewPeerConnection(w.peerConf)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (w *WebRTCFactory) VideoTrack(codec string) (*webrtc.TrackLocalStaticSample, error) {
	switch codec {
	case "h264":
		codec = webrtc.MimeTypeH264
	default:
		return nil, errors.New("unsupported video codec")
	}
	return webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{
		MimeType: codec,
	}, "video", "nesgo-webrtc")
}

func (w *WebRTCFactory) AudioTrack(codec string) (*webrtc.TrackLocalStaticSample, error) {
	switch codec {
	case "opus":
		codec = webrtc.MimeTypeOpus
	default:
		return nil, errors.New("unsupported video codec")
	}
	return webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{
		MimeType: codec,
	}, "audio", "nesgo-webrtc")
}
