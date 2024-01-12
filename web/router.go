package web

import (
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/config"
	"github.com/stellarisJAY/nesgo/web/middleware"
	"github.com/stellarisJAY/nesgo/web/service"
	"net/http"
)

var (
	gameService *service.GameService
	userService = service.NewUserService()
	roomService = service.NewRoomService()
)

func setupRouter() *gin.Engine {
	conf := config.GetEmulatorConfig()
	gameService = service.NewGameService(conf)
	r := gin.Default()
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("web/ui/*.html")
	r.POST("/user/register", userService.Register)
	r.POST("/user/login", userService.Login)
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})
	r.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, "room.html", gin.H{})
	})
	r.GET("/game/:name", gameService.HandleGamePage)

	// javascript files
	r.StaticFS("/assets", http.Dir("web/ui/assets"))
	r.StaticFS("/scripts", http.Dir("web/ui/scripts"))

	authorized := r.Group("/", middleware.ParseQueryToken, middleware.JWTAuth)
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
	authorized.POST("/room/:roomId/join", roomService.JoinRoom)
	authorized.GET("/room/:roomId/members", roomService.ListRoomMembers)
	authorized.GET("/room/:roomId/info", roomService.GetRoomInfo)

	authorized.POST("/room/:roomId/start", roomService.StartGame)
	authorized.GET("/ws/room/:roomId", roomService.HandleWebsocket)
	return r
}
