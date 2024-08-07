package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
)

type Connection struct {
	pc          *webrtc.PeerConnection
	videoTrack  *webrtc.TrackLocalStaticSample
	audioTrack  *webrtc.TrackLocalStaticSample
	dataChannel *webrtc.DataChannel
	userId      int64
}

func (g *GameInstance) NewConnection(userId int64) (*Connection, string, error) {
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return nil, "", fmt.Errorf("new peer connection error: %v", err)
	}

	defer func() {
		if err := recover(); err != nil {
			log.NewHelper(log.DefaultLogger).Errorf("new peer connection error: %v", err)
			_ = pc.Close()
		}
	}()

	// create video and audio tracks
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "nesgo")
	if err != nil {
		panic(fmt.Errorf("create video track error: %v", err))
	}
	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "nesgo")
	if err != nil {
		panic(fmt.Errorf("create audio track error: %v", err))
	}
	if _, err := pc.AddTrack(videoTrack); err != nil {
		panic(fmt.Errorf("add video track error: %v", err))
	}
	if _, err := pc.AddTrack(audioTrack); err != nil {
		panic(fmt.Errorf("add audio track error: %v", err))
	}
	// create dataChannel to transfer control messages
	dataChannel, err := pc.CreateDataChannel("control-channel", nil)
	if err != nil {
		panic(fmt.Errorf("create data channel error: %v", err))
	}
	// create sdp offer and set local description
	offer, err := pc.CreateOffer(nil)
	if err != nil {
		panic(fmt.Errorf("create sdp offer error: %v", err))
	}
	if err := pc.SetLocalDescription(offer); err != nil {
		panic(fmt.Errorf("set local description error: %v", err))
	}

	conn := &Connection{
		pc:          pc,
		videoTrack:  videoTrack,
		audioTrack:  audioTrack,
		dataChannel: dataChannel,
		userId:      userId,
	}

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		conn.OnPeerConnectionState(state, g)
	})
	pc.OnICEConnectionStateChange(conn.OnICEStateChange)
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		g.onDataChannelMessage(userId, msg.Data)
	})
	g.mutex.Lock()
	g.connections[userId] = conn
	g.mutex.Unlock()
	return conn, offer.SDP, nil
}

func (c *Connection) OnPeerConnectionState(state webrtc.PeerConnectionState, instance *GameInstance) {
	// TODO Handle conn state change
	switch state {
	case webrtc.PeerConnectionStateConnected:

	case webrtc.PeerConnectionStateClosed:
		instance.closeConnection(c)
	default:
	}
}

func (c *Connection) OnICEStateChange(state webrtc.ICEConnectionState) {
	// TODO Log ice state change
}
