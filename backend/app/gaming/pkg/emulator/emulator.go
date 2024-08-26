package emulator

import (
	"context"
	"errors"
)

type IGameFileRepo interface {
	// GetGameData 获取游戏文件数据
	GetGameData(ctx context.Context, game string) ([]byte, error)
}

// IEmulator 模拟器接口层，用于适配更多种类的模拟器
type IEmulator interface {
	// Start 启动模拟器
	Start() error
	// Pause 暂停模拟器
	Pause() error
	// Resume 继续运行模拟器
	Resume() error
	// Save 保存游戏
	Save() (IEmulatorSave, error)
	// LoadSave 加载存档，如果需要切换游戏，使用gameFileRepo加载游戏文件
	LoadSave(save IEmulatorSave, gameFileRepo IGameFileRepo) error
	// Restart 重启模拟器
	Restart(options IEmulatorOptions) error
	// Stop 结束模拟器
	Stop() error
}

// IEmulatorOptions 模拟器选项接口
type IEmulatorOptions interface {
	Game() string
	GameData() []byte
	AudioSampleRate() int
	AudioSampleChan() chan float32
}

type IEmulatorSave interface {
	GameName() string
	SaveData() []byte
}

type BaseEmulatorSave struct {
	game string
	data []byte
}

var (
	ErrorEmulatorNotSupported = errors.New("nes not supported")
)

func MakeEmulator(emulatorType string, options IEmulatorOptions) (IEmulator, error) {
	switch emulatorType {
	case "nes":
		return makeNESEmulatorAdapter(options)
	default:
		return nil, ErrorEmulatorNotSupported
	}
}

func makeNESEmulatorAdapter(options IEmulatorOptions) (IEmulator, error) {
	e, err := makeNESEmulator(options)
	if err != nil {
		return nil, err
	}
	return &NesEmulatorAdapter{
		e:       e,
		options: options,
	}, nil
}

func (s *BaseEmulatorSave) GameName() string {
	return s.game
}

func (s *BaseEmulatorSave) SaveData() []byte {
	return s.data
}
