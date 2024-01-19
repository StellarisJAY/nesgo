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

type RoomSignal struct {
	signalNumber byte
	payload      []byte
	callback     func(error)
}

const (
	SignalClose byte = iota
	SignalRestart
)

type RoomSession struct {
	connections      map[string]*RoomConnection
	newConnChan      chan *RoomConnection
	closeChan        chan RoomSignal
	writeChan        chan []byte // 模拟器输出channel
	closedConnChan   chan *RoomConnection
	controlChan      chan *MessageWrapper     // 模拟器输入channel
	changeMemberChan chan MemberChangeMessage // 房间成员发生变化channel，用于修改权限和删除成员连接
	game             string
	e                *emulator.Emulator
	emulatorCancel   context.CancelFunc
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

type MemberChangeMessage struct {
	action MemberControlAction
	m      *room.Member
}

type MemberControlAction byte

const (
	DeleteMember = iota
	ChangeMemberType
)

func newRoomSession(roomId int64, game string) (*RoomSession, error) {
	conf := config.GetEmulatorConfig()
	game = filepath.Join(conf.GameDirectory, game)
	rs := &RoomSession{
		game:             game,
		connections:      make(map[string]*RoomConnection),
		newConnChan:      make(chan *RoomConnection, 16),
		closedConnChan:   make(chan *RoomConnection, 16),
		controlChan:      make(chan *MessageWrapper, 128),
		closeChan:        make(chan RoomSignal),
		writeChan:        make(chan []byte, 32),
		changeMemberChan: make(chan MemberChangeMessage, 16),
	}
	e, err := emulator.NewEmulator(game, conf, rs.emulatorRenderCallback)
	if err != nil {
		return nil, err
	}
	rs.e = e
	return rs, nil
}

func (rs *RoomSession) StartGame() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	rs.emulatorCancel = cancelFunc
	go rs.e.LoadAndRun(ctx, false)
	// 模拟器初始暂停，等待第一个websocket连接来resume
	rs.e.Pause()
}

func (rs *RoomSession) Restart(game string) error {
	result := make(chan error)
	rs.closeChan <- RoomSignal{
		signalNumber: SignalRestart,
		payload:      []byte(game),
		callback: func(err error) {
			result <- err
		},
	}
	return <-result
}

func (rs *RoomSession) restart(game string) error {
	conf := config.GetEmulatorConfig()
	game = filepath.Join(conf.GameDirectory, game)
	e, err := emulator.NewEmulator(game, conf, rs.emulatorRenderCallback)
	if err != nil {
		return err
	}
	rs.emulatorCancel()
	rs.e = e
	ctx, cancelFunc := context.WithCancel(context.Background())
	rs.emulatorCancel = cancelFunc
	go rs.e.LoadAndRun(ctx, false)
	return nil
}

func (rs *RoomSession) Close() {
	close(rs.closeChan)
}

func (rs *RoomSession) ControlLoop() {
	for {
		select {
		case signal := <-rs.closeChan:
			if signal.signalNumber == SignalClose {
				if rs.emulatorCancel != nil {
					rs.emulatorCancel()
				}
				return
			} else if signal.signalNumber == SignalRestart {
				err := rs.restart(string(signal.payload))
				signal.callback(err)
			}
		case conn := <-rs.closedConnChan:
			log.Println("conn closed, addr:", conn.conn.RemoteAddr().String())
			delete(rs.connections, conn.conn.RemoteAddr().String())
			// 当前房间没有活跃的控制连接，暂停模拟器
			livingControlConnections := rs.filterConnections(func(conn *RoomConnection) bool {
				return conn.m.MemberType == room.MemberTypeOwner || conn.m.MemberType == room.MemberTypeGamer
			})
			if len(livingControlConnections) == 0 {
				rs.e.Pause()
			}
		case msg := <-rs.changeMemberChan:
			conn, ok := rs.findMemberConnection(msg.m.UserId)
			if !ok {
				continue
			}
			switch msg.action {
			case ChangeMemberType:
				conn.m = msg.m
			case DeleteMember:
				delete(rs.connections, conn.conn.RemoteAddr().String())
				if err := conn.conn.Close(); err != nil {
					log.Println("kick member close conn error:", err)
				}
			default:
			}
		case msg := <-rs.controlChan:
			// ignore message from watcher
			if msg.from.MemberType >= room.MemberTypeWatcher {
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
			controlConnections := rs.filterConnections(func(conn *RoomConnection) bool {
				return conn.m.MemberType == room.MemberTypeOwner || conn.m.MemberType == room.MemberTypeGamer
			})
			if len(controlConnections) > 0 {
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

func (rs *RoomSession) filterConnections(filterer func(conn *RoomConnection) bool) []*RoomConnection {
	connections := make([]*RoomConnection, 0, len(rs.connections))
	for _, conn := range rs.connections {
		if filterer(conn) {
			connections = append(connections, conn)
		}
	}
	return connections
}

func (rs *RoomSession) findMemberConnection(id int64) (*RoomConnection, bool) {
	for _, conn := range rs.connections {
		if conn.m.UserId == id {
			return conn, true
		}
	}
	return nil, false
}
