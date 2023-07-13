package main

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/stellarisJAY/nesgo/ppu"
	"github.com/veandco/go-sdl2/sdl"
	"io"
	"os"
	"testing"
	"unsafe"
)

func loadNES(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open game file %w", err)
	}
	program, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read game file error %w", err)
	}
	return program, nil
}

func TestPPU_RenderTile(t *testing.T) {
	program, err := loadNES("games/PacMan.nes")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	rom := bus.NewROM(program)
	p := ppu.NewPPU(rom.GetChrROM(), byte(rom.GetMirroring()))
	p.DisplayAllTiles(1)
	window, renderer, err := initSDL()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, ppu.WIDTH, ppu.HEIGHT)
	defer texture.Destroy()
	defer window.Destroy()
	frame := p.FrameData()
	_ = texture.Update(nil, unsafe.Pointer(&frame[0]), ppu.WIDTH*3)
	_ = renderer.Copy(texture, nil, nil)
	renderer.Present()

	select {}
}
