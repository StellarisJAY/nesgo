package cpu

import (
	"fmt"
)

const (
	MemorySize      int = 1 << 16 // 内存大小，64KiB
	ProgramBaseAddr     = 0x0600  // 程序代码加载到0x8000地址
	OutputBaseAddr      = 0x0200
	OutputEndAddr       = 0x0600
)

// CallbackFunc 每条指令执行前的callback，返回false将结束处理器循环
type CallbackFunc func(*Processor) bool

type Processor struct {
	regA      byte
	regX      byte
	regY      byte
	regStatus byte
	pc        uint16
	memory    [MemorySize]byte
}

func NewProcessor() Processor {
	return Processor{}
}

func (p *Processor) LoadAndRun(program []byte) {
	p.loadProgram(program)
	p.reset()
	p.run()
}

func (p *Processor) LoadAndRunWithCallback(program []byte, callback CallbackFunc) {
	p.loadProgram(program)
	p.reset()
	p.runWithCallback(callback)
}

func (p *Processor) loadProgram(program []byte) {
	copy(p.memory[ProgramBaseAddr:], program)
	p.writeMemUint16(0xFFFC, ProgramBaseAddr)
}

func (p *Processor) reset() {
	p.regX = 0
	p.regA = 0
	p.regY = 0
	p.regStatus = 0
	p.pc = p.readMemUint16(0xFFFC)
}

func (p *Processor) run() {
	for {
		opCode := p.readMemUint8(p.pc)
		p.pc++
		originalPc := p.pc
		instruction, ok := Instructions[opCode]
		if !ok {
			panic(fmt.Errorf("unknown instruction: 0x%x", opCode))
		}

		switch opCode {
		case BRK:
			return
		case NOP:
			continue
		case INX:
			p.inx()
		case INY:
			p.iny()
		default:
			instruction.handler(p, instruction)
		}
		if p.pc == originalPc {
			p.pc += uint16(instruction.length - 1)
		}
	}
}

func (p *Processor) runWithCallback(callback CallbackFunc) {
	for {
		if !callback(p) {
			break
		}
		opCode := p.readMemUint8(p.pc)
		p.pc++
		originalPc := p.pc
		instruction, ok := Instructions[opCode]
		if !ok {
			panic(fmt.Errorf("unknown instruction: 0x%x", opCode))
		}
		switch opCode {
		case BRK:
			return
		case NOP:
			continue
		case INX:
			p.inx()
		case INY:
			p.iny()
		default:
			instruction.handler(p, instruction)
		}
		if p.pc == originalPc {
			p.pc += uint16(instruction.length - 1)
		}
	}
}

func (p *Processor) readMemUint8(addr uint16) byte {
	return p.memory[addr]
}

func (p *Processor) writeMemUint8(addr uint16, val byte) {
	p.memory[addr] = val
}

// 小端序读取16bits内存
func (p *Processor) readMemUint16(addr uint16) uint16 {
	low := uint16(p.readMemUint8(addr))
	high := uint16(p.readMemUint8(addr + 1))
	return (high << 8) | low
}

// 小端序写入16bits内存
func (p *Processor) writeMemUint16(addr uint16, val uint16) {
	low := byte(val & 0xff)
	high := byte(val >> 8)
	p.writeMemUint8(addr, low)
	p.writeMemUint8(addr+1, high)
}

func (p *Processor) getMemoryAddress(mode AddressMode) uint16 {
	var addr uint16
	switch mode {
	case Immediate:
		addr = p.pc
	case ZeroPage:
		addr = uint16(p.readMemUint8(p.pc))
	case Absolute:
		addr = p.readMemUint16(p.pc)
	case ZeroPageX:
		addr = uint16(p.readMemUint8(p.pc))
		addr += uint16(p.regX)
	case ZeroPageY:
		addr = uint16(p.readMemUint8(p.pc))
		addr += uint16(p.regY)
	case AbsoluteX:
		addr = p.readMemUint16(p.pc)
		addr += uint16(p.regX)
	case AbsoluteY:
		addr = p.readMemUint16(p.pc)
		addr += uint16(p.regY)
	case IndirectX:
		ptr := uint16(p.readMemUint8(p.pc))
		addr = p.readMemUint16(ptr + uint16(p.regX))
	case IndirectY:
		ptr := uint16(p.readMemUint8(p.pc))
		addr = p.readMemUint16(ptr + uint16(p.regY))
	case NoneAddressing:
		addr = p.pc
	}
	return addr
}

// GetMemoryRange 获取start到end范围内的内存切片
func (p *Processor) GetMemoryRange(start, end uint16) []byte {
	return p.memory[start:end]
}
