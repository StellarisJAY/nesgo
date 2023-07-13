package cpu

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/bus"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math/rand"
	"time"
)

// 内存布局
// 屏幕像素点：[0x0200, 0x0600)，共32x32个像素点，一行32个像素
// 程序代码：[0x8000, ...)
// 栈：[0x0100, 0x01FF)，共256字节
// 上一个Input：0xFF
// 随机数：0xFE
const (
	ProgramEntryPoint        = 0x0600
	PrgROMAddr               = 0x8000 // 程序代码加载到0x8000地址的PrgROM
	OutputBaseAddr           = 0x0200
	OutputEndAddr            = 0x0600
	StackBase                = 0x0100
	StackReset               = 0xFF
	Input             uint16 = 0xFF
	RandomNumber      uint16 = 0xFE

	PrgROMEntryPointAddr   uint16 = 0xFFFC // 程序entry point在ROM的地址
	PrgROMInterruptHandler uint16 = 0xFFFA // 程序中断处理函数地址
)

const (
	CarryStatus            byte = 1 << 1
	ZeroStatus             byte = 1 << 2
	InterruptDisableStatus byte = 1 << 3
	DecimalModeStatus      byte = 1 << 4
	BreakStatus            byte = 1 << 5
	OverflowStatus         byte = 1 << 6
	NegativeStatus         byte = 1 << 7
)

const (
	ActionUp    byte = 0x77
	ActionDown  byte = 0x73
	ActionLeft  byte = 0x61
	ActionRight byte = 0x64
)

// CallbackFunc 每条指令执行前的callback，返回false将结束处理器循环
type CallbackFunc func(*Processor) bool

type Processor struct {
	regA      byte
	regX      byte
	regY      byte
	regStatus byte     // resStatus 状态寄存器，记录上一条指令的状态
	pc        uint16   // pc 程序计数器
	sp        byte     // sp 栈指针，记录栈地址的低8位，高位固定为0x0100
	bus       *bus.Bus // bus 总线，通过总线访问内存或mmio寄存器
	randNum   *rand.Rand
}

func NewProcessor() Processor {
	source := rand.NewSource(time.Now().UnixMilli())
	return Processor{randNum: rand.New(source), bus: bus.NewBusWithNoROM()}
}

func NewProcessorWithROM(bus *bus.Bus) Processor {
	source := rand.NewSource(time.Now().UnixMilli())
	return Processor{randNum: rand.New(source), bus: bus}
}

func (p *Processor) LoadAndRun(program []byte) {
	p.loadProgram(program)
	p.reset()
	p.run()
}

func (p *Processor) LoadAndRunWithCallback(prevExec, afterExec CallbackFunc) {
	p.reset()
	p.runWithCallback(prevExec, afterExec)
}

func (p *Processor) loadProgram(program []byte) {
	p.bus.WriteRAM(ProgramEntryPoint, program)
	p.writeMemUint16(0xFFFC, ProgramEntryPoint)
}

func (p *Processor) reset() {
	p.regX = 0
	p.regA = 0
	p.regY = 0
	p.regStatus = 0
	p.sp = StackReset
	// 从ROM读取程序的entry point
	p.pc = p.readMemUint16(PrgROMEntryPointAddr)
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
			p.regStatus |= BreakStatus
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

func (p *Processor) runWithCallback(prevExec, afterExec CallbackFunc) {
	for {
		if p.bus.PollNMIInterrupt() {
			p.HandleInterrupt()
		}
		if !prevExec(p) {
			break
		}
		p.writeMemUint8(RandomNumber, byte(2+p.randNum.Intn(13)))
		opCode := p.readMemUint8(p.pc)
		p.pc++
		originalPc := p.pc
		instruction, ok := Instructions[opCode]
		if !ok {
			panic(fmt.Errorf("unknown instruction at %d: 0x%x", originalPc-1-PrgROMAddr, opCode))
		}
		switch opCode {
		case BRK:
			p.regStatus |= BreakStatus
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
		afterExec(p)
		p.bus.Tick(uint64(instruction.cycle))
	}
}

func (p *Processor) HandleKeyboardEvent(event *sdl.KeyboardEvent) {
	var action byte
	switch event.Keysym.Scancode {
	case sdl.SCANCODE_W:
		action = ActionUp
	case sdl.SCANCODE_S:
		action = ActionDown
	case sdl.SCANCODE_A:
		action = ActionLeft
	case sdl.SCANCODE_D:
		action = ActionRight
	default:
		return
	}
	p.writeMemUint8(Input, action)
}

func (p *Processor) HandleInterrupt() {
	// 保存PC
	p.stackPush(byte(p.pc & 0xff))
	p.stackPush(byte(p.pc >> 8))
	status := p.regStatus
	status &= (^BreakStatus)
	//status |= Break2Status
	// 保存状态，中断关闭
	p.stackPush(status)
	p.regStatus |= InterruptDisableStatus

	p.bus.Tick(2)
	// 跳转到中断处理
	p.pc = p.readMemUint16(PrgROMInterruptHandler)
	log.Printf("interrupt handle: 0x%x\n", p.pc-PrgROMAddr)
}

func (p *Processor) readMemUint8(addr uint16) byte {
	return p.bus.ReadMemUint8(addr)
}

func (p *Processor) writeMemUint8(addr uint16, val byte) {
	p.bus.WriteMemUint8(addr, val)
}

// 小端序读取16bits内存
func (p *Processor) readMemUint16(addr uint16) uint16 {
	low := p.readMemUint8(addr)
	high := p.readMemUint8(addr + 1)
	return uint16(high)<<8 + uint16(low)
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
	return p.bus.GetRAMRange(start, end)
}

func clc(p *Processor, _ Instruction) {
	p.regStatus &= ^CarryStatus
}

func cld(p *Processor, _ Instruction) {
	p.regStatus &= ^DecimalModeStatus
}

func cli(p *Processor, _ Instruction) {
	p.regStatus &= ^InterruptDisableStatus
}

func clv(p *Processor, _ Instruction) {
	p.regStatus &= ^OverflowStatus
}

func sec(p *Processor, _ Instruction) {
	p.regStatus |= CarryStatus
}

func sed(p *Processor, _ Instruction) {
	p.regStatus |= DecimalModeStatus
}
func sei(p *Processor, _ Instruction) {
	p.regStatus |= InterruptDisableStatus
}
