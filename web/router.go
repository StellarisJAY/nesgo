package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stellarisJAY/nesgo/emulator/config"
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
	if config.GetEmulatorConfig().Debug {
		r.Use(noCache)
	}
	r.NoRoute(func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.StaticFS("/assets", http.Dir("ui/dist/assets"))
	r.StaticFS("/scripts", http.Dir("web/ui/scripts"))
	r.LoadHTMLFiles("ui/dist/index.html")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	//// html pages
	//page := r.Group("/")
	//r.LoadHTMLGlob("web/ui/*.html")
	//{
	//	page.GET("/login", func(c *gin.Context) {
	//		c.HTML(200, "login.html", gin.H{})
	//	})
	//	page.GET("/room/:roomId", func(c *gin.Context) {
	//		c.HTML(200, "room.html", gin.H{})
	//	})
	//	page.GET("/home", func(c *gin.Context) {
	//		c.HTML(200, "home.html", gin.H{})
	//	})
	//	page.GET("/register", func(c *gin.Context) {
	//		c.HTML(200, "register.html", nil)
	//	})
	//}

	// web api
	api := r.Group("/api")
	{
		api.POST("/user/register", userService.Register)
		api.POST("/user/login", userService.Login)
		// only authorized user can access these apis:
		authorized := api.Group("/", middleware.ParseQueryToken, middleware.AuthHandler)
		{
			authorized.POST("/room", roomService.CreateRoom)
			authorized.GET("/room/list", roomService.ListAllRooms)
			authorized.GET("/room/list/joined", roomService.ListJoinedRooms)
			authorized.POST("/room/:roomId/join", roomService.JoinRoom)
			authorized.GET("/room/:roomId/info", roomService.GetRoomInfo)
			authorized.GET("/games", gameService.ListGames)
			authorized.GET("/game/:name", gameService.GetGameInfo)
			authorized.GET("/room/search", roomService.Search)
		}
		// only room host can access these apis:
		hostApis := authorized.Group("/", roomService.HostAccessible())
		{
			hostApis.POST("/room/:roomId/restart", roomService.Restart)
			hostApis.POST("/room/:roomId/quickSave", roomService.QuickSave)
			hostApis.POST("/room/:roomId/load/:saveId", roomService.QuickLoad)
			hostApis.POST("/room/:roomId/control/transfer", roomService.TransferControl)
			hostApis.POST("/room/:roomId/member/kick", roomService.KickMember)
			hostApis.POST("/room/:roomId/role", roomService.AlterRole)
			hostApis.POST("/room/:roomId/delete", roomService.DeleteRoom)
			hostApis.POST("/room/:roomId/saves/:saveId/delete", roomService.DeleteSave)
		}
		// only room member can access these apis:
		roomMemberApis := authorized.Group("/", roomService.MemberAccessible())
		{
			roomMemberApis.GET("/room/:roomId/members", roomService.ListRoomMembers)
			roomMemberApis.GET("/room/:roomId/member", roomService.GetRoomMemberSelf)
			roomMemberApis.GET("/room/:roomId/saves", roomService.ListSaves)
			roomMemberApis.POST("/room/:roomId/leave", roomService.Leave)
		}
	}

	// websocket
	ws := r.Group("/ws", middleware.ParseQueryToken, middleware.AuthHandler)
	{
		ws.GET("/room/:roomId/rtc", roomService.MemberAccessible(), roomService.ConnectRTCRoomSession)
	}

	return r
}

func noCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
}
