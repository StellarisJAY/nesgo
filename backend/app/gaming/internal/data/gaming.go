package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/biz"
	etcdAPI "go.etcd.io/etcd/client/v3"
	"sync"
)

type gameInstanceRepo struct {
	instances map[int64]*biz.GameInstance
	mutex     *sync.RWMutex
	logger    *log.Helper
	etcdCli   *etcdAPI.Client
	leaseID   etcdAPI.LeaseID
}

func NewGameInstanceRepo(etcdCli *etcdAPI.Client, logger log.Logger) biz.GameInstanceRepo {
	// 创建lease，用来与房间session绑定，lease失效自动删除session
	lease, err := etcdCli.Lease.Grant(context.Background(), 30)
	if err != nil {
		panic(err)
	}
	respChan, err := etcdCli.Lease.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case _, ok := <-respChan:
				if !ok {
					log.NewHelper(log.With(logger, "module", "keepalive")).Infof("lease %d keepalive closed", lease.ID)
				}
			}
		}
	}()
	return &gameInstanceRepo{
		instances: make(map[int64]*biz.GameInstance),
		mutex:     &sync.RWMutex{},
		logger:    log.NewHelper(log.With(logger, "module", "data/gameInstance")),
		etcdCli:   etcdCli,
		leaseID:   lease.ID,
	}
}

func (g *gameInstanceRepo) CreateGameInstance(ctx context.Context, game *biz.GameInstance) (int64, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.instances[game.RoomId] = game
	return int64(g.leaseID), nil
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
