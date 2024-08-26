package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/pkg/codec"
	"github.com/stellarisJAY/nesgo/nes"
	"github.com/stellarisJAY/nesgo/nes/bus"
	"github.com/stellarisJAY/nesgo/nes/config"
	"github.com/stellarisJAY/nesgo/nes/ppu"
)

const (
	MsgPlayerControlButtonPressed byte = iota
	MsgPlayerControlButtonReleased
	MsgChat
	MsgNewConn
	MsgCloseConn
	MsgSetController1
	MsgSetController2
	MsgResetController
	MsgPeerConnected
	MsgPauseEmulator
	MsgResumeEmulator
	MsgSaveGame
	MsgLoadSave
	MsgRestartEmulator
	MsgPing
	MsgSetGraphicOptions
	MsgSetEmulatorSpeed
)

const (
	InstanceStatusRunning int32 = iota
	InstanceStatusStopped
)

const (
	MessageTargetEmulator int64 = iota
)

type ConsumerResult struct {
	Success bool
	Error   error
	Message string
	Data    any
}

type Message struct {
	Type      byte  `json:"type"`
	From      int64 `json:"from"`
	To        int64 `json:"to"`
	Timestamp int64 `json:"timestamp"`
	Data      any   `json:"data"`

	resultChan chan ConsumerResult
}

type GameInstance struct {
	RoomId          int64
	e               *nes.Emulator
	emulatorCancel  context.CancelFunc
	game            string
	videoEncoder    codec.IVideoEncoder
	audioSampleRate int
	audioEncoder    codec.IAudioEncoder
	audioSampleChan chan float32 // audioSampleChan 音频输出channel
	controller1     int64        // controller1 模拟器P1控制权玩家
	controller2     int64        // controller2 模拟器P2控制权玩家

	messageChan chan *Message
	cancel      context.CancelFunc

	connections map[int64]*Connection
	mutex       *sync.RWMutex

	destroyed bool
	LeaseID   int64
	status    atomic.Int32

	createTime time.Time
	InstanceId string

	allConnCloseCallback func(instance *GameInstance)

	enhancedFrame    *image.YCbCr
	enhanceFrameOpen bool
	frameEnhancer    func(frame *ppu.Frame) *ppu.Frame

	reverseColorOpen bool
	grayscaleOpen    bool
}

func (g *GameInstance) enhanceFrame(frame *ppu.Frame) *ppu.Frame {
	original := frame.YCbCr()
	enhancedFrame := g.enhancedFrame
	for y := 0; y < ppu.HEIGHT; y++ {
		for x := 0; x < ppu.WIDTH; x++ {
			// 分辨率放大到原来的两倍，每个像素变成四个像素
			offset := original.YOffset(x, y)
			cOffset := original.COffset(x, y)
			dx, dy := x*2, y*2
			// TODO 优化性能，减少CPU占用
			enhancedFrame.Y[enhancedFrame.YOffset(dx+1, dy+1)] = original.Y[offset]
			enhancedFrame.Y[enhancedFrame.YOffset(dx+1, dy)] = original.Y[offset]
			enhancedFrame.Y[enhancedFrame.YOffset(dx, dy+1)] = original.Y[offset]
			enhancedFrame.Y[enhancedFrame.YOffset(dx, dy)] = original.Y[offset]

			enhancedFrame.Cb[enhancedFrame.COffset(dx+1, dy+1)] = original.Cb[cOffset]
			enhancedFrame.Cb[enhancedFrame.COffset(dx+1, dy)] = original.Cb[cOffset]
			enhancedFrame.Cb[enhancedFrame.COffset(dx, dy+1)] = original.Cb[cOffset]
			enhancedFrame.Cb[enhancedFrame.COffset(dx, dy)] = original.Cb[cOffset]

			enhancedFrame.Cr[enhancedFrame.COffset(dx+1, dy+1)] = original.Cr[cOffset]
			enhancedFrame.Cr[enhancedFrame.COffset(dx+1, dy)] = original.Cr[cOffset]
			enhancedFrame.Cr[enhancedFrame.COffset(dx, dy+1)] = original.Cr[cOffset]
			enhancedFrame.Cr[enhancedFrame.COffset(dx, dy)] = original.Cr[cOffset]
		}
	}
	return ppu.NewCustomSizeFrame(enhancedFrame)
}

