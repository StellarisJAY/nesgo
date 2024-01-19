package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/cartridge"
	"github.com/stellarisJAY/nesgo/web/model/game"
)

type GameService struct{}

func NewGameService() *GameService {
	return &GameService{}
}

func (g *GameService) GetGameInfo(c *gin.Context) {
	info, err := game.GetGameInfo(c.Param("name"))
	if err != nil {
		if errors.Is(err, cartridge.ErrUnknownCartridgeFormat) || errors.Is(err, cartridge.ErrUnsupportedMapper) {
			c.JSON(200, JSONResp{
				Status:  500,
				Message: err.Error(),
			})
			return
		}
		panic(err)
	}
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data:    info,
	})
}

func (g *GameService) ListGames(c *gin.Context) {
	games, err := game.ListGames()
	if err != nil {
		panic(err)
	}
	c.JSON(200, JSONResp{
		Status:  200,
		Message: "ok",
		Data:    games,
	})
}
