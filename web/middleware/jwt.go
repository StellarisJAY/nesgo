package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stellarisJAY/nesgo/web/config"
	"github.com/stellarisJAY/nesgo/web/model/db"
	"github.com/stellarisJAY/nesgo/web/model/user"
	"strconv"
	"time"
)

const tokenExpire = time.Hour
const tokenVersion = 1

var (
	ErrInvalidToken = errors.New("invalid token string")
	ErrExpiredToken = errors.New("token expired")
)
var secret []byte

type Claim struct {
	UserId     int64
	UserName   string
	ExpireTime time.Time
	Issuer     string
	Version    int
}

func (c *Claim) Valid() error {
	if c.Issuer != "nesgo" && c.Version != tokenVersion {
		return ErrInvalidToken
	}
	return nil
}

func init() {
	secret = []byte(config.GetConfig().JwtSecret)
}

func ParseQueryToken(c *gin.Context) {
	value := c.Query("auth")
	if value != "" && c.Request.Header.Get("Authorization") == "" {
		c.Request.Header.Set("Authorization", value)
	}
	c.Next()
}

func JWTAuth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.Abort()
		c.Status(403)
		return
	}
	claim, err := parseToken(tokenString)
	// valid token
	if err == nil {
		c.AddParam("uid", strconv.FormatInt(claim.UserId, 10))
		c.Next()
		return
	}
	// expired token
	if errors.Is(err, ErrExpiredToken) {
		key := "user_" + strconv.FormatInt(claim.UserId, 10)
		client := db.GetRedis()
		_, err := client.Get(key).Result()
		if errors.Is(err, redis.Nil) {
			// login expired
			c.JSON(403, gin.H{"status": 403, "message": "login expired"})
			c.Abort()
		} else {
			// renew access token
			tokenString, _ := GenerateToken(&user.User{
				Id:   claim.UserId,
				Name: claim.UserName,
			})
			c.Header("refresh-token", tokenString)
			c.AddParam("uid", strconv.FormatInt(claim.UserId, 10))
			c.Next()
			return
		}
	} else {
		c.Status(403)
		c.Abort()
	}
}

func GenerateToken(usr *user.User) (string, error) {
	expireTime := time.Now().Add(tokenExpire)
	claim := Claim{
		UserId:     usr.Id,
		UserName:   usr.Name,
		Issuer:     "nesgo",
		ExpireTime: expireTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim)
	signedString, err := token.SignedString(secret)
	return signedString, err
}

func parseToken(tokenString string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(_ *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claim, ok := token.Claims.(*Claim); ok && token.Valid {
		if time.Now().After(claim.ExpireTime) {
			return claim, ErrExpiredToken
		}
		return claim, nil
	} else {
		return nil, ErrInvalidToken
	}
}
