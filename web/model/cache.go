package model

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/stellarisJAY/nesgo/web/model/db"
	"log"
)

func CacheGet[T any](key string, missFunc func(k string) (*T, error)) (*T, error) {
	cmd := db.GetRedis().Get(key)
	value, err := cmd.Result()
	if errors.Is(err, redis.Nil) {
		val, err := missFunc(key)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(val)
		_, err = db.GetRedis().Set(key, string(data), 0).Result()
		if err != nil {
			log.Println("cache set error", err)
		}
		return val, nil
	} else if err != nil {
		return nil, err
	} else {
		val := new(T)
		_ = json.Unmarshal([]byte(value), val)
		return val, nil
	}
}

func CacheHashGet[T any](cKey, hKey string, missFunc func(cKey, hKey string) (*T, error)) (*T, error) {
	result, err := db.GetRedis().HGet(cKey, hKey).Result()
	if errors.Is(err, redis.Nil) {
		val, err := missFunc(cKey, hKey)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(val)
		_, err = db.GetRedis().HSet(cKey, hKey, data).Result()
		if err != nil {
			log.Println("cache hset error", err)
		}
		return val, nil
	} else if err != nil {
		return nil, err
	} else {
		t := new(T)
		_ = json.Unmarshal([]byte(result), t)
		return t, nil
	}
}

func CacheDelete(key string) error {
	_, err := db.GetRedis().Del(key).Result()
	return err
}

func CacheHashDelete(cKey, hKey string) error {
	_, err := db.GetRedis().HDel(cKey, hKey).Result()
	return err
}
