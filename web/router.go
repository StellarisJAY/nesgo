package web

import (
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/config"
	"sync"
)

var gameService *GameService

func setupRouter(config config.Config) *gin.Engine {
	gameService = &GameService{
		mutex:    &sync.Mutex{},
		sessions: make(map[string]*GameSession),
		conf:     config,
	}
	r := gin.Default()
	r.LoadHTMLGlob("web/ui/*")
	// game page
	r.GET("/game/:name", gameService.HandleGamePage)
	// game session websocket
	r.GET("/ws/:id", gameService.HandleWebsocket)
	return r
}
