package web

import (
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/web/middleware"
	"github.com/stellarisJAY/nesgo/web/service"
)

var (
	gameService *service.GameService
	userService = service.NewUserService()
	roomService = service.NewRoomService()
)

func setupRouter(config config.Config) *gin.Engine {
	gameService = service.NewGameService(config)
	r := gin.Default()
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("web/ui/*")
	r.POST("/user/register", userService.Register)
	r.POST("/user/login", userService.Login)

	authorized := r.Group("/", middleware.JWTAuth)
	// game page
	authorized.GET("/game/:name", gameService.HandleGamePage)
	// game session websocket
	authorized.GET("/ws/:id", gameService.HandleWebsocket)
	// boost game session
	authorized.POST("/game/:id/boost/:rate", gameService.HandleCPUBoost)
	authorized.POST("/game/:id/pause", gameService.HandlePauseGame)
	// pause & resume
	authorized.POST("/game/:id/resume", gameService.HandleResumeGame)
	authorized.POST("/game/:id/reverse", gameService.HandleReverseOne)

	authorized.POST("/room", roomService.CreateRoom)
	authorized.GET("/room/list", roomService.ListOwningRooms)
	return r
}
