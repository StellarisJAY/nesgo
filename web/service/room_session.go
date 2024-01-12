package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/stellarisJAY/nesgo/web/model/room"
	"log"
	"path/filepath"
)

type RoomConnection struct {
	conn *websocket.Conn
	m    *room.Member
}

type RoomSession struct {
	connections    map[string]*RoomConnection
	newConnChan    chan *RoomConnection
	closeChan      chan struct{}
	writeChan      chan []byte // 模拟器输出channel
	closedConnChan chan *RoomConnection
	controlChan    chan *MessageWrapper // 模拟器输入channel
	game           string
	e              *emulator.Emulator
	emulatorCancel context.CancelFunc
}

const (
	ActionButtonDown int = iota
	ActionButtonUp
)

type ControlMessage struct {
	KeyCode string `json:"KeyCode"`
	Action  int    `json:"Action"`
}

type MessageWrapper struct {
	ControlMessage
	fromAddr string
	from     *room.Member
}

func newRoomSession(roomId int64, game string) *RoomSession {
	conf := config.GetEmulatorConfig()
	game = filepath.Join(conf.GameDirectory, game)
	rs := &RoomSession{
		game:           game,
		connections:    make(map[string]*RoomConnection),
		newConnChan:    make(chan *RoomConnection, 16),
		closedConnChan: make(chan *RoomConnection, 16),
		controlChan:    make(chan *MessageWrapper, 128),
		closeChan:      make(chan struct{}),
		writeChan:      make(chan []byte, 32),
	}
	rs.e = emulator.NewEmulator(game, conf, rs.emulatorRenderCallback)
	return rs
}

func (rs *RoomSession) StartGame() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	rs.emulatorCancel = cancelFunc
	go rs.e.LoadAndRun(ctx, false)
	// 模拟器初始暂停，等待第一个websocket连接来resume
	rs.e.Pause()
}

func (rs *RoomSession) Close() {
	close(rs.closeChan)
}

func (rs *RoomSession) ControlLoop() {
	for {
		select {
		case <-rs.closeChan:
			if rs.emulatorCancel != nil {
				rs.emulatorCancel()
			}
			return
		case conn := <-rs.closedConnChan:
			log.Println("conn closed, addr:", conn.conn.RemoteAddr().String())
			delete(rs.connections, conn.conn.RemoteAddr().String())
			// 当前房间没有活跃连接，暂停模拟器
			if len(rs.connections) == 0 {
				rs.e.Pause()
			}
		case msg := <-rs.controlChan:
			// ignore message from watcher
			if msg.from.MemberType == room.MemberTypeWatcher {
				continue
			}
			// ignore message from dead connections
			if _, ok := rs.connections[msg.fromAddr]; !ok {
				continue
			}
			input := msg.ControlMessage
			switch input.KeyCode {
			case "KeyA":
				rs.e.SetJoyPadButtonPressed(bus.Left, input.Action == ActionButtonDown)
			case "KeyD":
				rs.e.SetJoyPadButtonPressed(bus.Right, input.Action == ActionButtonDown)
			case "KeyW":
				rs.e.SetJoyPadButtonPressed(bus.Up, input.Action == ActionButtonDown)
			case "KeyS":
				rs.e.SetJoyPadButtonPressed(bus.Down, input.Action == ActionButtonDown)
			case "Space":
				rs.e.SetJoyPadButtonPressed(bus.ButtonA, input.Action == ActionButtonDown)
			case "KeyJ":
				rs.e.SetJoyPadButtonPressed(bus.ButtonB, input.Action == ActionButtonDown)
			case "Enter":
				rs.e.SetJoyPadButtonPressed(bus.Start, input.Action == ActionButtonDown)
			default:
			}
		case data := <-rs.writeChan:
			for addr, c := range rs.connections {
				err := c.conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
					log.Println("ws write to", addr, "error:", err)
				}
			}
		case conn := <-rs.newConnChan:
			rs.connections[conn.conn.RemoteAddr().String()] = conn
			// 新连接建立，判断是否需要启动模拟器
			if len(rs.connections) == 1 {
				rs.e.Resume()
			}
			go conn.HandleRead(rs.OnConnectionClose, rs.controlChan)
		}
	}
}

func (rs *RoomSession) OnConnectionClose(c *RoomConnection) {
	rs.closedConnChan <- c
}

func (rc *RoomConnection) HandleRead(closeCallback func(*RoomConnection), controlChan chan *MessageWrapper) {
	for {
		msgType, payload, err := rc.conn.ReadMessage()
		if err != nil {
			log.Println("ws conn read error", err)
			closeCallback(rc)
			return
		}
		// todo handle inputs
		switch msgType {
		case websocket.BinaryMessage:
		case websocket.TextMessage:
			var msg ControlMessage
			err := json.Unmarshal(payload, &msg)
			if err != nil {
				log.Println("unmarshal control msg error:", err)
				continue
			}
			controlChan <- &MessageWrapper{
				ControlMessage: msg,
				fromAddr:       rc.conn.RemoteAddr().String(),
				from:           rc.m,
			}
		}
	}
}

func (rs *RoomSession) emulatorRenderCallback(p *ppu.PPU) {
	p.Render()
	data := p.CompressedFrameData()
	rs.writeChan <- data
}
