package biz

import (
	"context"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisjay/nesgo/backend/app/gaming/pkg/codec"
)

type RoomSession struct {
	roomId          int64
	e               *emulator.Emulator
	emulatorCancel  context.CancelFunc
	game            string
	videoEncoder    codec.IVideoEncoder
	audioSampleRate int
	audioEncoder    codec.IAudioEncoder
	audioSampleChan chan float32 // audioSampleChan 音频输出channel
	controller1     int64        // controller1 模拟器P1控制权玩家
	controller2     int64        // controller2 模拟器P2控制权玩家
}

type RoomSessionUseCase struct {
	sessions map[int64]*RoomSession
}
