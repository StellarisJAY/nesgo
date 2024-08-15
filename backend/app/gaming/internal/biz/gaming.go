package biz

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
	v1 "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/pkg/codec"
	"github.com/stellarisJAY/nesgo/emulator"
	"github.com/stellarisJAY/nesgo/emulator/config"
	"github.com/stellarisJAY/nesgo/emulator/ppu"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"runtime"
	"sync"
	"time"
)

type GameSave struct {
	Id     int64  `json:"id"`
	RoomId int64  `json:"roomId"`
	Game   string `json:"game"`
	Data   string `json:"data"`
}

type GameFileMetadata struct {
	Name       string `json:"name"`
	Mapper     string `json:"mapper"`
	Mirroring  string `json:"mirroring"`
	Size       int32  `json:"size"`
	UploadTime int64  `json:"uploadTime"`
}

type GameInstanceStats struct {
	RoomId            int64         `json:"roomId"`
	Connections       int           `json:"connections"`
	ActiveConnections int           `json:"activeConnections"`
	Game              string        `json:"game"`
	Uptime            time.Duration `json:"uptime"`
}

type EndpointStats struct {
	EmulatorCount int32 `json:"emulatorCount"`
	CpuUsage      int32 `json:"cpuUsage"`
	MemoryUsed    int64 `json:"memoryUsed"`
	MemoryTotal   int64 `json:"memoryTotal"`
	Uptime        int64 `json:"uptime"`
}

type GameInstanceRepo interface {
	CreateGameInstance(ctx context.Context, game *GameInstance) (int64, error)
	DeleteGameInstance(ctx context.Context, roomId int64) error
	GetGameInstance(ctx context.Context, roomId int64) (*GameInstance, error)
	ListGameInstances(ctx context.Context) ([]*GameInstance, error)
}

type GameFileRepo interface {
	GetGameData(ctx context.Context, game string) ([]byte, error)
	UploadGameData(ctx context.Context, game string, data []byte, metadata *GameFileMetadata) error
	ListGames(ctx context.Context, page, pageSize int) ([]*GameFileMetadata, int, error)
	DeleteGames(ctx context.Context, games []string) (int, error)
	GetSavedGame(ctx context.Context, id int64) (*GameSave, error)
	SaveGame(ctx context.Context, save *GameSave) error
}

type GameInstanceUseCase struct {
	repo         GameInstanceRepo
	gameFileRepo GameFileRepo
	logger       *log.Helper
	startupTime  time.Time
}

func NewGameInstanceUseCase(repo GameInstanceRepo, gameFileRepo GameFileRepo, logger log.Logger) *GameInstanceUseCase {
	return &GameInstanceUseCase{
		repo:         repo,
		gameFileRepo: gameFileRepo,
		logger:       log.NewHelper(log.With(logger, "module", "biz/gameInstance")),
		startupTime:  time.Now(),
	}
}

// CreateGameInstance 创建模拟器实例，第一个连接房间并创建会话的操作会创建模拟器实例
// 调用者必须持有房间会话的分布式锁，保证只创建一次
func (uc *GameInstanceUseCase) CreateGameInstance(ctx context.Context, roomId int64, game string) (*GameInstance, error) {
	gameData, err := uc.gameFileRepo.GetGameData(ctx, game)
	if errors.Is(err, gridfs.ErrFileNotFound) {
		return nil, v1.ErrorGameFileNotFound("game file not found")
	}
	if err != nil {
		return nil, v1.ErrorCreateGameInstanceFailed("gem game file failed: %v", err)
	}
	instance := GameInstance{
		RoomId:      roomId,
		game:        game,
		messageChan: make(chan *Message, 64),
		mutex:       &sync.RWMutex{},
		connections: make(map[int64]*Connection),
		createTime:  time.Now(),
	}
	emulatorConfig := config.Config{
		Game:        game,
		EnableTrace: false,
		Disassemble: false,
		MuteApu:     false,
		Debug:       false,
	}
	instance.audioSampleRate = 48000
	instance.audioSampleChan = make(chan float32, instance.audioSampleRate)
	// create emulator
	renderCallback := func(ppu *ppu.PPU) {
		instance.RenderCallback(ppu, uc.logger)
	}
	e, err := emulator.NewEmulatorWithGameData(gameData, emulatorConfig, renderCallback, instance.audioSampleChan, instance.audioSampleRate)
	if err != nil {
		return nil, v1.ErrorCreateGameInstanceFailed("create emulator failed: %v", err)
	}
	instance.e = e
	// create video and audio encoder
	videoEncoder, err := codec.NewVideoEncoder("h264")
	if err != nil {
		return nil, v1.ErrorCreateGameInstanceFailed("create video encoder failed: %v", err)
	}
	audioEncoder, err := codec.NewAudioEncoder(instance.audioSampleRate)
	if err != nil {
		return nil, v1.ErrorGameFileNotFound("create audio encoder failed: %v", err)
	}
	instance.videoEncoder = videoEncoder
	instance.audioEncoder = audioEncoder

	// start emulator
	emulatorCtx, cancel := context.WithCancel(context.Background())
	instance.emulatorCancel = cancel
	uc.logger.Infof("emulator created, roomId: %d", roomId)
	go instance.e.LoadAndRun(emulatorCtx, false)
	instance.e.Pause()
	// collect audio samples
	go instance.audioSampleListener(emulatorCtx, uc.logger)

	// start message consumer
	msgConsumerCtx, msgConsumerCancel := context.WithCancel(context.Background())
	go instance.messageConsumer(msgConsumerCtx)
	instance.cancel = msgConsumerCancel

	leaseID, _ := uc.repo.CreateGameInstance(ctx, &instance)
	instance.LeaseID = leaseID
	return &instance, nil
}

