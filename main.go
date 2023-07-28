package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type Configs struct {
	game        string // game 游戏文件路径
	enableTrace bool   // enableTrace 是否在控制台打印trace
	disassemble bool
}

func parseConfigs() (conf Configs) {
	flag.StringVar(&conf.game, "game", "", "game file path")
	flag.BoolVar(&conf.enableTrace, "trace", false, "enable debug tracing")
	flag.BoolVar(&conf.disassemble, "disassemble", false, "disassemble program")
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
	conf := parseConfigs()
	var program []byte
	if p, err := readProgramFile(conf.game); err != nil {
		panic(err)
	} else {
		program = p
	}
	emulator := NewEmulator(program)
	if conf.disassemble {
		emulator.Disassemble()
	} else {
		emulator.LoadAndRun(conf.enableTrace)
	}
}
