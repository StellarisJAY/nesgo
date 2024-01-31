package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gen2brain/x264-go"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"github.com/stellarisJAY/nesgo/web/network"
	"log"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type MsgWithConnectionInfo struct {
	Message
	RTCRoomConnection
}

type Message struct {
	Type byte   `json:"type"`
	Data []byte `json:"data"`
}

const (
	MessageSDPOffer byte = iota
	MessageSDPAnswer
	MessageICECandidate
	MessageGameButtonPressed
	MessageGameButtonReleased
)

type RTCRoomConnection struct {
	MemberId     int64
	wsConn       *WebsocketConn
	rtcConn      *webrtc.PeerConnection
	track        *webrtc.TrackLocalStaticSample
	videoEncoder *x264.Encoder // 每个连接独占一个视频编码器 和 buffer
	videoBuffer  *bytes.Buffer

	connected *atomic.Bool
}

type WebsocketConn struct {
	Member *room.Member
	Conn   *websocket.Conn
}

type RTCRoomSession struct {
	m           *sync.Mutex
	members     map[int64]*room.Member
	connections map[int64]*RTCRoomConnection
	e           *emulator.Emulator
	signalChan  chan Signal
	cancel      context.CancelFunc

	videoEncoder *x264.Encoder
	videoBuffer  *bytes.Buffer
	encoderOpts  *x264.Options

	wsMessageChan chan MsgWithConnectionInfo
}

type Signal struct {
	Type byte
	Data any
}

const (
	SignalNewConnection byte = iota
	SignalWebsocketClose
)

var rtcFactory *network.WebRTCFactory

func init() {
	factory, err := network.NewWebRTCFactory(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	})
	if err != nil {
		panic(err)
	}
	rtcFactory = factory
}

func NewRTCRoomSession(game string) (*RTCRoomSession, error) {
	buffer := bytes.NewBuffer([]byte{})
	opts := &x264.Options{
		Width:     ppu.WIDTH,
		Height:    ppu.HEIGHT,
		Tune:      "zerolatency",
		Preset:    "ultrafast",
		FrameRate: 60,
		Profile:   "baseline",
		LogLevel:  x264.LogError,
	}
	encoder, err := x264.NewEncoder(buffer, opts)
	if err != nil {
		return nil, err
	}
	rs := &RTCRoomSession{
		m:             &sync.Mutex{},
		members:       make(map[int64]*room.Member),
		connections:   make(map[int64]*RTCRoomConnection),
		signalChan:    make(chan Signal),
		videoBuffer:   buffer,
		videoEncoder:  encoder,
		wsMessageChan: make(chan MsgWithConnectionInfo),
		encoderOpts:   opts,
	}
	game = filepath.Join(config.GetEmulatorConfig().GameDirectory, game)
	e, err := emulator.NewEmulator(game, config.GetEmulatorConfig(), rs.renderCallback)
	if err != nil {
		return nil, err
	}
	rs.e = e
	// 创建模拟器线程，暂停模拟器等待第一个连接唤醒
	go rs.e.LoadAndRun(context.Background(), false)
	rs.e.Pause()
	return rs, nil
}

func (r *RTCRoomSession) ControlLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case signal := <-r.signalChan:
			switch signal.Type {
			case SignalNewConnection:
				if conn, ok := signal.Data.(*WebsocketConn); ok {
					r.onNewConnection(ctx, conn)
				}
			case SignalWebsocketClose:
				if conn, ok := signal.Data.(*WebsocketConn); ok {
					r.onWebsocketConnClose(conn)
				}
			}
		case msg := <-r.wsMessageChan:
			if msg.wsConn.Member.MemberType == room.MemberTypeWatcher {
				continue
			}
			if msg.Type == MessageGameButtonReleased || msg.Type == MessageGameButtonPressed {
				pressed := msg.Type == MessageGameButtonPressed
				switch string(msg.Data) {
				case "KeyA":
					r.e.SetJoyPadButtonPressed(bus.Left, pressed)
				case "KeyD":
					r.e.SetJoyPadButtonPressed(bus.Right, pressed)
				case "KeyW":
					r.e.SetJoyPadButtonPressed(bus.Up, pressed)
				case "KeyS":
					r.e.SetJoyPadButtonPressed(bus.Down, pressed)
				case "Space":
					r.e.SetJoyPadButtonPressed(bus.ButtonA, pressed)
				case "KeyJ":
					r.e.SetJoyPadButtonPressed(bus.ButtonB, pressed)
				case "Enter":
					r.e.SetJoyPadButtonPressed(bus.Start, pressed)
				default:
				}
			}
		}
	}
}

func (r *RTCRoomSession) renderCallback(p *ppu.PPU) {
	p.Render()
	frame := &x264.YCbCr{YCbCr: p.Frame().YCbCr()}
	for _, conn := range r.connections {
		// peer conn没有建立连接不编码视频，否则会导致I帧没有发送到客户端
		if !conn.connected.Load() {
			continue
		}
		if err := conn.videoEncoder.Encode(frame); err != nil {
			log.Println("encoder error:", err)
			continue
		}
		if err := conn.videoEncoder.Flush(); err != nil {
			log.Println("flush encoder error:", err)
			continue
		}
		data := conn.videoBuffer.Bytes()
		if err := conn.track.WriteSample(media.Sample{
			Data:     data,
			Duration: 2 * time.Millisecond, // todo 根据帧率设置Duration
		}); err != nil {
			log.Println("write video sample error:", err)
			return
		}
		conn.videoBuffer.Reset()
	}
}

