package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/stellarisJAY/nesgo/web/codec"
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
	videoTrack   *webrtc.TrackLocalStaticSample
	audioTrack   *webrtc.TrackLocalStaticSample
	videoEncoder codec.IVideoEncoder // 每个连接独占一个视频编码器 和 buffer
	connected    *atomic.Bool
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

	wsMessageChan chan MsgWithConnectionInfo

	emulatorCancel  context.CancelFunc
	game            string
	audioSampleRate int
	audioEncoder    codec.IAudioEncoder
	audioSampleChan chan float32
}

type Signal struct {
	Type byte
	Data any
}

type restartEmulatorRequest struct {
	game     string
	respChan chan error
}

type saveGameRequest struct {
	errChan  chan error
	respChan chan []byte
}

type loadSavedGameRequest struct {
	data     []byte
	respChan chan error
}

const (
	SignalNewConnection byte = iota
	SignalWebsocketClose
	SignalRestartEmulator
	SignalSaveGame
	SignalLoadSavedGame
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
	const sampleRate = 48000
	rs := &RTCRoomSession{
		m:               &sync.Mutex{},
		members:         make(map[int64]*room.Member),
		connections:     make(map[int64]*RTCRoomConnection),
		signalChan:      make(chan Signal),
		wsMessageChan:   make(chan MsgWithConnectionInfo),
		game:            game,
		audioSampleRate: sampleRate,
	}
	game = filepath.Join(config.GetEmulatorConfig().GameDirectory, game)
	sampleChan := make(chan float32, sampleRate)
	e, err := emulator.NewEmulator(game, config.GetEmulatorConfig(), rs.renderCallback, sampleChan, sampleRate)
	if err != nil {
		return nil, err
	}
	rs.e = e
	rs.audioSampleChan = sampleChan
	audioEncoder, err := codec.NewAudioEncoder(sampleRate)
	if err != nil {
		return nil, err
	}
	rs.audioEncoder = audioEncoder
	ctx, cancelFunc := context.WithCancel(context.Background())
	// 创建模拟器线程，暂停模拟器等待第一个连接唤醒
	go rs.e.LoadAndRun(ctx, false)
	rs.emulatorCancel = cancelFunc
	rs.e.Pause()
	return rs, nil
}

func (r *RTCRoomSession) ControlLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case signal := <-r.signalChan:
			r.handleSignal(ctx, signal)
		case msg := <-r.wsMessageChan:
			if msg.wsConn.Member.MemberType == room.MemberTypeWatcher {
				continue
			}
			if msg.Type == MessageGameButtonReleased || msg.Type == MessageGameButtonPressed {
				r.handleGameButtonMsg(msg.Message)
			}
		}
	}
}

func (r *RTCRoomSession) handleSignal(ctx context.Context, signal Signal) {
	switch signal.Type {
	case SignalNewConnection:
		if conn, ok := signal.Data.(*WebsocketConn); ok {
			r.onNewConnection(ctx, conn)
		}
	case SignalWebsocketClose:
		if conn, ok := signal.Data.(*WebsocketConn); ok {
			r.onWebsocketConnClose(conn)
		}
	case SignalRestartEmulator:
		req := signal.Data.(*restartEmulatorRequest)
		if err := r.restart(req.game); err != nil {
			req.respChan <- err
		} else {
			close(req.respChan)
		}
	case SignalSaveGame:
		req := signal.Data.(*saveGameRequest)
		if data, err := r.save(); err != nil {
			req.errChan <- err
		} else {
			req.respChan <- data
		}
	case SignalLoadSavedGame:
		req := signal.Data.(*loadSavedGameRequest)
		if err := r.loadSavedGame(req.data); err != nil {
			req.respChan <- err
		} else {
			close(req.respChan)
		}
	}
}

func (r *RTCRoomSession) handleGameButtonMsg(msg Message) {
	pressed := msg.Type == MessageGameButtonPressed
	switch string(msg.Data) {
	case "Left":
		r.e.SetJoyPadButtonPressed(bus.Left, pressed)
	case "Right":
		r.e.SetJoyPadButtonPressed(bus.Right, pressed)
	case "Up":
		r.e.SetJoyPadButtonPressed(bus.Up, pressed)
	case "Down":
		r.e.SetJoyPadButtonPressed(bus.Down, pressed)
	case "A":
		r.e.SetJoyPadButtonPressed(bus.ButtonA, pressed)
	case "B":
		r.e.SetJoyPadButtonPressed(bus.ButtonB, pressed)
	case "Start":
		r.e.SetJoyPadButtonPressed(bus.Start, pressed)
	case "Select":
		r.e.SetJoyPadButtonPressed(bus.Select, pressed)
	default:
	}
}

