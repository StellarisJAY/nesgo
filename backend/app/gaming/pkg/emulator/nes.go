package emulator

import (
	"context"

	"github.com/stellarisJAY/nesgo/nes"
	nesConfig "github.com/stellarisJAY/nesgo/nes/config"
	"github.com/stellarisJAY/nesgo/nes/ppu"
)

type NesEmulatorAdapter struct {
	e          *nes.Emulator
	options    IEmulatorOptions
	cancelFunc context.CancelFunc
}

type NesEmulatorOptions struct {
	NesGame               string
	NesGameData           []byte
	OutputAudioSampleRate int
	OutputAudioSampleChan chan float32
	NesRenderCallback     func(frame *ppu.Frame)
}

func (n *NesEmulatorOptions) GameData() []byte {
	return n.NesGameData
}

func (n *NesEmulatorOptions) AudioSampleRate() int {
	return n.OutputAudioSampleRate
}

func (n *NesEmulatorOptions) AudioSampleChan() chan float32 {
	return n.OutputAudioSampleChan
}

func (n *NesEmulatorOptions) Game() string {
	return n.NesGame
}

// Start 启动NES模拟器，创建单独的goroutine运行CPU循环，使用context打断
func (n *NesEmulatorAdapter) Start() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go n.e.LoadAndRun(ctx, false)
	n.cancelFunc = cancelFunc
	return nil
}

// Pause 暂停NES模拟器
func (n *NesEmulatorAdapter) Pause() error {
	n.e.Pause()
	return nil
}

// Resume 恢复NES模拟器
func (n *NesEmulatorAdapter) Resume() error {
	n.e.Resume()
	return nil
}

// Restart 重启NES模拟器，结束旧模拟器goroutine，创建并运行新模拟器
func (n *NesEmulatorAdapter) Restart(options IEmulatorOptions) error {
	n.cancelFunc()
	e, err := makeNESEmulator(options)
	if err != nil {
		return err
	}
	n.e = e
	ctx, cancelFunc := context.WithCancel(context.Background())
	go n.e.LoadAndRun(ctx, false)
	n.cancelFunc = cancelFunc
	n.options = options
	return nil
}

// Stop 关闭NES模拟器
func (n *NesEmulatorAdapter) Stop() error {
	n.cancelFunc()
	return nil
}

// Save 获取模拟器存档数据
func (n *NesEmulatorAdapter) Save() (IEmulatorSave, error) {
	n.e.Pause()
	defer n.e.Resume()
	s, err := n.e.GetSaveData()
	if err != nil {
		return nil, err
	}
	return &BaseEmulatorSave{
		data: s,
	}, nil
}

func (n *NesEmulatorAdapter) LoadSave(save IEmulatorSave, gameFileRepo IGameFileRepo) error {
	n.e.Pause()
	defer n.e.Resume()
	// 如果游戏名称不同，需要加载新的游戏文件，并重启模拟器
	if n.options.Game() != save.GameName() {
		gameData, err := gameFileRepo.GetGameData(context.Background(), save.GameName())
		if err != nil {
			return err
		}
		oldOpts := n.options.(*NesEmulatorOptions)
		err = n.Restart(&NesEmulatorOptions{
			NesGame:               save.GameName(),
			NesGameData:           gameData,
			OutputAudioSampleRate: oldOpts.OutputAudioSampleRate,
			OutputAudioSampleChan: oldOpts.OutputAudioSampleChan,
			NesRenderCallback:     oldOpts.NesRenderCallback,
		})
		if err != nil {
			return err
		}
		return n.e.Load(save.SaveData())
	} else {
		return n.e.Load(save.SaveData())
	}
}

func makeNESEmulator(options IEmulatorOptions) (*nes.Emulator, error) {
	configs := nesConfig.Config{
		Game:               options.Game(),
		EnableTrace:        false,
		Disassemble:        false,
		SnapshotSerializer: "json",
		MuteApu:            false,
		Debug:              false,
	}
	renderCallback := func(p *ppu.PPU) {
		options.(*NesEmulatorOptions).NesRenderCallback(p.Frame())
	}
	e, err := nes.NewEmulatorWithGameData(options.GameData(), configs, renderCallback, options.AudioSampleChan(), options.AudioSampleRate())
	if err != nil {
		return nil, err
	}
	return e, nil
}
