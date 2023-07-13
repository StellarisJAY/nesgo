package main

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/cpu"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"time"
	"unsafe"
)

const RenderInterval = 20000 * time.Nanosecond
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
	w, err := sdl.CreateWindow("nesgo", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, ppu.WIDTH*SCALE, ppu.HEIGHT*SCALE, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, nil, fmt.Errorf("sdl create window error %w", err)
	}
	r, err := sdl.CreateRenderer(w, 0, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, nil, fmt.Errorf("sdl get renderer error %w", err)
	}
	_ = r.SetScale(SCALE, SCALE)
	return w, r, nil
}

func (e *Emulator) LoadAndRun() {
	defer e.texture.Destroy()
	defer e.window.Destroy()
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			select {}
		}
	}()
	// 运行program，callback进行屏幕渲染
	e.processor.LoadAndRunWithCallback(e.handleEvents, logInstruction)
}

func (e *Emulator) RendererCallback(p *ppu.PPU) {
	p.Render()
	frame := p.FrameData()
	_ = e.texture.Update(nil, unsafe.Pointer(&frame[0]), ppu.WIDTH*3)
	_ = e.renderer.Copy(e.texture, nil, nil)
	e.renderer.Present()
}

func (e *Emulator) handleEvents(p *cpu.Processor) bool {
	for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
		switch ev.(type) {
		case *sdl.QuitEvent:
			return false
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
	return true
}

func logInstruction(instruction cpu.Instruction, pc uint16, args []byte) {
	var addrMode string
	switch instruction.AddrMode {
	case cpu.NoneAddressing:
		addrMode = "none"
	case cpu.Immediate:
		addrMode = "im"
	case cpu.ZeroPage:
		addrMode = "zp"
	case cpu.ZeroPageX:
		addrMode = "zpx"
	case cpu.ZeroPageY:
		addrMode = "zpy"
	case cpu.Absolute:
		addrMode = "abs"
	case cpu.AbsoluteX:
		addrMode = "abx"
	case cpu.AbsoluteY:
		addrMode = "aby"
	case cpu.IndirectX:
		addrMode = "inx"
	case cpu.IndirectY:
		addrMode = "iny"
	}
	log.Printf("0x%x : %s-%s : [%d, %d], [%d, %d]\n", pc, instruction.Name, addrMode, args[0], args[1], int8(args[0]), int8(args[1]))
}
