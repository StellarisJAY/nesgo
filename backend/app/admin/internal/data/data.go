package data

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	gaming "github.com/stellarisJAY/nesgo/backend/api/gaming/service/v1"
	roomAPI "github.com/stellarisJAY/nesgo/backend/api/room/service/v1"
	userAPI "github.com/stellarisJAY/nesgo/backend/api/user/service/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGameFileRepo, NewAdminRepo)

// Data .
type Data struct {
	gamingCli gaming.GamingClient
	roomCli   roomAPI.RoomClient
	userCli   userAPI.UserClient
	db        *gorm.DB
	logger    *log.Helper
}

// NewData .
func NewData(c *conf.Data, r registry.Discovery, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data"))
	gamingConn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///nesgo.service.gaming"),
		grpc.WithDiscovery(r))
	if err != nil {
		panic(err)
	}

	roomConn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///nesgo.service.room"),
		grpc.WithDiscovery(r))
	if err != nil {
		panic(err)
	}

	userConn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///nesgo.service.user"),
		grpc.WithDiscovery(r))
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(c.Database.Source))
	if err != nil {
		panic(err)
	}
	if !db.Migrator().HasTable(&Admin{}) {
		if err := db.Migrator().CreateTable(&Admin{}); err != nil {
			panic(err)
		}
	}
	err = db.Model(&Admin{}).Where("name=?", "admin").First(&Admin{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		sum := md5.Sum([]byte("admin"))
		password := hex.EncodeToString(sum[:])
		err = db.Model(&Admin{}).Create(&Admin{
			Name:      "admin",
			Password:  password,
			CreatedAt: time.Now(),
		}).Error
		if err != nil {
			panic(err)
		}
	}
	cleanup := func() {
		logHelper.Info("closing data resources")
		_ = gamingConn.Close()
		_ = roomConn.Close()
		_ = userConn.Close()
	}
	return &Data{
		logger:    logHelper,
		gamingCli: gaming.NewGamingClient(gamingConn),
		roomCli:   roomAPI.NewRoomClient(gamingConn),
		userCli:   userAPI.NewUserClient(userConn),
		db:        db,
	}, cleanup, nil
}