func (g *GameInstance) RenderCallback(ppu *ppu.PPU, logger *log.Helper) {
	frame := g.frameEnhancer(ppu.Frame())
	data, release, err := g.videoEncoder.Encode(frame)
	if err != nil {
		logger.Error("encode frame error", "err", err)
		return
	}
	defer release()
	sample := media.Sample{Data: data, Duration: 15 * time.Millisecond, Timestamp: time.Now()}
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	for _, conn := range g.connections {
		if conn.pc.ConnectionState() != webrtc.PeerConnectionStateConnected {
			continue
		}
		if err := conn.videoTrack.WriteSample(sample); err != nil {
			logger.Errorf("write sample error: %v", err)
		}
	}
	return
}

func (g *GameInstance) audioSampleListener(ctx context.Context, logger *log.Helper) {
	// 每5毫秒发送一次，根据采样率计算buffer大小
	buffer := make([]float32, 0, g.audioSampleRate*5/1000)
	for {
		select {
		case <-ctx.Done():
			return
		case s := <-g.audioSampleChan:
			buffer = append(buffer, s)
			if len(buffer) == cap(buffer) {
				g.sendAudioSamples(buffer, logger)
				buffer = buffer[:0]
			}
		}
	}
}

func (g *GameInstance) sendAudioSamples(buffer []float32, logger *log.Helper) {
	data, err := g.audioEncoder.Encode(buffer)
	if err != nil {
		logger.Error("encode audio samples error: ", err)
	}
	sample := media.Sample{Data: data, Timestamp: time.Now()}
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	for _, conn := range g.connections {
		if conn.pc.ConnectionState() != webrtc.PeerConnectionStateConnected {
			continue
		}
		err := conn.audioTrack.WriteSample(sample)
		if err != nil {
			logger.Errorf("write sample error: %v", err)
		}
	}
}

func (g *GameInstance) onDataChannelMessage(userId int64, msg *Message) {
	msg.From = userId
	g.messageChan <- msg
}

func (g *GameInstance) messageConsumer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-g.messageChan:
			switch msg.Type {
			case MsgPlayerControlButtonPressed:
				fallthrough
			case MsgPlayerControlButtonReleased:
				keyCode := msg.Data.(string)
				g.handlePlayerControl(keyCode, msg.Type, msg.From)
			case MsgNewConn:
				msg.resultChan <- g.handleMsgNewConn(msg.Data.(*Connection))
			case MsgPeerConnected:
				g.handlePeerConnected(msg.Data.(*Connection))
			case MsgCloseConn:
				g.handleMsgCloseConn(msg.Data.(*Connection))
			case MsgSetController1:
				msg.resultChan <- g.handleSetController(msg.Data.(int64), 0)
			case MsgSetController2:
				msg.resultChan <- g.handleSetController(msg.Data.(int64), 1)
			case MsgResetController:
				msg.resultChan <- g.handleResetController(msg.Data.(int64))
			case MsgSaveGame:
				msg.resultChan <- g.handleSaveGame()
			case MsgLoadSave:
				msg.resultChan <- g.handleLoadSave(msg.Data.(*gameSaveLoader))
			case MsgRestartEmulator:
				msg.resultChan <- g.handleRestartEmulator(msg.Data.(*emulatorRestartRequest))
			case MsgChat:
				g.handleChat(msg)
			case MsgSetGraphicOptions:
				msg.resultChan <- g.setGraphicOptions(msg.Data.(*GraphicOptions))
			case MsgSetEmulatorSpeed:
				msg.resultChan <- g.setEmulatorSpeed(msg.Data.(float64))
			default: // TODO handle unknown message
			}
		}
	}
}

func (g *GameInstance) onConnected(conn *Connection) {
	g.messageChan <- &Message{
		Type: MsgPeerConnected,
		Data: conn,
	}
}

func (g *GameInstance) closeConnection(conn *Connection) {
	g.messageChan <- &Message{
		Type: MsgCloseConn,
		Data: conn,
	}
}

