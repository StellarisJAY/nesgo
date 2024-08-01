package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/pkg/codec"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisJAY/nesgo/emulator/config"
	"github.com/stellarisJAY/nesgo/emulator/ppu"
)

type GameInstance struct {
	RoomId          int64
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

type GameSave struct {
	Id     int64  `json:"id"`
	RoomId int64  `json:"roomId"`
	Game   string `json:"game"`
	Data   string `json:"data"`
}

type GameInstanceRepo interface {
	CreateGameInstance(ctx context.Context, game *GameInstance) error
	DeleteGameInstance(ctx context.Context, roomId int64) error
	GetGameInstance(ctx context.Context, roomId int64) (*GameInstance, error)
}

type GameFileRepo interface {
	GetGameData(ctx context.Context, game string) ([]byte, error)
	GetSavedGame(ctx context.Context, id int64) (*GameSave, error)
	SaveGame(ctx context.Context, save *GameSave) error
}

type GameInstanceUseCase struct {
	repo         GameInstanceRepo
	gameFileRepo GameFileRepo
	logger       *log.Helper
}

func NewGameInstanceUseCase(repo GameInstanceRepo, gameFileRepo GameFileRepo, logger log.Logger) *GameInstanceUseCase {
	return &GameInstanceUseCase{
		repo:         repo,
		gameFileRepo: gameFileRepo,
		logger:       log.NewHelper(log.With(logger, "module", "biz/gameInstance")),
	}
}

func (uc *GameInstanceUseCase) CreateGameInstance(ctx context.Context, roomId int64, game string) (*GameInstance, error) {
	gameData, err := uc.gameFileRepo.GetGameData(ctx, game)
	if err != nil {
		return nil, err
	}

	instance := GameInstance{
		RoomId: roomId,
		game:   game,
	}
	emulatorConfig := config.Config{
		Game:        game,
		EnableTrace: false,
		Disassemble: false,
		MuteApu:     false,
		Debug:       false,
	}
	instance.audioSampleRate = 44100
	instance.audioSampleChan = make(chan float32, instance.audioSampleRate)
	renderCallback := func(ppu *ppu.PPU) {
		instance.RenderCallback(ppu, uc.logger)
	}
	e, err := emulator.NewEmulatorWithGameData(gameData, emulatorConfig, renderCallback, instance.audioSampleChan, instance.audioSampleRate)
	if err != nil {
		return nil, err
	}
	instance.e = e
	videoEncoder, err := codec.NewVideoEncoder("h264")
	if err != nil {
		return nil, err
	}
	audioEncoder, err := codec.NewAudioEncoder(instance.audioSampleRate)
	if err != nil {
		return nil, err
	}
	instance.videoEncoder = videoEncoder
	instance.audioEncoder = audioEncoder

	emulatorCtx, cancel := context.WithCancel(context.Background())
	instance.emulatorCancel = cancel
	go instance.e.LoadAndRun(emulatorCtx, false)

	_ = uc.repo.CreateGameInstance(ctx, &instance)
	return &instance, nil
}
