package biz

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
	"sync"
)

type Connection struct {
	pc          *webrtc.PeerConnection
	videoTrack  *webrtc.TrackLocalStaticSample
	audioTrack  *webrtc.TrackLocalStaticSample
	dataChannel *webrtc.DataChannel
	userId      int64

	localCandidates []*webrtc.ICECandidate
	mutex           *sync.Mutex
}

func (g *GameInstance) NewConnection(userId int64, stunServer string) (*Connection, string, error) {
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{stunServer},
			},
		},
	})
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
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "nesgo_video")
	if err != nil {
		panic(fmt.Errorf("create video track error: %v", err))
	}
	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "nesgo_audio")
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
		mutex:       &sync.Mutex{},
	}

	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			g.mutex.Lock()
			conn.localCandidates = append(conn.localCandidates, candidate)
			g.mutex.Unlock()
		}
	})

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		conn.OnPeerConnectionState(state, g)
	})
	pc.OnICEConnectionStateChange(conn.OnICEStateChange)
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		conn.OnDataChannelMessage(g, msg)
	})
	rch := make(chan ConsumerResult)
	g.messageChan <- &Message{
		Type:       MsgNewConn,
		Data:       conn,
		resultChan: rch,
	}
	<-rch
	return conn, offer.SDP, nil
}

func (c *Connection) OnPeerConnectionState(state webrtc.PeerConnectionState, instance *GameInstance) {
	switch state {
	case webrtc.PeerConnectionStateConnected:
		instance.onConnected(c)
	case webrtc.PeerConnectionStateFailed:
		instance.closeConnection(c)
	case webrtc.PeerConnectionStateDisconnected:
		instance.closeConnection(c)
	case webrtc.PeerConnectionStateClosed:
		instance.closeConnection(c)
	default:
	}
}

func (c *Connection) OnICEStateChange(_ webrtc.ICEConnectionState) {

}

func (c *Connection) OnDataChannelMessage(instance *GameInstance, msg webrtc.DataChannelMessage) {
	m := &Message{}
	err := json.Unmarshal(msg.Data, m)
	if err != nil {
		return
	}
	if m.Type == MsgPing {
		_ = c.dataChannel.Send(msg.Data)
	} else {
		instance.onDataChannelMessage(c.userId, m)
	}
}

func (c *Connection) Close() {
	_ = c.dataChannel.Close()
	_ = c.pc.Close()
}
