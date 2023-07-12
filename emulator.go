package main

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/cpu"
	"github.com/stellarisJAY/nesgo/mem"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/veandco/go-sdl2/sdl"
	"time"
	"unsafe"
)

const RenderInterval = 20000 * time.Nanosecond

type Emulator struct {
	processor *cpu.Processor
	rom       *mem.ROM
	window    *sdl.Window
	frame     []byte
	renderer  *sdl.Renderer
}

func NewEmulator(nesData []byte) *Emulator {
	rom := mem.NewROM(nesData)
	processor := cpu.NewProcessorWithROM(rom)
	window, renderer, err := initSDL()
	if err != nil {
		panic(err)
	}
	frame := make([]byte, 32*32*3)
	return &Emulator{
		processor: &processor,
		rom:       rom,
		window:    window,
		frame:     frame,
		renderer:  renderer,
	}
}

func initSDL() (*sdl.Window, *sdl.Renderer, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, nil, fmt.Errorf("init sdl error %w", err)
	}
	w, err := sdl.CreateWindow("nesgo", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 32*20, 32*20, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, nil, fmt.Errorf("sdl create window error %w", err)
	}
	r, err := sdl.CreateRenderer(w, 0, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, nil, fmt.Errorf("sdl get renderer error %w", err)
	}
	_ = r.SetScale(20, 20)
	return w, r, nil
}

func (e *Emulator) LoadAndRun() {
	renderer, frame, processor := e.renderer, e.frame, e.processor
	// 用texture表示整个32x32屏幕
	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, 32, 32)
	defer texture.Destroy()
	defer e.window.Destroy()

	// 运行program，callback进行屏幕渲染
	processor.LoadAndRunWithCallback(handleEvents, func(p *cpu.Processor) bool {
		// 从内存读取屏幕数据，如果发生更新就刷新屏幕像素
		updated := ppu.ReadAndUpdateScreen(p.GetMemoryRange(cpu.OutputBaseAddr, cpu.OutputEndAddr), frame)
		if updated {
			_ = texture.Update(nil, unsafe.Pointer(&frame[0]), 32*3)
			_ = e.renderer.Copy(texture, nil, nil)
			renderer.Present()
		}
		time.Sleep(RenderInterval)
		return true
	})
}

func handleEvents(p *cpu.Processor) bool {
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch e.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			p.HandleKeyboardEvent(e.(*sdl.KeyboardEvent))
		default:
		}
	}
	return true
}
