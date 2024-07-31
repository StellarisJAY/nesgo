package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/emulator"
	v1 "github.com/stellarisjay/nesgo/backend/api/gaming/service/v1"
	"github.com/stellarisjay/nesgo/backend/app/gaming/pkg/codec"
)

type RoomGamingSession struct {
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

type RoomSession struct {
	RoomId   int64  `json:"roomId"`
	Endpoint string `json:"endpoint"`
}

type RoomSessionUseCase struct {
	sessions map[int64]*RoomGamingSession
	repo     RoomSessionRepo
	logger   *log.Helper
}

type RoomSessionRepo interface {
	CreateRoomSession(ctx context.Context, session *RoomSession) error
	GetRoomSession(ctx context.Context, roomId int64) (*RoomSession, error)
	DeleteRoomSession(ctx context.Context, roomId int64) error
}

func NewRoomSessionUseCase(roomSessionRepo RoomSessionRepo, logger log.Logger) *RoomSessionUseCase {
	return &RoomSessionUseCase{
		sessions: make(map[int64]*RoomGamingSession),
		repo:     roomSessionRepo,
		logger:   log.NewHelper(log.With(logger, "module", "biz/room_session")),
	}
}

func (uc *RoomSessionUseCase) GetOrCreateRoomSession(ctx context.Context, roomId int64) (*RoomSession, error) {
	session, err := uc.repo.GetRoomSession(ctx, roomId)
	if err != nil {
		return nil, v1.ErrorGetRoomSessionFailed("repo error: %v", err)
	}
	if session != nil {
		return session, nil
	}
	session.RoomId = roomId
	err = uc.repo.CreateRoomSession(ctx, session)
	if err != nil {
		return nil, v1.ErrorGetRoomSessionFailed("create room session repo error: %v", err)
	}
	return session, nil
}
