package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Config struct {
	game          string        // game 游戏文件路径
	enableTrace   bool          // enableTrace 是否在控制台打印trace
	disassemble   bool          // disassemble 打印程序的反汇编结果
	scale         int           // scale 屏幕放大尺寸，原始尺寸：256x240像素
	frameInterval time.Duration // frameInterval 每一帧画面渲染间隔时间
}

func parseConfig() (conf Config) {
	flag.StringVar(&conf.game, "game", "games/super.nes", "game file path")
	flag.BoolVar(&conf.enableTrace, "trace", false, "enable debug tracing")
	flag.BoolVar(&conf.disassemble, "disassemble", false, "disassemble program")
	flag.IntVar(&conf.scale, "scale", 1, "game screen scale")
	flag.DurationVar(&conf.frameInterval, "interval", 1*time.Millisecond, "interval between each frame")
	flag.Parse()
	return
}

func readProgramFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("can't open game file %s,  %w", fileName, err)
	}
	program, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read game file error %w", err)
	}
	log.Printf("loaded program file: %s, size: %d", fileName, len(program))
	return program, nil
}

func main() {
	conf := parseConfig()
	var program []byte
	if p, err := readProgramFile(conf.game); err != nil {
		panic(err)
	} else {
		program = p
	}
	emulator := NewEmulator(program, conf)
	if conf.disassemble {
		emulator.Disassemble()
	} else {
		emulator.LoadAndRun(conf.enableTrace)
	}
}
