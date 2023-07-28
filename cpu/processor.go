package cpu

import (
	"fmt"
	"github.com/stellarisJAY/nesgo/bus"
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
	StackBase                = 0x0100
	StackReset               = 0xFD
	RandomNumber      uint16 = 0xFE

	PrgROMEntryPointAddr   uint16 = 0xFFFC // 程序entry point在ROM的地址
	PrgROMInterruptHandler uint16 = 0xFFFA // 程序中断处理函数地址
)

const (
	CarryStatus            byte = 1 << 0
	ZeroStatus             byte = 1 << 1
	InterruptDisableStatus byte = 1 << 2
	DecimalModeStatus      byte = 1 << 3
	BreakStatus            byte = 1 << 4
	Break2Status           byte = 1 << 5
	OverflowStatus         byte = 1 << 6
	NegativeStatus         byte = 1 << 7
)

// CallbackFunc 每条指令执行前的callback，返回false将结束处理器循环
type CallbackFunc func(*Processor) bool
type InstructionCallback func(*Processor, *Instruction)

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

func NewProcessorWithROM(bus *bus.Bus) *Processor {
	source := rand.NewSource(time.Now().UnixMilli())
	return &Processor{randNum: rand.New(source), bus: bus}
}

func (p *Processor) LoadAndRun(program []byte) {
	p.loadProgram(program)
	p.reset()
	p.run()
}

func (p *Processor) LoadAndRunWithCallback(callback InstructionCallback) {
	p.reset()
	p.runWithCallback(callback)
}

func (p *Processor) loadProgram(program []byte) {
	p.bus.WriteRAM(ProgramEntryPoint, program)
	p.writeMemUint16(0xFFFC, ProgramEntryPoint)
}

func (p *Processor) reset() {
	p.regX = 0
	p.regA = 0
	p.regY = 0
	p.regStatus = 0b100100
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
			p.pc += uint16(instruction.Length - 1)
		}
	}
}

func (p *Processor) runWithCallback(callback InstructionCallback) {
	for {
		if p.bus.PollNMIInterrupt() {
			p.HandleInterrupt()
		}
		// 在0xFE保存0~255随机数
		p.writeMemUint8(RandomNumber, byte(p.randNum.Intn(256)))
		opCode := p.readMemUint8(p.pc)
		p.pc++
		originalPc := p.pc
		instruction, ok := Instructions[opCode]
		if !ok {
			panic(fmt.Errorf("unknown instruction at %04x: 0x%x", originalPc-1, opCode))
		}
		callback(p, instruction)
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
			p.pc += uint16(instruction.Length - 1)
		}
		p.bus.Tick(uint64(instruction.Cycle))
	}
}

func (p *Processor) Disassemble(callback InstructionCallback) {
	p.pc = 0x8000
	for p.pc >= 0x8000 {
		opCode := p.readMemUint8(p.pc)
		p.pc += 1
		instruction, ok := Instructions[opCode]
		if !ok {
			panic(fmt.Errorf("unknown instruction at %04x: 0x%x", p.pc-1, opCode))
		}
		callback(p, instruction)
		p.pc += uint16(instruction.Length) - 1
	}
}

func (p *Processor) HandleInterrupt() {
	ra := p.pc
	// 保存PC
	p.stackPush(byte(ra >> 8))
	p.stackPush(byte(ra & 0xff))
	status := p.regStatus
	status |= Break2Status
	// 保存状态，中断关闭
	p.stackPush(status)
	p.regStatus |= InterruptDisableStatus
	p.bus.Tick(2)
	// 跳转到中断处理
	p.pc = p.readMemUint16(PrgROMInterruptHandler)
	fmt.Printf("interrupt: %x\n", p.readMemUint16(PrgROMInterruptHandler))
}

func (p *Processor) GetArgAddress(pc uint16, mode AddressMode) uint16 {
	return p.getAddress(pc, mode)
}

func (p *Processor) ReadMem8(addr uint16) byte {
	return p.readMemUint8(addr)
}

func (p *Processor) ReadMem16(addr uint16) uint16 {
	return p.readMemUint16(addr)
}

func (p *Processor) ProgramCounter() uint16 {
	return p.pc
}

func (p *Processor) readMemUint8(addr uint16) byte {
	return p.bus.ReadMemUint8(addr)
}

func (p *Processor) writeMemUint8(addr uint16, val byte) {
	p.bus.WriteMemUint8(addr, val)
}

func (p *Processor) Cycles() uint64 {
	return p.bus.Cycles()
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

func (p *Processor) getAddress(pc uint16, mode AddressMode) uint16 {
	var addr uint16
	switch mode {
	case Immediate:
		addr = p.pc
	case ZeroPage:
		addr = uint16(p.readMemUint8(pc))
	case Absolute:
		addr = p.readMemUint16(p.pc)
	case ZeroPageX:
		addr = uint16(p.readMemUint8(pc) + p.regX)
	case ZeroPageY:
		addr = uint16(p.readMemUint8(pc) + p.regY)
	case AbsoluteX:
		addr = p.readMemUint16(pc)
		addr += uint16(p.regX)
	case AbsoluteY:
		addr = p.readMemUint16(pc)
		addr += uint16(p.regY)
	case IndirectX:
		base := p.readMemUint8(pc)
		ptr := base + p.regX
		// ptr必须是8位地址
		low := p.readMemUint8(uint16(ptr))
		high := p.readMemUint8(uint16(ptr + 1))
		addr = uint16(high)<<8 + uint16(low)
	case IndirectY:
		base := p.readMemUint8(pc)
		low := p.readMemUint8(uint16(base))
		high := p.readMemUint8(uint16(base + 1))
		addr = uint16(high)<<8 + uint16(low) + uint16(p.regY)
	case NoneAddressing:
		addr = pc
	}
	return addr
}

func (p *Processor) DumpRegisters() string {
	return fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X", p.regA, p.regX, p.regY, p.regStatus, p.sp)
}

func (p *Processor) getMemoryAddress(mode AddressMode) uint16 {
	return p.getAddress(p.pc, mode)
}

// GetMemoryRange 获取start到end范围内的内存切片
func (p *Processor) GetMemoryRange(start, end uint16) []byte {
	return p.bus.GetRAMRange(start, end)
}

func clc(p *Processor, _ *Instruction) {
	p.regStatus &= ^CarryStatus
}

func cld(p *Processor, _ *Instruction) {
	p.regStatus &= ^DecimalModeStatus
}

func cli(p *Processor, _ *Instruction) {
	p.regStatus &= ^InterruptDisableStatus
}

func clv(p *Processor, _ *Instruction) {
	p.regStatus &= ^OverflowStatus
}

func sec(p *Processor, _ *Instruction) {
	p.regStatus |= CarryStatus
}

func sed(p *Processor, _ *Instruction) {
	p.regStatus |= DecimalModeStatus
}
func sei(p *Processor, _ *Instruction) {
	p.regStatus |= InterruptDisableStatus
}
