package trace

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/cpu"
)

func Trace(p *cpu.Processor, instruction cpu.Instruction) {
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
	fmt.Printf("%4x\t%5s_%4s\t%s\n", pc-1, instruction.Name, addrMode, registers)
}
