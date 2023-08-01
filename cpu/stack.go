package cpu

func tsx(p *Processor, _ *Instruction) {
	p.regX = p.sp
	p.zeroOrNegativeStatus(p.regX)
}

func txs(p *Processor, _ *Instruction) {
	p.sp = p.regX
}

func pha(p *Processor, _ *Instruction) {
	p.stackPush(p.regA)
}

func php(p *Processor, _ *Instruction) {
	p.regStatus |= Break2Status
	p.regStatus |= BreakStatus
	p.stackPush(p.regStatus)
}

func pla(p *Processor, _ *Instruction) {
	p.regA = p.stackPop()
	p.zeroOrNegativeStatus(p.regA)
}

func plp(p *Processor, _ *Instruction) {
	p.regStatus = p.stackPop()
}

// stackPush 入栈一个byte，栈指针减小
func (p *Processor) stackPush(val byte) {
	p.writeMemUint8(uint16(p.sp)+StackBase, val)
	p.sp -= 1
}

// stackPop 弹出一个byte，栈指针增大
func (p *Processor) stackPop() byte {
	p.sp += 1
	val := p.readMemUint8(uint16(p.sp) + StackBase)
	return val
}

func (p *Processor) peek() byte {
	return p.readMemUint8(uint16(p.sp) + StackBase)
}
