package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/stellarisJAY/nesgo/web/model/db"
	"log"
	"path"
	"strconv"
	"time"
)

type Authorization struct {
	UserId    int64  `json:"user_id"`
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

func ParseQueryToken(c *gin.Context) {
	value := c.Query("auth")
	if value != "" && c.Request.Header.Get("Authorization") == "" {
		c.Request.Header.Set("Authorization", value)
	}
	c.Next()
}

func AuthHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	ip := c.RemoteIP()
	ua := c.GetHeader("User-Agent")

	result, err := db.GetRedis().Get(path.Join("auth", tokenString)).Result()
	if errors.Is(err, redis.Nil) {
		c.Status(401)
		c.Abort()
		return
	} else if err != nil {
		panic(err)
	}
	var auth Authorization
	_ = json.Unmarshal([]byte(result), &auth)

	if ip != auth.Ip || ua != auth.UserAgent {
		c.Status(401)
		c.Abort()
		return
	}

	ttl, err := db.GetRedis().TTL(tokenString).Result()
	if errors.Is(err, redis.Nil) {
		c.Status(401)
		c.Abort()
		return
	} else if err != nil {
		panic(err)
	}
	if ttl <= time.Minute*10 {
		_, err = db.GetRedis().Expire(tokenString, time.Minute*30).Result()
		if err != nil {
			log.Println("reset expire for auth error:", err)
		}
	}
	c.Set("uid", auth.UserId)
	c.Next()
}

func StoreAuthorization(a Authorization) (string, error) {
	data, _ := json.Marshal(a)
	sum := md5.Sum([]byte(a.Ip + strconv.FormatInt(time.Now().UnixMilli(), 16)))
	token := hex.EncodeToString(sum[:])
	if err := db.GetRedis().
		Set(path.Join("auth", token), data, time.Minute*30).
		Err(); err != nil {
		return "", err
	}
	return token, nil
}
