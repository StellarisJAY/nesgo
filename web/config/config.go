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

	FileStorageType          string `yaml:"file_storage_type"`
	HostFileSystemStorageDir string `yaml:"host_fs_store_dir"`

	TurnServerAddr string `yaml:"turn_server_addr"`
}

var config Config = Config{
	DatabaseURL:              "",
	DatabaseType:             "mysql",
	JwtSecret:                "123456",
	RedisAddr:                "localhost:6369",
	HostFileSystemStorageDir: "/var/tmp/nesgo/saves",
	TurnServerAddr:           "turn:192.168.0.107:3478",
}

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
