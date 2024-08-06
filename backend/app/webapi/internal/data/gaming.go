package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/biz"
)

type gamingRepo struct {
	data   *Data
	logger *log.Helper
}

func NewGamingRepo(data *Data, logger log.Logger) biz.GamingRepo {
	return &gamingRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/gaming")),
	}
}
