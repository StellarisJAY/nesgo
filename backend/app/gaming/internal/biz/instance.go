package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/pkg/codec"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisJAY/nesgo/emulator/bus"
	"github.com/stellarisJAY/nesgo/emulator/config"
	"github.com/stellarisJAY/nesgo/emulator/ppu"
	"sync"
	"sync/atomic"
	"time"
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
	e               *emulator.Emulator
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
}

func (g *GameInstance) RenderCallback(ppu *ppu.PPU, logger *log.Helper) {
	ppu.Render()
	frame := ppu.Frame()
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
	defer g.mutex.Unlock()
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
		return ConsumerResult{Error: fmt.Errorf("restart emulator error: %v", err)}
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
		g.RenderCallback(p, log.NewHelper(log.With(log.DefaultLogger, "module", "emulator")))
	}
	e, err := emulator.NewEmulatorWithGameData(gameData, emulatorConfig, renderCallback, g.audioSampleChan, g.audioSampleRate)
	if err != nil {
		return fmt.Errorf("create new emulator error: %v", err)
	}
	g.e = e
	// 启动新模拟器goroutine
	ctx, cancelFunc := context.WithCancel(context.Background())
	g.emulatorCancel = cancelFunc
	go e.LoadAndRun(ctx, false)
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
