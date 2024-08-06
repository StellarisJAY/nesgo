//go:build sdl

package emulator

import (
	"context"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/stellarisJAY/nesgo/emulator/apu"
	"github.com/stellarisJAY/nesgo/emulator/bus"
	"github.com/stellarisJAY/nesgo/emulator/cartridge"
	"github.com/stellarisJAY/nesgo/emulator/config"
	"github.com/stellarisJAY/nesgo/emulator/cpu"
	"github.com/stellarisJAY/nesgo/emulator/ppu"
	"github.com/stellarisJAY/nesgo/emulator/trace"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// Emulator sdl render emulator
type Emulator struct {
	RawEmulator
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	keyMap   map[sdl.Scancode]bus.JoyPadButton
	audio    *Audio
}

func (e *Emulator) init() {
	s, err := os.Stat(e.config.SaveDirectory)
	if os.IsNotExist(err) {
		err := os.MkdirAll(e.config.SaveDirectory, 0755)
		if err != nil {
			log.Println(err)
			panic("Unable to find or create save directory")
		}
	} else if err != nil {
		panic("Unable to find save directory")
	} else if !s.IsDir() {
		panic("Provided save directory is not a directory")
	}
}

func NewEmulator(nesData []byte, conf config.Config) *Emulator {
	c, err := cartridge.MakeCartridge(nesData)
	if err != nil {
		panic(err)
	}
	e := &Emulator{
		RawEmulator: RawEmulator{
			cartridge: c,
			config:    conf,
			m:         &sync.Mutex{},
		},
	}
	e.joyPad1 = bus.NewJoyPad()
	e.joyPad2 = bus.NewJoyPad()
	e.ppu = ppu.NewPPU(e.cartridge.GetChrBank, e.cartridge.GetMirroring, e.cartridge.WriteCHR)
	e.apu = apu.NewBasicAPU()
	e.bus = bus.NewBus(e.cartridge, e.ppu, e.RendererCallback, e.joyPad1, e.joyPad2, e.apu)
	e.processor = cpu.NewProcessor(e.bus)
	e.keyMap = make(map[sdl.Scancode]bus.JoyPadButton)
	if !conf.MuteApu {
		if err := portaudio.Initialize(); err != nil {
			panic(err)
		}
		e.audio = NewAudio()
		if err := e.audio.Start(); err != nil {
			panic(err)
		}
		e.apu.SetRates(bus.CPUFrequency, e.audio.sampleRate)
		e.apu.SetMemReader(e.bus.ReadMemUint8)
		e.apu.SetOutputChan(e.audio.sampleChan)
	} else {
		e.apu.Mute()
	}

	scale := int32(e.config.Scale)
	window, renderer, err := initSDL(scale)
	if err != nil {
		panic(err)
	}
	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, ppu.WIDTH, ppu.HEIGHT)
	e.window, e.renderer, e.texture = window, renderer, texture
	e.init()
	return e
}

func initSDL(scale int32) (*sdl.Window, *sdl.Renderer, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, nil, fmt.Errorf("init sdl error %w", err)
	}
	w, err := sdl.CreateWindow("NesGO", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, ppu.WIDTH*scale, ppu.HEIGHT*scale, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, nil, fmt.Errorf("sdl create window error %w", err)
	}
	r, err := sdl.CreateRenderer(w, 0, sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, nil, fmt.Errorf("sdl get renderer error %w", err)
	}
	_ = r.SetScale(float32(scale), float32(scale))
	return w, r, nil
}

// loadKeyMap 加载key绑定
func (e *Emulator) loadKeyMap() {
	// todo 配置文件
	e.keyMap[sdl.SCANCODE_W] = bus.Up
	e.keyMap[sdl.SCANCODE_S] = bus.Down
	e.keyMap[sdl.SCANCODE_A] = bus.Left
	e.keyMap[sdl.SCANCODE_D] = bus.Right

	e.keyMap[sdl.SCANCODE_SPACE] = bus.ButtonA
	e.keyMap[sdl.SCANCODE_F] = bus.ButtonB
	e.keyMap[sdl.SCANCODE_TAB] = bus.Select
	e.keyMap[sdl.SCANCODE_RETURN] = bus.Start
}

