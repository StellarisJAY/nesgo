package emulator

import (
	"context"
	"errors"
	"image"
)

const (
	EmulatorTypeNES = "NES"
)

// IFrame 模拟器输出画面接口
type IFrame interface {
	// Width 画面宽度
	Width() int
	// Height 画面高度
	Height() int
	// YCbCr 画面YUV格式数据
	YCbCr() *image.YCbCr
	// Read 读取画面数据， func() 用于释放资源
	Read() (image.Image, func(), error)
}

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
	// SubmitInput 提交模拟器按键输入
	SubmitInput(controllerId int, keyCode string, pressed bool)
	// SetGraphicOptions 修改模拟器画面设置
	SetGraphicOptions(*GraphicOptions)

	GetCPUBoostRate() float64
	SetCPUBoostRate(float64) float64
	// OutputResolution 获取模拟器输出分辨率
	OutputResolution() (width, height int)
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

type GraphicOptions struct {
	ReverseColor bool
	Grayscale    bool
}

type BaseEmulatorSave struct {
	Game string
	Data []byte
}

type BaseFrame struct {
	image  *image.YCbCr
	width  int
	height int
}

var (
	ErrorEmulatorNotSupported = errors.New("emulator not supported")
)

func MakeEmulator(emulatorType string, options IEmulatorOptions) (IEmulator, error) {
	switch emulatorType {
	case EmulatorTypeNES:
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
	return s.Game
}

func (s *BaseEmulatorSave) SaveData() []byte {
	return s.Data
}

func MakeBaseFrame(image *image.YCbCr, width, height int) IFrame {
	return &BaseFrame{
		image:  image,
		width:  width,
		height: height,
	}
}

func (b *BaseFrame) Width() int {
	return b.width
}

func (b *BaseFrame) Height() int {
	return b.height
}

func (b *BaseFrame) YCbCr() *image.YCbCr {
	return b.image
}

func (b *BaseFrame) Read() (image.Image, func(), error) {
	return b.image, func() {}, nil
}
