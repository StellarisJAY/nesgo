package db

import (
	"github.com/go-redis/redis"
	"github.com/stellarisJAY/nesgo/web/config"
)

var cli *redis.Client

func init() {
	conf := config.GetConfig()
	cli = redis.NewClient(&redis.Options{
		Addr: conf.RedisAddr,
	})
}

func GetRedis() *redis.Client {
	return cli
}
