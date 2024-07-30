package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"github.com/stellarisjay/nesgo/backend/app/user/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo)

// Data .
type Data struct {
	rdb    *redis.Client
	db     *gorm.DB
	logger *log.Helper
}

func NewData(c *conf.Data) (*Data, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: c.Redis.Addr,
	})
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if !db.Migrator().HasTable(&User{}) {
		if err := db.Migrator().CreateTable(&User{}); err != nil {
			panic(err)
		}
	}

	logger := log.NewHelper(log.DefaultLogger)
	cleanup := func() {
		logger.Info("cleaning up database connections")
		_ = rdb.Close()
	}
	return &Data{
		rdb:    rdb,
		db:     db,
		logger: logger,
	}, cleanup, nil
}