func (g *GameInstance) filterConnection(filter func(*Connection) bool) []*Connection {
	result := make([]*Connection, 0, len(g.connections))
	for _, conn := range g.connections {
		if filter(conn) {
			result = append(result, conn)
		}
	}
	return result
}

func (g *GameInstance) handleMsgNewConn(conn *Connection) ConsumerResult {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if old, ok := g.connections[conn.userId]; ok {
		old.Close()
		delete(g.connections, conn.userId)
	}
	g.connections[conn.userId] = conn
	return ConsumerResult{Success: true}
}

func (g *GameInstance) handlePeerConnected(_ *Connection) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	active := g.filterConnection(func(c *Connection) bool {
		return c.pc.ConnectionState() == webrtc.PeerConnectionStateConnected
	})
	// 首个活跃连接，开启模拟器
	if len(active) == 1 {
		g.e.Resume()
	}
}

func (g *GameInstance) handleMsgCloseConn(conn *Connection) {
	g.mutex.Lock()
	// 被动关闭连接，可能是因为新连接挤掉旧连接，需要避免删除新连接
	if cur, ok := g.connections[conn.userId]; ok {
		if cur.pc.ConnectionState() == webrtc.PeerConnectionStateClosed ||
			cur.pc.ConnectionState() == webrtc.PeerConnectionStateFailed ||
			cur.pc.ConnectionState() == webrtc.PeerConnectionStateDisconnected {
			delete(g.connections, conn.userId)
		}
	}
	active := g.filterConnection(func(conn *Connection) bool {
		return conn.pc.ConnectionState() == webrtc.PeerConnectionStateConnected
	})
	// 在Pause之前必须释放连接列表锁，避免 模拟器goroutine和messageConsumer死锁
	// 死锁循环等待：模拟器RenderCallback等待获取g.mutex, 之后消费processor.channel(无缓冲通道)
	//             closeConn获取到了g.mutex, 之后向processor.channel发送消息。
	// 模拟器等待g.mutex, closeConn等待processor.channel，循环等待导致死锁
	g.mutex.Unlock()
	// 没有活跃连接，暂停模拟器
	if len(active) == 0 {
		g.e.Pause()
		g.allConnCloseCallback(g)
	}
}

func (g *GameInstance) handlePlayerControl(keyCode string, action byte, player int64) {
	if player != g.controller1 && player != g.controller2 {
		return
	}
	id := 1
	if player == g.controller2 {
		id = 2
	}
	switch keyCode {
	case "Up":
		g.e.SetJoyPadButtonPressed(id, bus.Up, action == MsgPlayerControlButtonPressed)
	case "Down":
		g.e.SetJoyPadButtonPressed(id, bus.Down, action == MsgPlayerControlButtonPressed)
	case "Left":
		g.e.SetJoyPadButtonPressed(id, bus.Left, action == MsgPlayerControlButtonPressed)
	case "Right":
		g.e.SetJoyPadButtonPressed(id, bus.Right, action == MsgPlayerControlButtonPressed)
	case "A":
		g.e.SetJoyPadButtonPressed(id, bus.ButtonA, action == MsgPlayerControlButtonPressed)
	case "B":
		g.e.SetJoyPadButtonPressed(id, bus.ButtonB, action == MsgPlayerControlButtonPressed)
	case "Select":
		g.e.SetJoyPadButtonPressed(id, bus.Select, action == MsgPlayerControlButtonPressed)
	case "Start":
		g.e.SetJoyPadButtonPressed(id, bus.Start, action == MsgPlayerControlButtonPressed)
	}
}

func (g *GameInstance) handleSetController(playerId int64, id int) ConsumerResult {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	if _, ok := g.connections[playerId]; !ok {
		return ConsumerResult{Success: false}
	}
	if id == 0 {
		g.controller1 = playerId
		if g.controller2 == playerId {
			g.controller2 = 0
		}
	} else {
		g.controller2 = playerId
		if g.controller1 == playerId {
			g.controller1 = 0
		}
	}
	return ConsumerResult{Success: true}
}

func (g *GameInstance) handleResetController(playerId int64) ConsumerResult {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	if _, ok := g.connections[playerId]; !ok {
		return ConsumerResult{Success: false}
	}
	if g.controller1 == playerId {
		g.controller1 = 0
	}
	if g.controller2 == playerId {
		g.controller2 = 0
	}
	return ConsumerResult{Success: true}
}

