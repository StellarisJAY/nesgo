package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pion/webrtc/v3"
	"log"
	"path/filepath"
	"sync/atomic"
	"time"

	_ "github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisJAY/nesgo/emulator/bus"
	"github.com/stellarisJAY/nesgo/emulator/config"
	"github.com/stellarisJAY/nesgo/emulator/ppu"
	"github.com/stellarisJAY/nesgo/web/codec"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"github.com/stellarisJAY/nesgo/web/network"
	"github.com/stellarisJAY/nesgo/web/util/future"
)

type WebRTCRoomSession struct {
	cancel             context.CancelFunc
	members            map[int64]*room.Member
	connections        map[int64]*RoomConn
	signalChan         chan Signal              // signalChan 传递模拟器信号、连接信号
	dataChannelMsgChan chan MessageWithConnInfo // dataChannelMsgChan 用于接收webrtc datachannel
	e                  *emulator.Emulator
	emulatorCancel     context.CancelFunc
	game               string
	videoEncoder       codec.IVideoEncoder
	audioSampleRate    int
	audioEncoder       codec.IAudioEncoder
	audioSampleChan    chan float32 // audioSampleChan 音频输出channel
	controller1        int64        // controller1 模拟器P1控制权玩家
	controller2        int64        // controller2 模拟器P2控制权玩家
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
	game     string
	respChan chan error
}

type transferControlRequest struct {
	memberId int64
	control1 bool
	control2 bool
	future   *future.Future[struct{}]
}

type kickMemberRequest struct {
	memberId int64
	future   *future.Future[struct{}]
}

type alterRoleRequest struct {
	memberId int64
	role     byte
	future   *future.Future[struct{}]
}

type ChatMessage struct {
	From    int64  `json:"from"`
	Content string `json:"content"`
}

const (
	SignalNewConnection byte = iota
	SignalWebsocketClose
	SignalPeerConnected
	SignalPeerDisconnected
	SignalRestartEmulator
	SignalSaveGame
	SignalLoadSavedGame
	SignalTransferControl
	SignalKickMember
	SignalAlterRole

	VideoCodec = "h264"
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

func NewRTCRoomSession(game string) (*WebRTCRoomSession, error) {
	const sampleRate = 48000
	rs := &WebRTCRoomSession{
		members:            make(map[int64]*room.Member),
		connections:        make(map[int64]*RoomConn),
		signalChan:         make(chan Signal),
		dataChannelMsgChan: make(chan MessageWithConnInfo),
		game:               game,
		audioSampleRate:    sampleRate,
	}
	game = filepath.Join(config.GetEmulatorConfig().GameDirectory, game)
	sampleChan := make(chan float32, sampleRate)
	e, err := emulator.NewEmulator(game, config.GetEmulatorConfig(), rs.renderCallback, sampleChan, sampleRate)
	if err != nil {
		return nil, err
	}
	rs.e = e
	encoder, err := codec.NewVideoEncoder(VideoCodec)
	if err != nil {
		return nil, err
	}
	rs.videoEncoder = encoder
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

// ControlLoop 所有房间相关的消息和信号，都由control 循环处理，避免线程安全问题
func (r *WebRTCRoomSession) ControlLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case signal := <-r.signalChan:
			r.handleSignal(ctx, signal)
		case msg := <-r.dataChannelMsgChan:
			r.handleDataChannelMessage(msg)
		}
	}
}

