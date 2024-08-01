package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/emulator/ppu"
)

func (g *GameInstance) RenderCallback(ppu *ppu.PPU, logger *log.Helper) {
	frame := ppu.Frame()
	_, release, err := g.videoEncoder.Encode(frame)
	if err != nil {
		logger.Error("encode frame error", "err", err)
		return
	}
	defer release()
	// TODO broadcast video frame
	return
}