func (g *GameInstance) DumpStats() *GameInstanceStats {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	active := g.filterConnection(func(c *Connection) bool {
		return c.pc.ConnectionState() == webrtc.PeerConnectionStateConnected
	})
	return &GameInstanceStats{
		RoomId:            g.RoomId,
		Connections:       len(g.connections),
		ActiveConnections: len(active),
		Game:              g.game,
		Uptime:            time.Now().Sub(g.createTime),
	}
}

func (g *GameInstance) DeleteConnection(userId int64) {
	g.mutex.Lock()
	conn, ok := g.connections[userId]
	if ok {
		delete(g.connections, userId)
	}
	g.mutex.Unlock()
	conn.Close()
}

func (g *GameInstance) handleSaveGame() ConsumerResult {
	g.e.Pause()
	defer g.e.Resume()
	data, err := g.e.GetSaveData()
	success := err == nil
	return ConsumerResult{Success: success, Data: data, Error: err}
}

// handleLoadSave 加载存档，如果存档的游戏与当前模拟器运行的游戏不同，需要先重启模拟器
func (g *GameInstance) handleLoadSave(loader *gameSaveLoader) ConsumerResult {
	if g.game == loader.game {
		g.e.Pause()
		defer g.e.Resume()
		err := g.e.Load(loader.data)
		return ConsumerResult{Success: err == nil, Error: err}
	}
	// 加载存档需要切换游戏，从repo获取游戏数据
	data, err := loader.gameFileRepo.GetGameData(context.Background(), loader.game)
	if err != nil {
		return ConsumerResult{Error: err}
	}
	// 加载新游戏，重启模拟器
	err = g.restartEmulator(loader.game, data)
	if err != nil {
		return ConsumerResult{Error: fmt.Errorf("restart nes error: %v", err)}
	}
	// 模拟器加载存档数据
	g.e.Pause()
	defer g.e.Resume()
	err = g.e.Load(loader.data)
	return ConsumerResult{Success: err == nil, Error: err}
}

func (g *GameInstance) handleRestartEmulator(request *emulatorRestartRequest) ConsumerResult {
	err := g.restartEmulator(request.game, request.gameData)
	return ConsumerResult{Success: err == nil, Error: err}
}

func (g *GameInstance) restartEmulator(game string, gameData []byte) error {
	// 结束旧模拟器goroutine
	g.emulatorCancel()

	// 清空上一个模拟器输出的还未发送的音频
LOOP:
	for {
		select {
		case <-g.audioSampleChan:
		default:
			break LOOP
		}
	}

	g.game = game
	// 创建新模拟器
	emulatorConfig := config.Config{
		Game:               game,
		EnableTrace:        false,
		Disassemble:        false,
		MuteApu:            false,
		Debug:              false,
		SnapshotSerializer: "gob",
	}
	renderCallback := func(p *ppu.PPU) {
		g.RenderCallback(p, log.NewHelper(log.With(log.DefaultLogger, "module", "nes")))
	}
	e, err := nes.NewEmulatorWithGameData(gameData, emulatorConfig, renderCallback, g.audioSampleChan, g.audioSampleRate)
	if err != nil {
		return fmt.Errorf("create new nes error: %v", err)
	}
	g.e = e
	// 启动新模拟器goroutine
	ctx, cancelFunc := context.WithCancel(context.Background())
	g.emulatorCancel = cancelFunc
	go e.LoadAndRun(ctx, false)
	// emulatorCancel后会结束旧的音频发送器，需要重启音频发送goroutine
	go g.audioSampleListener(ctx, log.NewHelper(log.With(log.DefaultLogger, "module", "audioSender")))
	return nil
}

func (g *GameInstance) SaveGame() (*GameSave, error) {
	ch := make(chan ConsumerResult)
	g.messageChan <- &Message{Type: MsgSaveGame, resultChan: ch}
	result := <-ch
	close(ch)
	if result.Success {
		data := result.Data.([]byte)
		return &GameSave{
			Game:       g.game,
			Data:       data,
			CreateTime: time.Now().UnixMilli(),
		}, nil
	} else {
		return nil, result.Error
	}
}

