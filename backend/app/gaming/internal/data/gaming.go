package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/biz"
	"sync"
)

type gameInstanceRepo struct {
	instances map[int64]*biz.GameInstance
	mutex     *sync.RWMutex
	logger    *log.Helper
}

func NewGameInstanceRepo(logger log.Logger) biz.GameInstanceRepo {
	return &gameInstanceRepo{
		instances: make(map[int64]*biz.GameInstance),
		mutex:     &sync.RWMutex{},
		logger:    log.NewHelper(log.With(logger, "module", "data/gameInstance")),
	}
}

func (g *gameInstanceRepo) CreateGameInstance(ctx context.Context, game *biz.GameInstance) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.instances[game.RoomId] = game
	return nil
}

func (g *gameInstanceRepo) DeleteGameInstance(ctx context.Context, roomId int64) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	delete(g.instances, roomId)
	return nil
}

func (g *gameInstanceRepo) GetGameInstance(ctx context.Context, roomId int64) (*biz.GameInstance, error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.instances[roomId], nil
}
