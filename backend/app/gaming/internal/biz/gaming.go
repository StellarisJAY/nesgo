package biz

import (
	"context"
	"encoding/json"
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
	v1 "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/conf"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/pkg/codec"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/pkg/emulator"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

const (
	DefaultAudioSampleRate = 48000
)

type GameSave struct {
	Id         int64  `json:"id"`
	RoomId     int64  `json:"roomId"`
	Game       string `json:"game"`
	Data       []byte `json:"data"`
	CreateTime int64  `json:"createTime"`
	ExitSave   bool   `json:"exitSave"`
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

type GraphicOptions struct {
	HighResOpen  bool `json:"highResOpen"`
	ReverseColor bool `json:"reverseColor"`
	Grayscale    bool `json:"grayscale"`
}

type SupportedEmulator struct {
	Name                  string `json:"name"`
	SupportSaving         bool   `json:"saving"`
	SupportGraphicSetting bool   `json:"supportGraphicSetting"`
	Provider              string `json:"provider"`
	Games                 int32  `json:"games"`
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
	ListSaves(ctx context.Context, roomId int64, page, pageSize int32) ([]*GameSave, int32, error)
	DeleteSave(ctx context.Context, saveId int64) error
	GetExitSave(ctx context.Context, roomId int64) (*GameSave, error)
}

type RoomRepo interface {
	AddDeleteRoomSessionTask(ctx context.Context, roomId int64, instanceId string) error
}

type GameInstanceUseCase struct {
	repo         GameInstanceRepo
	gameFileRepo GameFileRepo
	roomRepo     RoomRepo
	logger       *log.Helper
	startupTime  time.Time
	stunServer   string
}

func NewGameInstanceUseCase(repo GameInstanceRepo, gameFileRepo GameFileRepo, roomRepo RoomRepo, c *conf.IceServer, logger log.Logger) *GameInstanceUseCase {
	return &GameInstanceUseCase{
		repo:         repo,
		gameFileRepo: gameFileRepo,
		roomRepo:     roomRepo,
		logger:       log.NewHelper(log.With(logger, "module", "biz/gameInstance")),
		startupTime:  time.Now(),
		stunServer:   c.StunServer,
	}
}

// CreateGameInstance 创建模拟器实例，第一个连接房间并创建会话的操作会创建模拟器实例
// 调用者必须持有房间会话的分布式锁，保证只创建一次
func (uc *GameInstanceUseCase) CreateGameInstance(ctx context.Context, roomId int64, game string, emulatorType string) (*GameInstance, error) {
	gameData, err := uc.gameFileRepo.GetGameData(ctx, game)
	if errors.Is(err, gridfs.ErrFileNotFound) {
		return nil, v1.ErrorGameFileNotFound("game file not found")
	}
	if err != nil {
		return nil, v1.ErrorCreateGameInstanceFailed("gem game file failed: %v", err)
	}

	instance := GameInstance{
		RoomId:               roomId,
		EmulatorType:         emulatorType,
		game:                 game,
		messageChan:          make(chan *Message, 64),
		mutex:                &sync.RWMutex{},
		connections:          make(map[int64]*Connection),
		createTime:           time.Now(),
		allConnCloseCallback: uc.OnGameInstanceConnectionsAllClosed,
		audioSampleRate:      DefaultAudioSampleRate,
		audioSampleChan:      make(chan float32, DefaultAudioSampleRate),
		frameEnhancer: func(frame emulator.IFrame) emulator.IFrame {
			return frame
		},
	}
	// TODO select emulator
	opts, err := instance.makeEmulatorOptions(emulatorType, game, gameData)
	if err != nil {
		return nil, v1.ErrorCreateGameInstanceFailed("make emulator options failed: %v", err)
	}
	// TODO select emulator
	e, err := emulator.MakeEmulator(emulatorType, opts)
	if err != nil {
		return nil, v1.ErrorCreateGameInstanceFailed("create emulator failed: %v", err)
	}
	instance.e = e
	// create video and audio encoder
	width, height := e.OutputResolution()
	videoEncoder, err := codec.NewVideoEncoder("vp8", width, height)
	if err != nil {
		return nil, v1.ErrorCreateGameInstanceFailed("create video encoder failed: %v", err)
	}
	audioEncoder, err := codec.NewAudioEncoder(instance.audioSampleRate)
	if err != nil {
		return nil, v1.ErrorGameFileNotFound("create audio encoder failed: %v", err)
	}
	instance.videoEncoder = videoEncoder
	instance.audioEncoder = audioEncoder

	// 启动模拟器，之后暂停等待第一个连接激活
	e.Start()
	e.Pause()

	// 消息处理器和音频处理器
	msgConsumerCtx, msgConsumerCancel := context.WithCancel(context.Background())
	go instance.messageConsumer(msgConsumerCtx)
	go instance.audioSampleListener(msgConsumerCtx, log.NewHelper(log.With(log.DefaultLogger, "module", "audioSender")))
	instance.cancel = msgConsumerCancel

	// 加载退出时的存档
	save, err := uc.gameFileRepo.GetExitSave(ctx, roomId)
	if save != nil {
		err := instance.LoadSave(save.Data, save.Game, uc.gameFileRepo)
		if err != nil {
			uc.logger.Errorf("start nes load exit save error: %v", err)
		}
	} else if err != nil {
		uc.logger.Errorf("start nes load exit save error: %v", err)
	}

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
	_, sdp, err := instance.NewConnection(userId, uc.stunServer)
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
	save, err := instance.SaveGame()
	if err != nil {
		return v1.ErrorSaveGameFailed("exit save game failed: %v", err)
	}
	instance.e.Stop()
	instance.cancel()
	instance.status.Store(InstanceStatusStopped)
	save.ExitSave = true
	save.RoomId = roomId
	uc.logger.Infof("delete game instance, saving on exit, roomId=%d", instance.RoomId)
	err = uc.gameFileRepo.SaveGame(ctx, save)
	if err != nil {
		return v1.ErrorSaveGameFailed("exit save game failed: %v", err)
	}
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
		Uptime:        time.Since(uc.startupTime).Milliseconds(),
	}, nil
}

func (uc *GameInstanceUseCase) SaveGame(ctx context.Context, roomId int64) error {
	instance, _ := uc.repo.GetGameInstance(ctx, roomId)
	if instance == nil {
		return v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	save, err := instance.SaveGame()
	if err != nil {
		return v1.ErrorSaveGameFailed("nes error: %v", err)
	}
	save.RoomId = roomId
	err = uc.gameFileRepo.SaveGame(ctx, save)
	if err != nil {
		return v1.ErrorSaveGameFailed("database error: %v", err)
	}
	return nil
}

func (uc *GameInstanceUseCase) LoadSave(ctx context.Context, saveId int64) error {
	save, err := uc.gameFileRepo.GetSavedGame(ctx, saveId)
	if err != nil {
		return v1.ErrorLoadSaveFailed("database error: %v", err)
	}
	instance, _ := uc.repo.GetGameInstance(ctx, save.RoomId)
	if instance == nil {
		return v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	err = instance.LoadSave(save.Data, save.Game, uc.gameFileRepo)
	if err != nil {
		return v1.ErrorLoadSaveFailed("nes error: %v", err)
	}
	return nil
}

func (uc *GameInstanceUseCase) RestartEmulator(ctx context.Context, roomId int64, game string) error {
	instance, _ := uc.repo.GetGameInstance(ctx, roomId)
	if instance == nil {
		return v1.ErrorGameInstanceNotAccessible("game instance not found")
	}

	data, err := uc.gameFileRepo.GetGameData(ctx, game)
	if err != nil {
		return v1.ErrorRestartFailed("get game data error: %v", err)
	}
	err = instance.RestartEmulator(game, data)
	if err != nil {
		return v1.ErrorRestartFailed("restart nes error: %v", err)
	}
	return nil
}

func (uc *GameInstanceUseCase) GetICECandidates(ctx context.Context, roomId, userId int64) ([]string, error) {
	instance, _ := uc.repo.GetGameInstance(ctx, roomId)
	if instance == nil {
		return nil, v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	instance.mutex.RLock()
	conn, ok := instance.connections[userId]
	instance.mutex.RUnlock()
	if !ok {
		return nil, v1.ErrorGameConnectionNotFound("game connection not found")
	}
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	result := make([]string, 0, len(conn.localCandidates))
	for _, candidate := range conn.localCandidates {
		candidateInit := candidate.ToJSON()
		data, _ := json.Marshal(candidateInit)
		result = append(result, string(data))
	}
	conn.localCandidates = make([]*webrtc.ICECandidate, 0)
	return result, nil
}

func (uc *GameInstanceUseCase) OnGameInstanceConnectionsAllClosed(instance *GameInstance) {
	uc.logger.Infof("connection all closed in room:%d", instance.RoomId)
	err := uc.roomRepo.AddDeleteRoomSessionTask(context.TODO(), instance.RoomId, instance.InstanceId)
	if err != nil {
		uc.logger.Errorf("add delete room session task failed: %v", err)
	}
}

func (uc *GameInstanceUseCase) GetGraphicOptions(ctx context.Context, roomId int64) (*GraphicOptions, error) {
	instance, _ := uc.repo.GetGameInstance(ctx, roomId)
	if instance == nil {
		return nil, v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	return &GraphicOptions{HighResOpen: instance.enhanceFrameOpen, ReverseColor: instance.reverseColorOpen, Grayscale: instance.grayscaleOpen}, nil
}

func (uc *GameInstanceUseCase) SetGraphicOptions(ctx context.Context, roomId int64, options *GraphicOptions) error {
	instance, _ := uc.repo.GetGameInstance(ctx, roomId)
	if instance == nil {
		return v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	instance.SetGraphicOptions(options)
	return nil
}

func (uc *GameInstanceUseCase) SetEmulatorSpeed(ctx context.Context, roomId int64, rate float64) (float64, error) {
	instance, _ := uc.repo.GetGameInstance(ctx, roomId)
	if instance == nil {
		return 0.0, v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	return instance.SetEmulatorSpeed(rate), nil
}

func (uc *GameInstanceUseCase) GetEmulatorSpeed(ctx context.Context, roomId int64) (float64, error) {
	instance, _ := uc.repo.GetGameInstance(ctx, roomId)
	if instance == nil {
		return 0.0, v1.ErrorGameInstanceNotAccessible("game instance not found")
	}
	return instance.e.GetCPUBoostRate(), nil
}

func (uc *GameInstanceUseCase) ListSupportedEmulators(ctx context.Context) ([]*SupportedEmulator, error) {
	emulators := []*SupportedEmulator{
		{
			Name:                  "NES",
			Provider:              "https://github.com/stellarisjay/nesgo",
			SupportSaving:         true,
			SupportGraphicSetting: true,
			Games:                 0,
		},
	}
	// TODO count emulator suppported games
	return emulators, nil
}
