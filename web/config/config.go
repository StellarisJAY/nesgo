package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DatabaseURL  string `yaml:"db_url"`
	DatabaseType string `yaml:"db_type"`
	JwtSecret    string `yaml:"jwt_secret"`
	RedisAddr    string `yaml:"redis_addr"`
}

var config Config = Config{}

func init() {
	parseConfigs()
}

func GetConfig() Config {
	return config
}

func parseConfigs() {
	file, err := os.OpenFile("config.yml", os.O_RDONLY, 0444)
	if err != nil {
		panic("Can't find config.yml")
	}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		panic(fmt.Errorf("parse config.yml error %s", err))
	}
}