func (r *WebRTCRoomSession) handleSignal(ctx context.Context, signal Signal) {
	switch signal.Type {
	case SignalNewConnection:
		if conn, ok := signal.Data.(*WebsocketConn); ok {
			r.onNewConnection(ctx, conn)
		}
	case SignalWebsocketClose:
		if conn, ok := signal.Data.(*WebsocketConn); ok {
			r.onWebsocketConnClose(conn)
		}
	case SignalPeerConnected:
		conn, _ := signal.Data.(*RoomConn)
		// 检查当前连接数量，第一个建立的连接启动模拟器
		controllers := r.filterMembers(isControllableMember)
		if len(controllers) == 1 {
			r.e.Resume()
		}
		conn.connected.Store(true)
	case SignalPeerDisconnected:
		conn, _ := signal.Data.(*RoomConn)
		delete(r.connections, conn.MemberId)
		delete(r.members, conn.MemberId)
		// 无可控制游戏的用户，暂停模拟器
		controllers := r.filterMembers(isControllableMember)
		if len(controllers) == 0 {
			r.e.Pause()
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
		if err := r.loadSavedGame(req.game, req.data); err != nil {
			req.respChan <- err
		} else {
			close(req.respChan)
		}
	case SignalTransferControl:
		req := signal.Data.(transferControlRequest)
		if err := r.transferControl(req.memberId, req.control1, req.control2); err != nil {
			req.future.Fail(err)
		} else {
			req.future.Success(&struct{}{})
		}
	case SignalKickMember:
		req := signal.Data.(kickMemberRequest)
		if err := r.kickMember(req.memberId); err != nil {
			req.future.Fail(err)
		} else {
			req.future.Success(&struct{}{})
		}
	case SignalAlterRole:
		req := signal.Data.(alterRoleRequest)
		_ = r.alterRole(req.memberId, req.role)
		req.future.Success(&struct{}{})
	}
}

func (r *WebRTCRoomSession) handleDataChannelMessage(msg MessageWithConnInfo) {
	switch msg.Type {
	case MessageGameButtonReleased, MessageGameButtonPressed:
		if msg.wsConn.Member.Role != room.RoleObserver {
			r.handleGameButtonMsg(msg)
		}
	case MessageChat:
		r.handleChatMessage(msg)
	default:
	}
}

func (r *WebRTCRoomSession) handleGameButtonMsg(msg MessageWithConnInfo) {
	var controlId int
	if msg.MemberId == r.controller1 {
		controlId = 1
	} else if msg.MemberId == r.controller2 {
		controlId = 2
	} else {
		return
	}
	pressed := msg.Type == MessageGameButtonPressed
	switch string(msg.Data) {
	case "Left":
		r.e.SetJoyPadButtonPressed(controlId, bus.Left, pressed)
	case "Right":
		r.e.SetJoyPadButtonPressed(controlId, bus.Right, pressed)
	case "Up":
		r.e.SetJoyPadButtonPressed(controlId, bus.Up, pressed)
	case "Down":
		r.e.SetJoyPadButtonPressed(controlId, bus.Down, pressed)
	case "A":
		r.e.SetJoyPadButtonPressed(controlId, bus.ButtonA, pressed)
	case "B":
		r.e.SetJoyPadButtonPressed(controlId, bus.ButtonB, pressed)
	case "Start":
		r.e.SetJoyPadButtonPressed(controlId, bus.Start, pressed)
	case "Select":
		r.e.SetJoyPadButtonPressed(controlId, bus.Select, pressed)
	default:
	}
}

func (r *WebRTCRoomSession) handleChatMessage(msg MessageWithConnInfo) {
	m, ok := r.members[msg.MemberId]
	if !ok {
		return
	}
	chat := ChatMessage{
		From:    m.UserId,
		Content: msg.Data,
	}
	data, _ := json.Marshal(chat)
	bytes, _ := json.Marshal(Message{
		Type: MessageChat,
		Data: string(data),
	})
	for _, conn := range r.connections {
		if !conn.connected.Load() {
			continue
		}
		if err := conn.dataChannel.SendText(string(bytes)); err != nil {
			log.Println("send chat message error:", err)
		}
	}
}

func (r *WebRTCRoomSession) Restart(game string) error {
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

func (r *WebRTCRoomSession) restart(game string) error {
	r.game = game
	game = filepath.Join(config.GetEmulatorConfig().GameDirectory, game)
	// 重启模拟器之前清空音频输出chan，避免上一个模拟器的音频在新模拟器播放
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

func (r *WebRTCRoomSession) Save() ([]byte, error) {
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

func (r *WebRTCRoomSession) save() ([]byte, error) {
	r.e.Pause()
	defer r.e.Resume()
	return r.e.GetSaveData()
}

func (r *WebRTCRoomSession) LoadSavedGame(game string, data []byte) error {
	respChan := make(chan error)
	r.signalChan <- Signal{
		Type: SignalLoadSavedGame,
		Data: &loadSavedGameRequest{data, game, respChan},
	}
	return <-respChan
}

func (r *WebRTCRoomSession) loadSavedGame(game string, data []byte) error {
	if r.game != game {
		if err := r.restart(game); err != nil {
			return err
		}
	}
	r.e.Pause()
	defer r.e.Resume()
	return r.e.Load(data)
}

func (r *WebRTCRoomSession) TransferControl(memberId int64, control1, control2 bool) error {
	req := transferControlRequest{
		memberId: memberId,
		control1: control1,
		control2: control2,
		future:   future.NewFuture[struct{}](),
	}
	r.signalChan <- Signal{
		Type: SignalTransferControl,
		Data: req,
	}
	_, err := req.future.Result()
	return err
}

func (r *WebRTCRoomSession) transferControl(memberId int64, control1, control2 bool) error {
	if _, ok := r.members[memberId]; !ok {
		return errors.New("member not connected")
	}
	if r.members[memberId].Role == room.RoleObserver {
		return errors.New("observer can't not gain control")
	} else {
		if control1 {
			r.controller1 = memberId
		} else if r.controller1 == memberId {
			r.controller1 = 0
		}
		if control2 {
			r.controller2 = memberId
		} else if r.controller2 == memberId {
			r.controller2 = 0
		}
		return nil
	}
}

func (r *WebRTCRoomSession) renderCallback(p *ppu.PPU) {
	p.Render()
	data, release, err := r.videoEncoder.Encode(p.Frame())
	if err != nil {
		log.Println("encoder error: ", err)
	}
	defer release()
	for _, conn := range r.connections {
		// peer conn没有建立连接不编码视频，否则会导致I帧没有发送到客户端
		if !conn.connected.Load() {
			continue
		}
		if err := conn.videoTrack.WriteSample(media.Sample{
			Data:     data,
			Duration: 2 * time.Millisecond, // todo 根据帧率设置Duration
		}); err != nil {
			log.Println("write video sample error:", err)
		}
		release()
	}
}

// onNewConnection 新websocket连接建立后， 创建webrtc连接
// 1. 创建webrtc连接，创建视频音频流track
// 1.5 创建turn服务器凭证，保存在redis中，将用户名密码对发送给客户端
// 2. 创建SDP Offer， 并发送给客户端，完成SDP协商
// 3. 创建DataChannel，设置Message回调和连接状态回调
// 4. 创建视频编码器， 每个连接使用独立的视频编码器
// 5. 转移模拟器控制权
func (r *WebRTCRoomSession) onNewConnection(ctx context.Context, wsConn *WebsocketConn) {
	const videoCodec = "h264"
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
	roomConnection := &RoomConn{
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
	if err := roomConnection.SendTurnServerInfo(); err != nil {
		panic(fmt.Errorf("unable to send turn server info, error: %w", err))
	}
	// 创建H264视频和opus音频流
	videoTrack, _ := rtcFactory.VideoTrack(videoCodec)
	if _, err := peer.AddTrack(videoTrack); err != nil {
		panic(fmt.Errorf("unable to add video track to peer, error: %w", err))
	}
	roomConnection.videoTrack = videoTrack
	audioTrack, _ := rtcFactory.AudioTrack("opus")
	if _, err := peer.AddTrack(audioTrack); err != nil {
		panic(fmt.Errorf("unable to add audio track to peer, error: %w", err))
	}
	roomConnection.audioTrack = audioTrack
	// 创建DataChannel，用于控制消息传递
	channel, err := peer.CreateDataChannel("control-channel", nil)
	if err != nil {
		panic(fmt.Errorf("unable to create control data channel, error: %w", err))
	}
	roomConnection.dataChannel = channel
	// 创建SDP Offer， 并设置LocalSDP
	sdp, err := peer.CreateOffer(nil)
	if err != nil {
		panic(fmt.Errorf("create sdp offer error: %w", err))
	}
	if err := peer.SetLocalDescription(sdp); err != nil {
		panic(fmt.Errorf("unable to set local description, error: %w", err))
	}
	// 发送SDP Offer到客户端
	data, _ := json.Marshal(sdp)
	if err := roomConnection.sendMessage(Message{
		Type: MessageSDPOffer,
		Data: string(data),
	}); err != nil {
		panic(fmt.Errorf("unable to send sdp offer, error: %w", err))
	}
	// DataChannel OnMessage
	channel.OnMessage(func(msg webrtc.DataChannelMessage) {
		roomConnection.OnDataChannelMessage(msg, r.dataChannelMsgChan)
	})
	// PeerConnection connection state change
	peer.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		roomConnection.onPeerConnectionState(state, r.signalChan)
	})
	peer.OnICEConnectionStateChange(roomConnection.onICEStateChange)

	controllableMembers := r.filterMembers(isControllableMember)
	r.members[wsConn.Member.UserId] = wsConn.Member
	r.connections[wsConn.Member.UserId] = roomConnection
	// 如果当前没有可控制游戏的玩家，且当前连接是可控制游戏的玩家，转移控制权
	if r.controller1 == 0 && len(controllableMembers) == 0 && wsConn.Member.Role == room.RoleHost {
		_ = r.transferControl(wsConn.Member.UserId, true, false)
	}

	go func() {
		roomConnection.Handle(context.WithoutCancel(ctx))
		r.signalChan <- Signal{
			Type: SignalWebsocketClose,
			Data: roomConnection.wsConn,
		}
	}()
}

func (r *WebRTCRoomSession) onWebsocketConnClose(wsConn *WebsocketConn) {

}

// audioSampleListener 模拟器的音频输出到chan中，该循环从channel读取音频数据，并按照采样率打包发送
func (r *WebRTCRoomSession) audioSampleListener(ctx context.Context) {
	// 每5毫秒发送一次，根据采样率计算buffer大小
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

func (r *WebRTCRoomSession) sendAudioSamples(samples []float32) {
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

func (r *WebRTCRoomSession) filterMembers(filterFunc func(room.Member) bool) []*room.Member {
	result := make([]*room.Member, 0, len(r.members))
	for _, m := range r.members {
		if filterFunc(*m) {
			result = append(result, m)
		}
	}
	return result
}

func isControllableMember(m room.Member) bool {
	return m.Role != room.RoleObserver
}

func (r *WebRTCRoomSession) Close() {
	r.cancel()
}

func (r *WebRTCRoomSession) onClose() {
	r.emulatorCancel()
	for _, conn := range r.connections {
		conn.Close()
	}
}

func (r *WebRTCRoomSession) KickMember(memberId int64) error {
	f := future.NewFuture[struct{}]()
	r.signalChan <- Signal{
		Type: SignalKickMember,
		Data: kickMemberRequest{
			memberId: memberId,
			future:   f,
		},
	}
	_, err := f.Result()
	return err
}

func (r *WebRTCRoomSession) kickMember(memberId int64) error {
	if conn, ok := r.connections[memberId]; ok {
		conn.Close()
	}
	delete(r.connections, memberId)
	delete(r.members, memberId)
	return nil
}

func (r *WebRTCRoomSession) AlterRole(memberId int64, role byte) error {
	f := future.NewFuture[struct{}]()
	r.signalChan <- Signal{
		Type: SignalAlterRole,
		Data: alterRoleRequest{
			memberId: memberId,
			role:     role,
			future:   f,
		},
	}
	_, err := f.Result()
	return err
}

func (r *WebRTCRoomSession) alterRole(memberId int64, role byte) error {
	m, ok := r.members[memberId]
	if ok {
		m.Role = role
	}
	return nil
}
