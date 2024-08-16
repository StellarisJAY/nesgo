package data

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/user/internal/biz"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userKeyboardBindingCollection = "user_keyboard_binding"

type userKeyboardBindingRepo struct {
	data   *Data
	logger *log.Helper
}

func NewUserKeyboardBindingRepo(data *Data, logger log.Logger) biz.UserKeyboardBindingRepo {
	return &userKeyboardBindingRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/UserKeyboardBinding")),
	}
}

func (r *userKeyboardBindingRepo) CreateKeyboardBinding(ctx context.Context, ub *biz.UserKeyboardBinding) error {
	collection := r.data.mongo.Database("nesgo").Collection(userKeyboardBindingCollection)
	ub.Id = r.data.snowflake.Generate().Int64()
	_, err := collection.InsertOne(ctx, ub)
	if err != nil {
		return err
	}
	return nil
}

func (r *userKeyboardBindingRepo) UpdateKeyboardBinding(ctx context.Context, ub *biz.UserKeyboardBinding) error {
	collection := r.data.mongo.Database("nesgo").Collection(userKeyboardBindingCollection)
	updates := bson.D{
		{"$set", bson.D{
			{"name", ub.Name},
			{"userId", ub.UserId},
			{"keyboardBindings", ub.KeyboardBindings},
		}},
	}
	_, err := collection.UpdateOne(ctx, bson.M{"id": ub.Id}, updates)
	return err
}

func (r *userKeyboardBindingRepo) DeleteKeyboardBinding(ctx context.Context, id int64) error {
	collection := r.data.mongo.Database("nesgo").Collection(userKeyboardBindingCollection)
	res := collection.FindOneAndDelete(ctx, bson.M{"id": id})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return nil
	}
	return res.Err()
}

func (r *userKeyboardBindingRepo) GetKeyboardBinding(ctx context.Context, id int64) (*biz.UserKeyboardBinding, error) {
	collection := r.data.mongo.Database("nesgo").Collection(userKeyboardBindingCollection)
	res := collection.FindOne(ctx, bson.M{"id": id})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return nil, nil
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	b := &biz.UserKeyboardBinding{}
	if err := res.Decode(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (r *userKeyboardBindingRepo) ListUserKeyboardBinding(ctx context.Context, userId int64, page, pageSize int32) ([]*biz.UserKeyboardBinding, int32, error) {
	collection := r.data.mongo.Database("nesgo").Collection(userKeyboardBindingCollection)
	total, err := collection.CountDocuments(ctx, bson.M{"userId": userId})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	offset := int64(page * pageSize)
	limit := int64(pageSize)
	opts := options.Find().SetSkip(offset).SetLimit(limit)
	cursor, err := collection.Find(ctx, bson.M{"userId": userId}, opts)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*biz.UserKeyboardBinding
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, 0, err
	}
	return results, int32(total), nil
}

func (r *userKeyboardBindingRepo) GetBindingByName(ctx context.Context, userId int64, name string) (*biz.UserKeyboardBinding, error) {
	collection := r.data.mongo.Database("nesgo").Collection(userKeyboardBindingCollection)
	res := collection.FindOne(ctx, bson.M{"userId": userId, "name": name})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return nil, nil
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	b := &biz.UserKeyboardBinding{}
	if err := res.Decode(b); err != nil {
		return nil, err
	}
	return b, nil
}