func (uc *GameInstanceUseCase) OpenGameConnection(ctx context.Context, roomId, userId int64) (string, error) {
	instance, err := uc.repo.GetGameInstance(ctx, roomId)
	if err != nil {
		return "", err
	}
	if instance == nil {
		return "", v1.ErrorUnknownError("game instance not found")
	}
	_, sdp, err := instance.NewConnection(userId)
	if err != nil {
		return "", v1.ErrorOpenGameConnectionFailed("open game connection failed: %v", err)
	}
	return sdp, nil
}

func (uc *GameInstanceUseCase) SDPAnswer(ctx context.Context, roomId, userId int64, sdpAnswer string) error {
	instance, err := uc.repo.GetGameInstance(ctx, roomId)
	if err != nil {
		return err
	}
	if instance == nil {
		return v1.ErrorUnknownError("game instance not found")
	}
	instance.mutex.RLock()
	conn, ok := instance.connections[userId]
	instance.mutex.RUnlock()
	if !ok {
		return v1.ErrorGameConnectionNotFound("game connection not found")
	}

	sdp := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  sdpAnswer,
	}
	err = conn.pc.SetRemoteDescription(sdp)
	if err != nil {
		return v1.ErrorSdpAnswerFailed("set remote sdp failed: %v", err)
	}
	return nil
}

func (uc *GameInstanceUseCase) ICECandidate(ctx context.Context, roomId, userId int64, candidate string) error {
	instance, err := uc.repo.GetGameInstance(ctx, roomId)
	if err != nil {
		return err
	}
	if instance == nil {
		return v1.ErrorUnknownError("game instance not found")
	}
	instance.mutex.RLock()
	conn, ok := instance.connections[userId]
	instance.mutex.RUnlock()
	if !ok {
		return v1.ErrorGameConnectionNotFound("game connection not found")
	}
	candidateInit := webrtc.ICECandidateInit{}
	err = json.Unmarshal([]byte(candidate), &candidateInit)
	if err != nil {
		return v1.ErrorIceCandidateFailed("unmarshal ice candidate failed: %v", err)
	}
	err = conn.pc.AddICECandidate(candidateInit)
	if err != nil {
		return v1.ErrorIceCandidateFailed("add ice candidate failed: %v", err)
	}
	return nil
}

// ReleaseGameInstance 释放模拟器实例，所有连接断开后延迟释放
// 调用者必须持有房间会话的分布式锁，保证没有新连接建立
func (uc *GameInstanceUseCase) ReleaseGameInstance(ctx context.Context, roomId int64, force bool) error {
	instance, _ := uc.repo.GetGameInstance(ctx, roomId)
	if instance == nil {
		return v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if instance.status.Load() == InstanceStatusStopped {
		return nil
	}
	if len(instance.connections) > 0 && !force {
		return v1.ErrorGameInstanceNotAccessible("can't release active game instance")
	}
	if len(instance.connections) > 0 && force {
		for _, conn := range instance.connections {
			conn.Close()
		}
	}
	instance.emulatorCancel()
	instance.cancel()
	instance.status.Store(InstanceStatusStopped)
	_ = uc.repo.DeleteGameInstance(ctx, roomId)
	return nil
}

func (uc *GameInstanceUseCase) SetController(ctx context.Context, roomId, playerId int64, controllerId int) error {
	instance, err := uc.repo.GetGameInstance(ctx, roomId)
	if err != nil {
		return err
	}
	if instance == nil {
		return v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	msgType := MsgResetController
	if controllerId == 1 {
		msgType = MsgSetController1
	}
	if controllerId == 2 {
		msgType = MsgSetController2
	}
	resultChan := make(chan ConsumerResult)
	instance.messageChan <- &Message{
		Data:       playerId,
		Type:       msgType,
		resultChan: resultChan,
	}
	res := <-resultChan
	if !res.Success {
		return v1.ErrorOperationFailed("set controller failed")
	}
	return nil
}

func (uc *GameInstanceUseCase) GetController(ctx context.Context, roomId int64) (int64, int64, error) {
	instance, err := uc.repo.GetGameInstance(ctx, roomId)
	if err != nil {
		return 0, 0, err
	}
	if instance == nil {
		return 0, 0, v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	return instance.controller1, instance.controller2, nil
}

func (uc *GameInstanceUseCase) GetGameInstanceStats(ctx context.Context, roomId int64) (*GameInstanceStats, error) {
	instance, err := uc.repo.GetGameInstance(ctx, roomId)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	return instance.DumpStats(), nil
}

func (uc *GameInstanceUseCase) ListGameInstances(ctx context.Context) ([]*GameInstanceStats, error) {
	instances, _ := uc.repo.ListGameInstances(ctx)
	result := make([]*GameInstanceStats, 0, len(instances))
	for _, instance := range instances {
		result = append(result, instance.DumpStats())
	}
	return result, nil
}

func (uc *GameInstanceUseCase) DeleteMemberConnection(ctx context.Context, roomId, userId int64) error {
	instance, _ := uc.repo.GetGameInstance(ctx, roomId)
	if instance == nil {
		return nil
	}
	instance.DeleteConnection(userId)
	return nil
}

func (uc *GameInstanceUseCase) GetEndpointStats(ctx context.Context) (*EndpointStats, error) {
	instances, _ := uc.repo.ListGameInstances(ctx)
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)
	// TODO collect OS stats
	return &EndpointStats{
		EmulatorCount: int32(len(instances)),
		CpuUsage:      0,
		MemoryUsed:    int64(memStats.Alloc),
		MemoryTotal:   int64(memStats.Sys),
		Uptime:        time.Now().Sub(uc.startupTime).Milliseconds(),
	}, nil
}
