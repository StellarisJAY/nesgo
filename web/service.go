package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisJAY/nesgo/ppu"
	"log"
	"math/rand"
	"path/filepath"
	"sync"
)

type GameService struct {
	mutex    *sync.Mutex
	sessions map[string]*GameSession
	conf     config.Config
}

type GameSession struct {
	id       string
	conn     *websocket.Conn
	e        *emulator.Emulator
	conf     config.Config
	gameName string

	m         *sync.Mutex
	keyInputs []ControlMessage
	cancel    context.CancelFunc
}

const (
	ActionButtonDown int = iota
)

type ControlMessage struct {
	KeyCode string `json:"KeyCode"`
	Action  int    `json:"Action"`
}

func (s *GameService) HandleGamePage(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				c.JSON(500, fmt.Sprintf("{\"error\":\"%s\"}", e.Error()))
			}
		}
	}()
	id := instanceId()
	gameName := c.Param("name")
	session := makeGameSession(id, gameName, gameService.conf)
	s.mutex.Lock()
	s.sessions[id] = session
	s.mutex.Unlock()
	log.Println("game session created, id:", id)
	c.HTML(200, "game.html", gin.H{"game": gameName, "id": id})
}

func (s *GameService) HandleWebsocket(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				c.JSON(500, fmt.Sprintf("{\"error\":\"%s\"}", e.Error()))
			}
		}
	}()
	id := c.Param("id")
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		panic(err)
	}
	s.mutex.Lock()
	session := s.sessions[id]
	s.mutex.Unlock()

	session.conn = conn
	ctx, cancelFunc := context.WithCancel(context.Background())
	session.cancel = cancelFunc

	conn.SetCloseHandler(func(code int, text string) error {
		session.onConnectionClose()
		s.mutex.Lock()
		delete(s.sessions, id)
		s.mutex.Unlock()
		return nil
	})
	session.startGame(ctx)
}

func makeGameSession(id, game string, conf config.Config) *GameSession {
	gameFile := filepath.Join(conf.GameDirectory, game)
	g := &GameSession{
		id:        id,
		conf:      conf,
		gameName:  game,
		keyInputs: make([]ControlMessage, 0),
		m:         &sync.Mutex{},
	}
	g.e = emulator.NewEmulator(gameFile, conf, g.renderCallback)
	return g
}

func (g *GameSession) startGame(ctx context.Context) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok {
					log.Println("emulator error, session:", g.id, "error: ", err)
					g.cancel()
				}
			}
		}()
		g.e.LoadAndRun(ctx, false)
		log.Println("emulator closed, session:", g.id)
	}()
	go func() {
		g.readLoop()
		log.Println("ws read loop closed, session:", g.id)
	}()
}

func (g *GameSession) onConnectionClose() {
	g.cancel()
}

// GameService renderCallback 渲染回调，发送画面给前端
func (g *GameSession) renderCallback(p *ppu.PPU) {
	g.pullInput()
	p.Render()
	_ = g.conn.WriteMessage(websocket.BinaryMessage, p.CompressedFrameData())
}

func (g *GameSession) pullInput() {
	g.m.Lock()
	for len(g.keyInputs) > 0 {
		input := g.keyInputs[0]
		switch input.KeyCode {
		case "KeyA":
			g.e.SetJoyPadButtonPressed(bus.Left, input.Action == ActionButtonDown)
		case "KeyD":
			g.e.SetJoyPadButtonPressed(bus.Right, input.Action == ActionButtonDown)
		case "KeyW":
			g.e.SetJoyPadButtonPressed(bus.Up, input.Action == ActionButtonDown)
		case "KeyS":
			g.e.SetJoyPadButtonPressed(bus.Down, input.Action == ActionButtonDown)
		case "Space":
			g.e.SetJoyPadButtonPressed(bus.ButtonA, input.Action == ActionButtonDown)
		case "KeyJ":
			g.e.SetJoyPadButtonPressed(bus.ButtonB, input.Action == ActionButtonDown)
		case "Enter":
			g.e.SetJoyPadButtonPressed(bus.Start, input.Action == ActionButtonDown)
		default:
		}
		g.keyInputs = g.keyInputs[1:]
	}
	g.m.Unlock()
}

func (g *GameSession) readLoop() {
	for {
		msgType, data, err := g.conn.ReadMessage()
		if err != nil {
			log.Println("ws read error, session:", g.id, "error: ", err)
			g.cancel()
			return
		}
		switch msgType {
		case websocket.BinaryMessage:
		case websocket.TextMessage:
			msg := ControlMessage{}
			_ = json.Unmarshal(data, &msg)
			g.m.Lock()
			g.keyInputs = append(g.keyInputs, msg)
			g.m.Unlock()
		default:
		}
	}
}

func instanceId() string {
	return fmt.Sprintf("%02x%02x%02x%02x",
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256))
}
