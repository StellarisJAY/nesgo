package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/user/internal/biz"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type macroRepo struct {
	data   *Data
	logger *log.Helper
}

const macroCollectionName = "macros"

func NewMacroRepo(data *Data, logger log.Logger) biz.MacroRepo {
	return &macroRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/macro")),
	}
}

func (m *macroRepo) CreateMacro(ctx context.Context, macro *biz.Macro) error {
	collection := m.data.mongo.Database("nesgo").Collection(macroCollectionName)
	macro.Id = m.data.snowflake.Generate().Int64()
	_, err := collection.InsertOne(ctx, macro)
	return err
}

func (m *macroRepo) GetMacro(ctx context.Context, id int64) (*biz.Macro, error) {
	collection := m.data.mongo.Database("nesgo").Collection(macroCollectionName)
	var macro biz.Macro
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&macro)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &macro, nil
}

func (m *macroRepo) ListMacro(ctx context.Context, userId int64, page, pageSize int32) ([]*biz.Macro, int32, error) {
	collection := m.data.mongo.Database("nesgo").Collection(macroCollectionName)
	var macros []*biz.Macro
	total, err := collection.CountDocuments(ctx, bson.M{"userId": userId})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	skip := int64(page * pageSize)
	limit := int64(pageSize)
	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := collection.Find(ctx, bson.M{"userId": userId}, opts)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, int32(total), nil
	}
	if err != nil {
		return nil, 0, err
	}
	err = cursor.All(ctx, &macros)
	if err != nil {
		return nil, 0, err
	}
	return macros, int32(total), nil
}

func (m *macroRepo) DeleteMacro(ctx context.Context, id int64) error {
	collection := m.data.mongo.Database("nesgo").Collection(macroCollectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (m *macroRepo) GetMacroByName(ctx context.Context, userId int64, name string) (*biz.Macro, error) {
	collection := m.data.mongo.Database("nesgo").Collection(macroCollectionName)
	var macro biz.Macro
	err := collection.FindOne(ctx, bson.M{"userId": userId, "name": name}).Decode(&macro)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &macro, nil
}
