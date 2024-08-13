package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"golang.org/x/net/context"
)

func Get[T any](_ context.Context, rdb *redis.Client, key string) (*T, error) {
	data, err := rdb.Get(key).Bytes()
	if err != nil {
		return nil, err
	}
	res := new(T)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func Set[T any](_ context.Context, rdb *redis.Client, key string, value *T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(key, data, 0).Err()
}

func Del(_ context.Context, rdb *redis.Client, key string) error {
	return rdb.Del(key).Err()
}