func (e *Emulator) LoadAndRun(ctx context.Context, enableTrace bool) {
	defer func() {
		if err := recover(); err != nil {
			e.onShutdown()
			panic(err)
		}
	}()
	e.loadKeyMap()
	if enableTrace {
		e.processor.LoadAndRunWithCallback(ctx, trace.Trace, nil)
	} else {
		e.processor.LoadAndRunWithCallback(ctx, nil, nil)
	}
}

func (e *Emulator) RendererCallback(p *ppu.PPU) {
	e.handleEvents()
	p.Render()
	frame := p.FrameData()
	_ = e.texture.Update(nil, unsafe.Pointer(&frame[0]), ppu.WIDTH*3)
	_ = e.renderer.Copy(e.texture, nil, nil)
	e.renderer.Present()
}

// readLatestSave 从存档文件夹读取当前游戏的最新存档文件
func (e *Emulator) readLatestSave() ([]byte, error) {
	dir := e.config.SaveDirectory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	currentGame := filepath.Base(e.config.Game)
	latestTime, latestName := time.UnixMilli(0), ""
	for _, entry := range entries {
		name := entry.Name()
		// 过滤文件名
		if strings.HasSuffix(name, ".save") && strings.HasPrefix(name, currentGame) {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			// 根据时间选最新文件
			if info.ModTime().After(latestTime) {
				latestTime, latestName = info.ModTime(), name
			}
		}
	}
	if latestName == "" {
		return nil, os.ErrNotExist
	}
	return os.ReadFile(filepath.Join(e.config.SaveDirectory, latestName))
}

func (e *Emulator) handleEvents() {
	for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
		switch ev.(type) {
		case *sdl.QuitEvent:
			e.onShutdown()
			os.Exit(0)
			return
		case *sdl.KeyboardEvent:
			event := ev.(*sdl.KeyboardEvent)
			if event.Keysym.Scancode == sdl.SCANCODE_ESCAPE {
				e.onShutdown()
				os.Exit(0)
				return
			}
			if event.Keysym.Scancode == sdl.SCANCODE_F1 && event.State == sdl.RELEASED {
				e.BoostCPU(0.5)
				continue
			}
			if event.Keysym.Scancode == sdl.SCANCODE_F2 && event.State == sdl.RELEASED {
				e.BoostCPU(-0.5)
				continue
			}
			// F5 快速存档
			if event.Keysym.Scancode == sdl.SCANCODE_F5 && event.State == sdl.RELEASED {
				if err := e.SaveToFile(); err != nil {
					log.Println("save game error:", err)
				}
				continue
			}
			// F8 加载最新存档
			if event.Keysym.Scancode == sdl.SCANCODE_F8 && event.State == sdl.RELEASED {
				data, err := e.readLatestSave()
				if err != nil {
					log.Println("can't load latest saved game, error:", err)
					continue
				}
				if err := e.Load(data); err != nil {
					log.Println("can't load saved game, error:", err)
				}
			}
			switch event.State {
			case sdl.PRESSED:
				if button, ok := e.keyMap[event.Keysym.Scancode]; ok {
					e.joyPad1.SetButtonPressed(button, true)
				}
			case sdl.RELEASED:
				if button, ok := e.keyMap[event.Keysym.Scancode]; ok {
					e.joyPad1.SetButtonPressed(button, false)
				}
			}
		default:
		}
	}
}

func (e *Emulator) onShutdown() {
	log.Println("shutting down nesGo emulator")
	_ = e.texture.Destroy()
	_ = e.renderer.Destroy()
	_ = e.window.Destroy()
}
