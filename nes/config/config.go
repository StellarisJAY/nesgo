package config

import (
	"flag"
	"time"
)

type Config struct {
	Game               string        // Game 游戏文件路径
	EnableTrace        bool          // EnableTrace 是否在控制台打印trace
	Disassemble        bool          // Disassemble 打印程序的反汇编结果
	Scale              int           // Scale 屏幕放大尺寸，原始尺寸：256x240像素
	FrameInterval      time.Duration // FrameInterval 每一帧画面渲染间隔时间
	ServerAddr         string
	GameDirectory      string
	SaveDirectory      string
	SnapshotSerializer string

	MuteApu bool

	Debug bool
}

var conf Config

func InitConfigs() {
	flag.StringVar(&conf.Game, "game", "games/super.nes", "Game file path")
	flag.BoolVar(&conf.EnableTrace, "trace", false, "enable debug tracing")
	flag.BoolVar(&conf.Disassemble, "disassemble", false, "Disassemble program")
	flag.IntVar(&conf.Scale, "scale", 1, "Game screen Scale")
	flag.DurationVar(&conf.FrameInterval, "interval", 1*time.Millisecond, "interval between each frame")
	flag.StringVar(&conf.ServerAddr, "addr", ":8080", "Web server addr")
	flag.StringVar(&conf.GameDirectory, "dir", "", "Game directory")
	flag.StringVar(&conf.SaveDirectory, "save-dir", "", "Game save directory")
	flag.StringVar(&conf.SnapshotSerializer, "serializer", "gob", "Game save serializer")
	flag.BoolVar(&conf.MuteApu, "mute", false, "Mute APU")
	flag.BoolVar(&conf.Debug, "debug", false, "use debug mode")
	flag.Parse()
}

func GetEmulatorConfig() Config {
	return conf
}
