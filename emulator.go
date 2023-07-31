package main

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/cartridge"
	"github.com/stellarisJAY/nesgo/cpu"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/stellarisJAY/nesgo/trace"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"os"
	"time"
	"unsafe"
)

// Emulator 模拟器
type Emulator struct {
	processor *cpu.Processor
	cartridge cartridge.Cartridge
	bus       *bus.Bus
	ppu       *ppu.PPU
	joyPad    *bus.JoyPad

	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture

	keyMap map[sdl.Scancode]bus.JoyPadButton

	config Config
}

func NewEmulator(nesData []byte, config Config) *Emulator {
	e := &Emulator{
		cartridge: cartridge.MakeCartridge(nesData),
		config:    config,
	}
	e.ppu = ppu.NewPPU(byte(e.cartridge.GetMirroring()), e.cartridge.GetChrBank, e.cartridge.GetMirroring)
	e.joyPad = bus.NewJoyPad()
	e.bus = bus.NewBus(e.cartridge, e.ppu, e.RendererCallback, e.joyPad)
	e.processor = cpu.NewProcessorWithROM(e.bus)
	e.keyMap = make(map[sdl.Scancode]bus.JoyPadButton)
	scale := int32(e.config.scale)
	window, renderer, err := initSDL(scale)
	if err != nil {
		panic(err)
	}
	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, ppu.WIDTH, ppu.HEIGHT)
	e.window, e.renderer, e.texture = window, renderer, texture
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

func (e *Emulator) LoadAndRun(enableTrace bool) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
			e.onShutdown()
		}
	}()
	e.loadKeyMap()
	if enableTrace {
		e.processor.LoadAndRunWithCallback(trace.Trace)
	} else {
		e.processor.LoadAndRunWithCallback(func(_ *cpu.Processor, _ *cpu.Instruction) {})
	}
}

func (e *Emulator) Disassemble() {
	e.processor.Disassemble(trace.PrintDisassemble)
}

func (e *Emulator) RendererCallback(p *ppu.PPU) {
	p.Render()
	frame := p.FrameData()
	_ = e.texture.Update(nil, unsafe.Pointer(&frame[0]), ppu.WIDTH*3)
	_ = e.renderer.Copy(e.texture, nil, nil)
	e.renderer.Present()
	e.handleEvents()
	time.Sleep(e.config.frameInterval)
	e.handleEvents()
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
			switch event.State {
			case sdl.PRESSED:
				if button, ok := e.keyMap[event.Keysym.Scancode]; ok {
					e.joyPad.SetButtonPressed(button, true)
				}
			case sdl.RELEASED:
				if button, ok := e.keyMap[event.Keysym.Scancode]; ok {
					e.joyPad.SetButtonPressed(button, false)
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
