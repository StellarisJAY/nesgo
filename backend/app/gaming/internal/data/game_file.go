package data

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/biz"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const (
	gameFileBucketName = "game_file"
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
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(gameFileBucketName))
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

func (g *gameFileRepo) UploadGameData(ctx context.Context, game string, data []byte, metadata *biz.GameFileMetadata) error {
	db := g.data.mongo.Database("nesgo")
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(gameFileBucketName))
	if err != nil {
		return err
	}
	err = bucket.GetFilesCollection().FindOne(ctx, bson.M{"filename": game}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		start := time.Now()
		_, err = bucket.UploadFromStream(game, bytes.NewBuffer(data), options.GridFSUpload().SetMetadata(metadata))
		g.logger.Infof("uploaded game ,time: %dms", time.Now().Sub(start).Milliseconds())
		return err
	} else {
		return err
	}
}

func (g *gameFileRepo) ListGames(ctx context.Context, page, pageSize int) ([]*biz.GameFileMetadata, int, error) {
	db := g.data.mongo.Database("nesgo")
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(gameFileBucketName))
	if err != nil {
		return nil, 0, err
	}

	total, err := bucket.GetFilesCollection().CountDocuments(ctx, bson.M{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []*biz.GameFileMetadata{}, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	offset := int64(page * pageSize)
	opts := options.Find().SetSkip(offset).SetLimit(int64(pageSize))
	cursor, err := bucket.GetFilesCollection().Find(ctx, bson.M{}, opts)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []*biz.GameFileMetadata{}, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var result []*biz.GameFileMetadata
	for cursor.Next(ctx) {
		metadata := new(biz.GameFileMetadata)
		err = cursor.Current.Lookup("metadata").Unmarshal(&metadata)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, metadata)
	}
	return result, int(total), nil
}

func (g *gameFileRepo) DeleteGames(ctx context.Context, games []string) (int, error) {
	db := g.data.mongo.Database("nesgo")
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(gameFileBucketName))
	if err != nil {
		return 0, err
	}
	cursor, err := bucket.GetFilesCollection().Find(ctx, bson.M{"filename": bson.M{"$in": games}})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return 0, nil
	}
	count := 0
	for cursor.Next(ctx) {
		id := cursor.Current.Lookup("_id").ObjectID()
		err = bucket.DeleteContext(ctx, id)
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

func (g *gameFileRepo) GetSavedGame(ctx context.Context, id int64) (*biz.GameSave, error) {
	//TODO implement me
	panic("implement me")
}

func (g *gameFileRepo) SaveGame(ctx context.Context, save *biz.GameSave) error {
	//TODO implement me
	panic("implement me")
}
