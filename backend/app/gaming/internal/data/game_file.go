package data

import (
	"bytes"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/biz"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type gameFileRepo struct {
	logger        *log.Helper
	data          *Data
	cacheMutex    *sync.RWMutex
	gameFileCache map[string][]byte
}

func NewGameFileRepo(data *Data, logger log.Logger) biz.GameFileRepo {
	return &gameFileRepo{
		logger:        log.NewHelper(log.With(logger, "module", "data/gameFile")),
		data:          data,
		cacheMutex:    &sync.RWMutex{},
		gameFileCache: make(map[string][]byte),
	}
}

func (g *gameFileRepo) GetGameData(ctx context.Context, game string) ([]byte, error) {
	g.cacheMutex.RLock()
	data, ok := g.gameFileCache[game]
	g.cacheMutex.RUnlock()
	if ok {
		return data, nil
	}
	g.cacheMutex.Lock()
	defer g.cacheMutex.Unlock()
	db := g.data.mongo.Database("nesgo")
	bucketName := "game_file"
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return nil, err
	}
	buffer := &bytes.Buffer{}
	_, err = bucket.DownloadToStreamByName(game, buffer)
	if err != nil {
		return nil, err
	}
	gameData := buffer.Bytes()
	g.gameFileCache[game] = gameData
	return gameData, nil
}

func (g *gameFileRepo) UploadGameData(ctx context.Context, game string, data []byte) error {
	db := g.data.mongo.Database("nesgo")
	bucketName := "game_file"
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return err
	}
	_, err = bucket.UploadFromStream(game, bytes.NewReader(data))
	return err
}

func (g *gameFileRepo) GetSavedGame(ctx context.Context, id int64) (*biz.GameSave, error) {
	//TODO implement me
	panic("implement me")
}

func (g *gameFileRepo) SaveGame(ctx context.Context, save *biz.GameSave) error {
	//TODO implement me
	panic("implement me")
}
