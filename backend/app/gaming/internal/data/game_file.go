package data

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
	saveFileBucketName = "save_file"
)

type gameFileRepo struct {
	logger        *log.Helper
	data          *Data
	cacheMutex    *sync.RWMutex
	gameFileCache map[string][]byte
}

type saveFileMetadata struct {
	RoomId     int64  `json:"roomId" bson:"roomId"`
	Id         int64  `json:"id" bson:"id"`
	Game       string `json:"game" bson:"game"`
	CreateTime int64  `json:"createTime" bson:"createTime"`
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

func saveFileName(id int64) string {
	return fmt.Sprintf("nesgo_save_%d", id)
}

func (g *gameFileRepo) GetSavedGame(ctx context.Context, id int64) (*biz.GameSave, error) {
	db := g.data.mongo.Database("nesgo")
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(saveFileBucketName))
	if err != nil {
		return nil, err
	}
	filename := saveFileName(id)
	save := &biz.GameSave{}
	result := bucket.GetFilesCollection().FindOne(ctx, bson.M{"filename": filename})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	raw, _ := result.Raw()
	_ = raw.Lookup("metadata").Unmarshal(save)

	buffer := &bytes.Buffer{}
	_, err = bucket.DownloadToStreamByName(filename, buffer)
	if err != nil {
		return nil, err
	}
	save.Data = buffer.Bytes()
	return save, nil
}

func (g *gameFileRepo) SaveGame(ctx context.Context, save *biz.GameSave) error {
	db := g.data.mongo.Database("nesgo")
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(saveFileBucketName))
	if err != nil {
		return err
	}
	id := g.data.snowflake.Generate().Int64()
	filename := saveFileName(id)
	save.Id = id
	reader := bytes.NewReader(save.Data)
	metadata := &saveFileMetadata{
		RoomId:     save.RoomId,
		Id:         save.Id,
		Game:       save.Game,
		CreateTime: save.CreateTime,
	}
	_, err = bucket.UploadFromStream(filename, reader, options.GridFSUpload().SetMetadata(metadata))
	return err
}

func (g *gameFileRepo) ListSaves(ctx context.Context, roomId int64, page, pageSize int32) ([]*biz.GameSave, int32, error) {
	db := g.data.mongo.Database("nesgo")
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(saveFileBucketName))
	if err != nil {
		return nil, 0, err
	}
	total, err := bucket.GetFilesCollection().CountDocuments(ctx, bson.M{"metadata.roomId": roomId})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().SetSkip(int64(page * pageSize)).SetLimit(int64(pageSize)).SetSort(bson.M{"metadata.createTime": -1})
	cursor, err := bucket.GetFilesCollection().Find(ctx, bson.M{"metadata.roomId": roomId}, opts)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, int32(total), nil
	}
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	result := make([]*biz.GameSave, 0, pageSize)
	for cursor.Next(ctx) {
		save := new(biz.GameSave)
		err := cursor.Current.Lookup("metadata").Unmarshal(save)
		if err == nil {
			result = append(result, save)
		}
	}
	return result, int32(total), nil
}
