package data

import (
	"github.com/go-redis/redis"
	"github.com/stellarisjay/nesgo/backend/app/room/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRoomRepo)

// Data .
type Data struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data"))
	db, err := gorm.Open(mysql.Open(c.Database.Source))
	if err != nil {
		panic(err)
	}

	if !db.Migrator().HasTable(&Room{}) {
		if err := db.Migrator().CreateTable(&Room{}); err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&RoomMember{}) {
		if err := db.Migrator().CreateTable(&RoomMember{}); err != nil {
			panic(err)
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: c.Redis.Addr,
	})
	cleanup := func() {
		logHelper.Info("closing the data resources")
		_ = rdb.Close()
	}
	return &Data{
		db:     db,
		rdb:    rdb,
		logger: logHelper,
	}, cleanup, nil
}
