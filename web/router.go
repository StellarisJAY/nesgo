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
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AllowHeaders = append(corsConf.AllowHeaders, "Authorization")
	r.Use(cors.New(corsConf), gin.Recovery())

	// static resources
	{
		r.StaticFS("/assets", http.Dir("web/ui/assets"))
		r.StaticFS("/scripts", http.Dir("web/ui/scripts"))
	}

	// html pages
	page := r.Group("/")
	r.LoadHTMLGlob("web/ui/*.html")
	{
		page.GET("/login", func(c *gin.Context) {
			c.HTML(200, "login.html", gin.H{})
		})
		page.GET("/room/:roomId", func(c *gin.Context) {
			c.HTML(200, "video.html", gin.H{})
		})
		page.GET("/home", func(c *gin.Context) {
			c.HTML(200, "home.html", gin.H{})
		})
		page.GET("/register", func(c *gin.Context) {
			c.HTML(200, "register.html", nil)
		})
	}

	// web api
	api := r.Group("/api")
	{
		api.POST("/user/register", userService.Register)
		api.POST("/user/login", userService.Login)
		authorized := api.Group("/", middleware.ParseQueryToken, middleware.JWTAuth)
		{
			authorized.POST("/room", roomService.CreateRoom)
			authorized.GET("/room/list", roomService.ListAllRooms)
			authorized.GET("/room/list/joined", roomService.ListJoinedRooms)
			authorized.POST("/room/:roomId/join", roomService.JoinRoom)
			authorized.GET("/room/:roomId/members", roomService.ListRoomMembers)
			authorized.GET("/room/:roomId/info", roomService.GetRoomInfo)
			authorized.GET("/room/:roomId/member", roomService.GetRoomMemberSelf)

			authorized.GET("/games", gameService.ListGames)
			authorized.GET("/game/:name", gameService.GetGameInfo)

			authorized.POST("/room/:roomId/restart", roomService.Restart)

			authorized.POST("/room/:roomId/quickSave", roomService.QuickSave)
			authorized.POST("/room/:roomId/quickLoad", roomService.QuickLoad)
		}
	}

	// websocket
	ws := r.Group("/ws", middleware.ParseQueryToken, middleware.JWTAuth)
	{
		ws.GET("/room/:roomId/rtc", roomService.ConnectRTCRoomSession)
	}

	return r
}
