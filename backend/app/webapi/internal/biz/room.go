package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type RoomUseCase struct {
	repo   RoomRepo
	logger *log.Helper
}

type RoomSession struct {
	RoomId   int64  `json:"roomId"`
	Endpoint string `json:"endpoint"`
}

type RoomRepo interface {
	GetRoomSession(ctx context.Context, roomId, userId int64) (*RoomSession, error)
}

func NewRoomUseCase(repo RoomRepo, logger log.Logger) *RoomUseCase {
	return &RoomUseCase{
		repo:   repo,
		logger: log.NewHelper(log.With(logger, "module", "biz/room")),
	}
}

func (uc *RoomUseCase) GetRoomSession(ctx context.Context, roomId, userId int64) (*RoomSession, error) {
	session, err := uc.repo.GetRoomSession(ctx, roomId, userId)
	return session, err
}
