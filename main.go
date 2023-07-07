package main

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/cpu"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/veandco/go-sdl2/sdl"
	"time"
	"unsafe"
)

// test_program1 死循环，不断切换(0,0)像素点的颜色
var test_program1 = []byte{
	cpu.LDX_IM, 0x1, cpu.STX_ABS, 0x00, 0x02,
	cpu.LDX_IM, 0x0, cpu.STX_ABS, 0x00, 0x02,
	cpu.JMP_A, 0x0, 0x06,
}

// test_program2, 白色像素点，从(0,0)移动到(15,0)
var test_program2 = []byte{
	cpu.LDA_IM, 0x1, cpu.LDX_IM, 0x0,
	cpu.STA_ABS, 0x0, 0x02,
	cpu.STX_ABS, 0x0, 0x02, cpu.STA_ABS, 0x1, 0x2,
	cpu.STX_ABS, 0x1, 0x02, cpu.STA_ABS, 0x2, 0x2,
	cpu.STX_ABS, 0x2, 0x02, cpu.STA_ABS, 0x3, 0x2,
	cpu.STX_ABS, 0x3, 0x02, cpu.STA_ABS, 0x4, 0x2,
	cpu.STX_ABS, 0x4, 0x02, cpu.STA_ABS, 0x5, 0x2,
	cpu.STX_ABS, 0x5, 0x02, cpu.STA_ABS, 0x6, 0x2,
	cpu.STX_ABS, 0x6, 0x02, cpu.STA_ABS, 0x7, 0x2,
	cpu.STX_ABS, 0x7, 0x02, cpu.STA_ABS, 0x8, 0x2,
	cpu.STX_ABS, 0x8, 0x02, cpu.STA_ABS, 0x9, 0x2,
	cpu.STX_ABS, 0x9, 0x02, cpu.STA_ABS, 0xa, 0x2,
	cpu.STX_ABS, 0xa, 0x02, cpu.STA_ABS, 0xb, 0x2,
	cpu.STX_ABS, 0xb, 0x02, cpu.STA_ABS, 0xc, 0x2,
	cpu.STX_ABS, 0xc, 0x02, cpu.STA_ABS, 0xd, 0x2,
	cpu.STX_ABS, 0xd, 0x02, cpu.STA_ABS, 0xe, 0x2,
	cpu.STX_ABS, 0xe, 0x02, cpu.STA_ABS, 0xf, 0x2,
}

var frame = make([]byte, 32*32*3) // 记录屏幕32x32个像素的RGBA颜色
var window *sdl.Window
var renderer *sdl.Renderer

const RenderInterval time.Duration = 100

// initSDL 初始化window和renderer
func initSDL() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("init sdl error %w", err)
	}
	w, err := sdl.CreateWindow("nesgo", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 32*20, 32*20, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("sdl create window error %w", err)
	}
	window = w
	r, err := sdl.CreateRenderer(window, 0, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return fmt.Errorf("sdl get renderer error %w", err)
	}
	_ = r.SetScale(20, 20)
	renderer = r
	return nil
}

func main() {
	processor := cpu.NewProcessor()
	if err := initSDL(); err != nil {
		panic(err)
	}
	// 用texture表示整个32x32屏幕
	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, 32, 32)
	defer texture.Destroy()
	defer window.Destroy()

	// 运行programe，callback进行屏幕渲染
	processor.LoadAndRunWithCallback(test_program1, func(p *cpu.Processor) bool {
		if handleEvents() {
			return false
		}
		// 从内存读取屏幕数据，如果发生更新就刷新屏幕像素
		updated := ppu.ReadAndUpdateScreen(p.GetMemoryRange(cpu.OutputBaseAddr, cpu.OutputEndAddr), frame)
		if updated {
			_ = texture.Update(nil, unsafe.Pointer(&frame[0]), 32*3)
			_ = renderer.Copy(texture, nil, nil)
			renderer.Present()
			time.Sleep(RenderInterval * time.Millisecond)
		}
		return true
	})
}

func handleEvents() (shutdown bool) {
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch e.(type) {
		case *sdl.QuitEvent:
			shutdown = true
		default:

		}
	}
	return
}
