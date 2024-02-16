package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/web/middleware"
	"github.com/stellarisJAY/nesgo/web/service"
	"net/http"
)

var (
	userService = service.NewUserService()
	roomService = service.NewRoomService()
	gameService = service.NewGameService()
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AllowHeaders = append(corsConf.AllowHeaders, "Authorization")
	r.Use(cors.New(corsConf), gin.Recovery())
	r.LoadHTMLGlob("web/ui/*.html")
	r.POST("/user/register", userService.Register)
	r.POST("/user/login", userService.Login)
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})
	r.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, "video.html", gin.H{})
	})

	// javascript files
	r.StaticFS("/assets", http.Dir("web/ui/assets"))
	r.StaticFS("/scripts", http.Dir("web/ui/scripts"))

	authorized := r.Group("/", middleware.ParseQueryToken, middleware.JWTAuth)

	authorized.POST("/room", roomService.CreateRoom)
	authorized.GET("/room/list", roomService.ListOwningRooms)
	authorized.POST("/room/:roomId/join", roomService.JoinRoom)
	authorized.GET("/room/:roomId/members", roomService.ListRoomMembers)
	authorized.GET("/room/:roomId/info", roomService.GetRoomInfo)

	authorized.POST("/room/:roomId/start", roomService.StartGame)
	authorized.GET("/ws/room/:roomId", roomService.HandleWebsocket)

	authorized.GET("/games", gameService.ListGames)
	authorized.GET("/game/:name", gameService.GetGameInfo)

	authorized.GET("/room/:roomId/rtc", roomService.ConnectRTCRoomSession)
	authorized.POST("/room/:roomId/restart", roomService.Restart)
	return r
}
