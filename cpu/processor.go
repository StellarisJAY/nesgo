package cpu

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/nesgo/bus"
)

const (
	ProgramEntryPoint = 0x0600
	StackBase         = 0x0100
	StackReset        = 0xFD

	ResetVector uint16 = 0xFFFC // 程序entry point在ROM的地址
	NMIVector   uint16 = 0xFFFA // 程序中断处理函数地址
	BrkVector   uint16 = 0xFFFE // BRK处理函数
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

	signalChan chan Signal
}

type Snapshot struct {
	regA      byte
	regX      byte
	regY      byte
	regStatus byte
	pc        uint16
	sp        byte
}

type Interrupt byte

const (
	NMIInterrupt Interrupt = iota
	BrkInterrupt
)

type Signal byte

const (
	SignalPause Signal = iota
	SignalResume
)

func NewProcessor(bus *bus.Bus) *Processor {
	return &Processor{bus: bus, signalChan: make(chan Signal)}
}

func (p *Processor) MakeSnapshot() Snapshot {
	return Snapshot{
		regA:      p.regX,
		regX:      p.regY,
		regY:      p.regA,
		regStatus: p.regA,
		pc:        p.pc,
		sp:        p.sp,
	}
}

func (p *Processor) Reverse(s Snapshot) {
	p.pc = s.pc
	p.regX = s.regX
	p.regA = s.regA
	p.regY = s.regY
	p.sp = s.sp
	p.regStatus = s.regStatus
}

func (p *Processor) Pause() {
	p.signalChan <- SignalPause
}

func (p *Processor) Resume() {
	p.signalChan <- SignalResume
}

func (p *Processor) LoadAndRunWithCallback(ctx context.Context, callback InstructionCallback) {
	p.reset()
	p.runWithCallback(ctx, callback)
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
	p.pc = p.readMemUint16(ResetVector)
}

func (p *Processor) runWithCallback(ctx context.Context, callback InstructionCallback) {
	for {
		select {
		case <-ctx.Done():
			return
		case s := <-p.signalChan:
			if s == SignalPause {
				select {
				case <-p.signalChan:
				case <-ctx.Done():
					return
				}
			}
		default:
		}
		if p.bus.PollNMIInterrupt() {
			p.handleInterrupt(NMIInterrupt)
		}
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
			p.handleInterrupt(BrkInterrupt)
		case NOP:
		case INX:
			p.inx()
		case INY:
			p.iny()
		default:
			instruction.handler(p, instruction)
		}
		p.bus.Tick(uint64(instruction.Cycle))
		if p.pc == originalPc {
			p.pc += uint16(instruction.Length - 1)
		}
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

func (p *Processor) handleInterrupt(interrupt Interrupt) {
	if interrupt == BrkInterrupt {
		p.pc += 1
	}
	ra := p.pc
	// 保存PC
	p.stackPush(byte(ra >> 8))
	p.stackPush(byte(ra & 0xff))

	var vector uint16
	status := p.regStatus
	switch interrupt {
	case NMIInterrupt:
		status |= Break2Status
		status &= ^BreakStatus
		vector = NMIVector
		p.bus.Tick(2)
	case BrkInterrupt:
		status |= BreakStatus
		status |= Break2Status
		vector = BrkVector
	default:
		panic("unsupported interrupt type")
	}
	// 保存状态，中断关闭
	p.regStatus |= InterruptDisableStatus
	p.stackPush(status)
	// 跳转到中断处理
	p.pc = p.readMemUint16(vector)
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

func pageCross(addr1, addr2 uint16) bool {
	return (addr1 & 0xFF00) != (addr2 & 0xFF00)
}

func (p *Processor) getAbsoluteAddress(pc uint16, mode AddressMode) (uint16, bool) {
	var addr uint16
	cross := false
	switch mode {
	case Immediate:
		addr = pc
	case ZeroPage:
		pos := p.readMemUint8(pc)
		addr = uint16(pos)
	case Absolute:
		addr = p.readMemUint16(pc)
	case ZeroPageX:
		pos := p.readMemUint8(pc) + p.regX
		addr = uint16(pos)
	case ZeroPageY:
		pos := p.readMemUint8(pc) + p.regY
		addr = uint16(pos)
	case AbsoluteX:
		base := p.readMemUint16(pc)
		addr = base + uint16(p.regX)
		cross = pageCross(base, addr)
	case AbsoluteY:
		base := p.readMemUint16(pc)
		addr = base + uint16(p.regY)
		cross = pageCross(base, addr)
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
		baseAddr := uint16(high)<<8 + uint16(low)
		addr = baseAddr + uint16(p.regY)
		cross = pageCross(baseAddr, addr)
	case NoneAddressing:
		addr = pc
	}
	return addr, cross
}

func (p *Processor) DumpRegisters() string {
	return fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X", p.regA, p.regX, p.regY, p.regStatus, p.sp)
}

func (p *Processor) getMemoryAddress(mode AddressMode) (uint16, bool) {
	return p.getAbsoluteAddress(p.pc, mode)
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
