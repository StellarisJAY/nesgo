package data

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGameInstanceRepo, NewGameFileRepo)

// Data .
type Data struct {
	mongo     *mongo.Client
	snowflake *snowflake.Node
}

// NewData .
func NewData(c *conf.Data, sc *conf.Server, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data"))
	clientOpts := options.Client().ApplyURI(c.Mongo.Addr)
	conn, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		logHelper.Info("closing the data resources")
		_ = conn.Disconnect(context.Background())
	}
	node, _ := snowflake.NewNode(sc.NodeId)
	return &Data{
		mongo:     conn,
		snowflake: node,
	}, cleanup, nil
}
