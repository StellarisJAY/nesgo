package db

import (
	"github.com/stellarisJAY/nesgo/web/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

type Engine struct {
	db *gorm.DB
}

var e = Engine{}

func init() {
	conf := config.GetConfig()
	var dialector gorm.Dialector
	switch strings.ToLower(conf.DatabaseType) {
	case "mysql":
		dialector = mysql.Open(conf.DatabaseURL)
	case "postgres":
		dialector = postgres.Open(conf.DatabaseURL)
	default:
		panic("unsupported database type")
	}
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("Can't open database")
	}
	e.db = db
}

func GetDB() *gorm.DB {
	return e.db
}
