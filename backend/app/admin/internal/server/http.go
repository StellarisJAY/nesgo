package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/admin/v1"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/biz"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/conf"
	"github.com/stellarisJAY/nesgo/backend/app/admin/internal/service"
)

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/nesgo.admin.v1.Admin/Login"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, as *service.AdminService, ac *conf.Auth, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			selector.Server( // jwt 验证
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte(ac.Secret), nil
				}, jwt.WithSigningMethod(jwt2.SigningMethodHS256), jwt.WithClaims(func() jwt2.Claims {
					return &biz.AdminClaims{}
				})),
			).Match(NewWhiteListMatcher()).Build(),
			logging.Server(logger),
		),
		http.Filter(handlers.CORS( // 浏览器跨域
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	route := srv.Route("/")
	route.POST("api/v1/admin/game/upload", as.HandleUploadGame)
	v1.RegisterAdminHTTPServer(srv, as)
	return srv
}