type gameSaveLoader struct {
	data         []byte       // 存档数据
	game         string       // 存档对应的游戏
	gameFileRepo GameFileRepo // 游戏文件repo，如果需要重启模拟器，需要repo下载游戏文件
}

type emulatorRestartRequest struct {
	game     string
	gameData []byte
}

func (g *GameInstance) LoadSave(data []byte, game string, gameFileRepo GameFileRepo) error {
	ch := make(chan ConsumerResult)
	loader := &gameSaveLoader{data, game, gameFileRepo}
	g.messageChan <- &Message{Type: MsgLoadSave, Data: loader, resultChan: ch}
	result := <-ch
	close(ch)
	if result.Success {
		return nil
	} else {
		return result.Error
	}
}

func (g *GameInstance) RestartEmulator(game string, gameData []byte) error {
	ch := make(chan ConsumerResult)
	request := &emulatorRestartRequest{game, gameData}
	g.messageChan <- &Message{Type: MsgRestartEmulator, Data: request, resultChan: ch}
	result := <-ch
	close(ch)
	if result.Success {
		return nil
	} else {
		return result.Error
	}
}

func (g *GameInstance) handleChat(msg *Message) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	resp := &Message{Type: MsgChat, From: msg.From, To: 0, Data: msg.Data}
	for _, conn := range g.connections {
		raw, _ := json.Marshal(resp)
		_ = conn.dataChannel.Send(raw)
	}
}

func (g *GameInstance) SetGraphicOptions(options *GraphicOptions) {
	resultCh := make(chan ConsumerResult)
	g.messageChan <- &Message{Type: MsgSetGraphicOptions, Data: options, resultChan: resultCh}
	_ = <-resultCh
	close(resultCh)
}

func (g *GameInstance) setGraphicOptions(options *GraphicOptions) ConsumerResult {
	g.e.Pause()
	defer g.e.Resume()
	defer func() {
		options.HighResOpen = g.enhanceFrameOpen
		options.ReverseColor = g.reverseColorOpen
	}()
	if g.enhanceFrameOpen != options.HighResOpen {
		g.enhanceFrameOpen = options.HighResOpen
		g.videoEncoder.Close()
		enhanceRate := 1
		if options.HighResOpen {
			enhanceRate = 2
			g.frameEnhancer = g.enhanceFrame
		} else {
			g.frameEnhancer = func(f *ppu.Frame) *ppu.Frame {
				return f
			}
		}
		enc, err := codec.NewVideoEncoder("vp8", ppu.WIDTH*enhanceRate, ppu.HEIGHT*enhanceRate)
		if err != nil {
			return ConsumerResult{Error: err}
		}
		g.videoEncoder = enc
	}
	if g.reverseColorOpen != options.ReverseColor {
		g.reverseColorOpen = options.ReverseColor
		if g.reverseColorOpen {
			g.e.Frame().UseReverseColorPreprocessor()
		} else {
			g.e.Frame().RemoveReverseColorPreprocessor()
		}
	}

	if g.grayscaleOpen != options.Grayscale {
		g.grayscaleOpen = options.Grayscale
		if g.grayscaleOpen {
			g.e.Frame().UseGrayscalePreprocessor()
		} else {
			g.e.Frame().RemoveGrayscalePreprocessor()
		}
	}

	options.Grayscale = g.grayscaleOpen
	options.HighResOpen = g.enhanceFrameOpen
	options.ReverseColor = g.reverseColorOpen
	return ConsumerResult{Success: true, Error: nil}
}

func (g *GameInstance) SetEmulatorSpeed(boostRate float64) float64 {
	resultCh := make(chan ConsumerResult)
	g.messageChan <- &Message{Type: MsgSetEmulatorSpeed, Data: boostRate, resultChan: resultCh}
	result := <-resultCh
	close(resultCh)
	return result.Data.(float64)

}

func (g *GameInstance) setEmulatorSpeed(boostRate float64) ConsumerResult {
	g.e.Pause()
	defer g.e.Resume()
	rate := g.e.SetCPUBoostRate(boostRate)
	return ConsumerResult{Success: true, Error: nil, Data: rate}
}