func (r *RTCRoomSession) onNewConnection(ctx context.Context, wsConn *WebsocketConn) {
	if _, ok := r.connections[wsConn.Member.UserId]; ok {
		log.Println("member already connected")
		_ = wsConn.Conn.Close()
		return
	}
	peer, err := rtcFactory.CreatePeerConnection()
	if err != nil {
		log.Println("create webrtc peer conn error:", err)
		_ = wsConn.Conn.Close()
		return
	}
	roomConnection := &RTCRoomConnection{
		MemberId:    wsConn.Member.UserId,
		wsConn:      wsConn,
		rtcConn:     peer,
		videoBuffer: bytes.NewBuffer([]byte{}),
		connected:   &atomic.Bool{},
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			_ = wsConn.Conn.Close()
			_ = peer.Close()
		}
	}()

	track, _ := rtcFactory.VideoTrack("h264")
	// add video track to peer connection
	if _, err := peer.AddTrack(track); err != nil {
		panic(fmt.Errorf("unable to add video track to peer, error: %w", err))
	}
	roomConnection.track = track
	// todo add audio track
	// create sdp offer and set local description
	sdp, err := peer.CreateOffer(nil)
	if err != nil {
		panic(fmt.Errorf("create sdp offer error: %w", err))
	}
	if err := peer.SetLocalDescription(sdp); err != nil {
		panic(fmt.Errorf("unable to set local description, error: %w", err))
	}
	// send sdp to client
	data, _ := json.Marshal(sdp)
	if err := roomConnection.sendMessage(Message{
		Type: MessageSDPOffer,
		Data: data,
	}); err != nil {
		panic(fmt.Errorf("unable to send sdp offer, error: %w", err))
	}

	peer.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Println("rtc conn state:", state)
		switch state {
		case webrtc.PeerConnectionStateConnected:
			// 检查当前连接数量，第一个建立的连接启动模拟器
			r.m.Lock()
			if len(r.connections) == 1 {
				r.e.Resume()
			}
			r.m.Unlock()
			roomConnection.connected.Store(true)
		case webrtc.PeerConnectionStateDisconnected:
			roomConnection.connected.Store(false)
			roomConnection.Close()
		default:
		}
	})

	peer.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		log.Println("ice conn state:", state)
	})

	encoder, err := x264.NewEncoder(roomConnection.videoBuffer, r.encoderOpts)
	if err != nil {
		panic(fmt.Errorf("unable to create encoder, error: %w", err))
	}
	roomConnection.videoEncoder = encoder
	r.m.Lock()
	r.members[wsConn.Member.UserId] = wsConn.Member
	r.connections[wsConn.Member.UserId] = roomConnection
	r.m.Unlock()
	go func() {
		roomConnection.Handle(context.WithoutCancel(ctx), r.wsMessageChan)
		r.signalChan <- Signal{
			Type: SignalWebsocketClose,
			Data: roomConnection.wsConn,
		}
	}()
}

func (r *RTCRoomSession) onWebsocketConnClose(wsConn *WebsocketConn) {
	r.m.Lock()
	delete(r.connections, wsConn.Member.UserId)
	delete(r.members, wsConn.Member.UserId)
	// 已经没有活跃连接，暂停模拟器线程
	if len(r.connections) == 0 {
		r.e.Pause()
	}
	r.m.Unlock()
}

func (rc *RTCRoomConnection) Handle(ctx context.Context, msgChan chan MsgWithConnectionInfo) {
	wsConn := rc.wsConn.Conn
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msgType, payload, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("ws read error:", err)
			return
		}
		switch msgType {
		case websocket.TextMessage:
			msg := Message{}
			if err := json.Unmarshal(payload, &msg); err != nil {
				log.Println("invalid message error:", err)
				continue
			}
			rc.HandleMessage(msg, msgChan)
		case websocket.BinaryMessage:
		}
	}
}

func (rc *RTCRoomConnection) sendMessage(msg Message) error {
	payload, _ := json.Marshal(msg)
	return rc.wsConn.Conn.WriteMessage(websocket.TextMessage, payload)
}

func (rc *RTCRoomConnection) Close() {
	_ = rc.wsConn.Conn.Close()
	_ = rc.rtcConn.Close()
}

func (rc *RTCRoomConnection) HandleMessage(msg Message, msgChan chan MsgWithConnectionInfo) {
	switch msg.Type {
	case MessageSDPAnswer:
		sdp := webrtc.SessionDescription{}
		_ = json.Unmarshal(msg.Data, &sdp)
		if err := rc.rtcConn.SetRemoteDescription(sdp); err != nil {
			log.Println("unable to set remote description, error:", err)
			rc.Close()
		}
	case MessageICECandidate:
		candidate := webrtc.ICECandidateInit{}
		_ = json.Unmarshal(msg.Data, &candidate)
		log.Println("candidate:", candidate)
		if err := rc.rtcConn.AddICECandidate(candidate); err != nil {
			log.Println("unable to add ICE candidate, error:", err)
			rc.Close()
		}
	case MessageGameButtonPressed, MessageGameButtonReleased:
		msgChan <- MsgWithConnectionInfo{
			Message:           msg,
			RTCRoomConnection: *rc,
		}
	default:
	}
}
