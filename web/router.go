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
	// boost game session
	r.POST("/game/:id/boost/:rate", gameService.HandleCPUBoost)

	r.POST("/game/:id/pause", gameService.HandlePauseGame)
	r.POST("/game/:id/resume", gameService.HandleResumeGame)

	r.POST("/game/:id/reverse", gameService.HandleReverseOne)
	return r
}
