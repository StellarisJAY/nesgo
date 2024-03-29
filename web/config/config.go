package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DatabaseURL  string `yaml:"db_url"`
	DatabaseType string `yaml:"db_type"`
	RedisAddr    string `yaml:"redis_addr"`

	FileStorageType          string `yaml:"file_storage_type"`
	HostFileSystemStorageDir string `yaml:"host_fs_store_dir"`

	TurnServer struct {
		Addr             string `yaml:"addr"`
		LongTermUser     string `yaml:"long_term_user"`
		LongTermPassword string `yaml:"long_term_password"`
	} `yaml:"turn_server"`
}

var config Config = Config{
	DatabaseURL:              "",
	DatabaseType:             "mysql",
	RedisAddr:                "localhost:6369",
	HostFileSystemStorageDir: "/var/tmp/nesgo/saves",
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