func (r *RTCRoomSession) Restart(game string) error {
	respChan := make(chan error)
	r.signalChan <- Signal{
		Type: SignalRestartEmulator,
		Data: &restartEmulatorRequest{
			game:     game,
			respChan: respChan,
		},
	}
	return <-respChan
}

func (r *RTCRoomSession) restart(game string) error {
	r.game = game
	game = filepath.Join(config.GetEmulatorConfig().GameDirectory, game)

LOOP:
	for {
		select {
		case <-r.audioSampleChan:
		default:
			break LOOP
		}
	}

	e, err := emulator.NewEmulator(game, config.GetEmulatorConfig(), r.renderCallback, r.audioSampleChan, r.audioSampleRate)
	if err != nil {
		return err
	}
	r.emulatorCancel()
	ctx, cancelFunc := context.WithCancel(context.Background())
	r.emulatorCancel = cancelFunc
	go e.LoadAndRun(ctx, false)
	if len(r.connections) == 0 {
		e.Pause()
	}
	r.e = e
	return nil
}

func (r *RTCRoomSession) Save() ([]byte, error) {
	respChan := make(chan []byte)
	errChan := make(chan error)
	r.signalChan <- Signal{
		Type: SignalSaveGame,
		Data: &saveGameRequest{errChan, respChan},
	}
	defer close(errChan)
	defer close(respChan)
	select {
	case data := <-respChan:
		return data, nil
	case err := <-errChan:
		return nil, err
	}
}

func (r *RTCRoomSession) save() ([]byte, error) {
	r.e.Pause()
	defer r.e.Resume()
	return r.e.GetSaveData()
}

func (r *RTCRoomSession) LoadSavedGame(data []byte) error {
	respChan := make(chan error)
	r.signalChan <- Signal{
		Type: SignalLoadSavedGame,
		Data: &loadSavedGameRequest{data, respChan},
	}
	return <-respChan
}

func (r *RTCRoomSession) loadSavedGame(data []byte) error {
	r.e.Pause()
	defer r.e.Resume()
	return r.e.Load(data)
}

func (r *RTCRoomSession) renderCallback(p *ppu.PPU) {
	p.Render()
	for _, conn := range r.connections {
		// peer conn没有建立连接不编码视频，否则会导致I帧没有发送到客户端
		if !conn.connected.Load() {
			continue
		}
		if err := conn.videoEncoder.Encode(p.Frame()); err != nil {
			log.Println("encoder error:", err)
			continue
		}
		if err := conn.videoEncoder.Flush(); err != nil {
			log.Println("flush encoder error:", err)
			continue
		}
		data := conn.videoEncoder.FlushBuffer()
		if err := conn.videoTrack.WriteSample(media.Sample{
			Data:     data,
			Duration: 2 * time.Millisecond, // todo 根据帧率设置Duration
		}); err != nil {
			log.Println("write video sample error:", err)
			return
		}
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
		MemberId:  wsConn.Member.UserId,
		wsConn:    wsConn,
		rtcConn:   peer,
		connected: &atomic.Bool{},
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			_ = wsConn.Conn.Close()
			_ = peer.Close()
		}
	}()

	videoTrack, _ := rtcFactory.VideoTrack("h264")
	// add video track to peer connection
	if _, err := peer.AddTrack(videoTrack); err != nil {
		panic(fmt.Errorf("unable to add video track to peer, error: %w", err))
	}
	roomConnection.videoTrack = videoTrack

	audioTrack, _ := rtcFactory.AudioTrack("opus")
	if _, err := peer.AddTrack(audioTrack); err != nil {
		panic(fmt.Errorf("unable to add audio track to peer, error: %w", err))
	}
	roomConnection.audioTrack = audioTrack
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

	encoder, err := codec.NewVideoEncoder("h264")
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

func (r *RTCRoomSession) audioSampleListener(ctx context.Context) {
	buffer := make([]float32, 0, r.audioSampleRate*5/1000)
	for {
		select {
		case <-ctx.Done():
			return
		case s := <-r.audioSampleChan:
			buffer = append(buffer, s)
			if len(buffer) == cap(buffer) {
				r.sendAudioSamples(buffer)
				buffer = buffer[:0]
			}
		}
	}
}

func (r *RTCRoomSession) sendAudioSamples(samples []float32) {
	frame, err := r.audioEncoder.Encode(samples)
	if err != nil {
		log.Println(err)
		return
	}
	for _, conn := range r.connections {
		if !conn.connected.Load() {
			continue
		}
		if err := conn.audioTrack.WriteSample(media.Sample{
			Data:      frame,
			Timestamp: time.Now(),
			Duration:  5 * time.Millisecond,
		}); err != nil {
			log.Println("send audio frame to:", conn.wsConn.Conn.RemoteAddr(), "error:", err)
		}
	}
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
