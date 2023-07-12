package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func parseArgsAndReadProgramFile(args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("no input game file")
	}
	file, err := os.Open(args[1])
	if err != nil {
		return nil, fmt.Errorf("can't open game file %w", err)
	}
	program, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read game file error %w", err)
	}
	log.Printf("loaded program: %s, size: %d", args[1], len(program))
	return program, nil
}

func main() {
	args := os.Args
	var program []byte
	if p, err := parseArgsAndReadProgramFile(args); err != nil {
		panic(err)
	} else {
		program = p
	}
	emulator := NewEmulator(program)
	emulator.LoadAndRun()
}
