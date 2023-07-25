package main

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/cpu"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/stellarisJAY/nesgo/trace"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"os"
	"unsafe"
)

const SCALE = 4

type Emulator struct {
	processor *cpu.Processor
	rom       *bus.ROM
	bus       *bus.Bus
	ppu       *ppu.PPU
	joyPad    *bus.JoyPad
	window    *sdl.Window
	renderer  *sdl.Renderer
	texture   *sdl.Texture
}

func NewEmulator(nesData []byte) *Emulator {
	e := &Emulator{
		rom: bus.NewROM(nesData),
	}
	e.ppu = ppu.NewPPU(e.rom.GetChrROM(), byte(e.rom.GetMirroring()))
	e.joyPad = bus.NewJoyPad()
	e.bus = bus.NewBus(e.rom, e.ppu, e.RendererCallback, e.joyPad)
	e.processor = cpu.NewProcessorWithROM(e.bus)

	window, renderer, err := initSDL()
	if err != nil {
		panic(err)
	}
	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, ppu.WIDTH, ppu.HEIGHT)
	e.window, e.renderer, e.texture = window, renderer, texture
	return e
}

func initSDL() (*sdl.Window, *sdl.Renderer, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, nil, fmt.Errorf("init sdl error %w", err)
	}
	w, err := sdl.CreateWindow("NesGO", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, ppu.WIDTH*SCALE, ppu.HEIGHT*SCALE, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, nil, fmt.Errorf("sdl create window error %w", err)
	}
	r, err := sdl.CreateRenderer(w, 0, sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, nil, fmt.Errorf("sdl get renderer error %w", err)
	}
	_ = r.SetScale(SCALE, SCALE)
	return w, r, nil
}

func (e *Emulator) LoadAndRun(trace bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			e.onShutdown()
		}
	}()
	if trace {
		e.processor.LoadAndRunWithCallback(logInstruction)
	} else {
		e.processor.LoadAndRunWithCallback(func(_ *cpu.Processor, _ *cpu.Instruction) {})
	}
}

func (e *Emulator) RendererCallback(p *ppu.PPU) {
	p.Render()
	frame := p.FrameData()
	_ = e.texture.Update(nil, unsafe.Pointer(&frame[0]), ppu.WIDTH*3)
	_ = e.renderer.Copy(e.texture, nil, nil)
	e.renderer.Present()
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
			switch event.Keysym.Scancode {
			case sdl.SCANCODE_UP:
				e.joyPad.SetButtonPressed(bus.Up, event.State == sdl.PRESSED)
			case sdl.SCANCODE_DOWN:
				e.joyPad.SetButtonPressed(bus.Down, event.State == sdl.PRESSED)
			case sdl.SCANCODE_LEFT:
				e.joyPad.SetButtonPressed(bus.Left, event.State == sdl.PRESSED)
			case sdl.SCANCODE_RIGHT:
				e.joyPad.SetButtonPressed(bus.Right, event.State == sdl.PRESSED)
			case sdl.SCANCODE_A:
				e.joyPad.SetButtonPressed(bus.ButtonA, event.State == sdl.PRESSED)
			case sdl.SCANCODE_B:
				e.joyPad.SetButtonPressed(bus.ButtonB, event.State == sdl.PRESSED)
			case sdl.SCANCODE_S:
				e.joyPad.SetButtonPressed(bus.Select, event.State == sdl.PRESSED)
			case sdl.SCANCODE_SPACE:
				e.joyPad.SetButtonPressed(bus.Start, event.State == sdl.PRESSED)
			default:
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

func logInstruction(p *cpu.Processor, instruction *cpu.Instruction) {
	trace.Trace(p, instruction)
}
