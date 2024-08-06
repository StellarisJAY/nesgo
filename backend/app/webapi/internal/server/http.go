package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/stellarisJAY/nesgo/backend/api/app/webapi/v1"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/conf"
	"github.com/stellarisJAY/nesgo/backend/app/webapi/internal/service"
)

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/login"] = struct{}{}
	whiteList["/register"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, ac *conf.Auth, webApiService *service.WebApiService, logger log.Logger) *http.Server {
	//tokenFunc := func(_ *jwt2.Token) (interface{}, error) {
	//	return []byte(ac.Secret), nil
	//}
	//claimsFunc := func() jwt2.Claims {
	//	return &jwt2.MapClaims{}
	//}
	//
	//selectorServ := selector.Server(jwt.Server(tokenFunc, jwt.WithSigningMethod(jwt2.SigningMethodES256), jwt.WithClaims(claimsFunc)))
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			//selectorServ.Match(NewWhiteListMatcher()).Build(),
		),
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
	v1.RegisterWebApiHTTPServer(srv, webApiService)
	return srv
}
