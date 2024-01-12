package service

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
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
	controlChan    chan *ControlMessage // 模拟器输入channel
	game           string
	e              *emulator.Emulator
	emulatorCancel context.CancelFunc
}

func newRoomSession(roomId int64, game string) *RoomSession {
	conf := config.GetEmulatorConfig()
	game = filepath.Join(conf.GameDirectory, game)
	rs := &RoomSession{
		game:           game,
		connections:    make(map[string]*RoomConnection),
		newConnChan:    make(chan *RoomConnection, 16),
		closedConnChan: make(chan *RoomConnection, 16),
		controlChan:    make(chan *ControlMessage, 128),
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
		case <-rs.controlChan:
			// todo handle input messages
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
			go conn.HandleRead(rs.OnConnectionClose)
		}
	}
}

func (rs *RoomSession) OnConnectionClose(c *RoomConnection) {
	rs.closedConnChan <- c
}

func (rc *RoomConnection) HandleRead(closeCallback func(*RoomConnection)) {
	for {
		msgType, _, err := rc.conn.ReadMessage()
		if err != nil {
			log.Println("ws conn read error", err)
			closeCallback(rc)
			return
		}
		// todo handle inputs
		switch msgType {
		case websocket.BinaryMessage:
		case websocket.TextMessage:
		}
	}
}

func (rs *RoomSession) emulatorRenderCallback(p *ppu.PPU) {
	p.Render()
	data := p.CompressedFrameData()
	rs.writeChan <- data
}
