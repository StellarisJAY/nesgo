package trace

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/nes/cpu"
	"strings"
)

func Trace(p *cpu.Processor, instruction *cpu.Instruction) {
	pc := p.ProgramCounter()
	var addrMode string
	registers := p.DumpRegisters()
	switch instruction.AddrMode {
	case cpu.NoneAddressing:
		addrMode = "none"
	case cpu.Immediate:
		addrMode = "im"
	case cpu.ZeroPage:
		addrMode = "zp"
	case cpu.ZeroPageX:
		addrMode = "zpx"
	case cpu.ZeroPageY:
		addrMode = "zpy"
	case cpu.Absolute:
		addrMode = "abs"
	case cpu.AbsoluteX:
		addrMode = "abx"
	case cpu.AbsoluteY:
		addrMode = "aby"
	case cpu.IndirectX:
		addrMode = "inx"
	case cpu.IndirectY:
		addrMode = "iny"
	}

	argc := uint16(instruction.Length - 1)
	args := strings.Builder{}
	var i uint16
	for i = 0; i < argc; i++ {
		args.WriteString(fmt.Sprintf("%02X ", p.ReadMem8(pc+i)))
	}
	fmt.Printf("%04X\t%5s_%4s\t%8s\t%s CYC:%d\n", pc-1, instruction.Name, addrMode, args.String(), registers, p.Cycles())
}

func PrintDisassemble(p *cpu.Processor, instruction *cpu.Instruction) {
	pc := p.ProgramCounter()
	var addrMode string
	switch instruction.AddrMode {
	case cpu.NoneAddressing:
		addrMode = "none"
	case cpu.Immediate:
		addrMode = "im"
	case cpu.ZeroPage:
		addrMode = "zp"
	case cpu.ZeroPageX:
		addrMode = "zpx"
	case cpu.ZeroPageY:
		addrMode = "zpy"
	case cpu.Absolute:
		addrMode = "abs"
	case cpu.AbsoluteX:
		addrMode = "abx"
	case cpu.AbsoluteY:
		addrMode = "aby"
	case cpu.IndirectX:
		addrMode = "inx"
	case cpu.IndirectY:
		addrMode = "iny"
	}

	argc := uint16(instruction.Length - 1)
	args := strings.Builder{}
	var i uint16
	for i = 0; i < argc; i++ {
		args.WriteString(fmt.Sprintf("%02X ", p.ReadMem8(pc+i)))
	}
	fmt.Printf("%04X\t%5s_%4s\t%8s\n", pc-1, instruction.Name, addrMode, args.String())
}
