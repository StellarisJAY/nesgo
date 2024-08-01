package data

import (
	"github.com/stellarisJAY/nesgo/backend/app/gaming/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGameInstanceRepo, NewGameFileRepo)

// Data .
type Data struct {
}

// NewData .
func NewData(c *conf.Registry, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "data"))
	cleanup := func() {
		logHelper.Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}
