package cpu

func tsx(p *Processor, _ Instruction) {
	p.regX = p.sp
	p.zeroOrNegativeStatus(p.regX)
}

func txs(p *Processor, _ Instruction) {
	p.sp = p.regX
	p.zeroOrNegativeStatus(p.sp)
}

func pha(p *Processor, _ Instruction) {
	p.stackPush(p.regA)
}

func php(p *Processor, _ Instruction) {
	p.stackPush(p.regStatus)
}

func pla(p *Processor, _ Instruction) {
	p.regA = p.stackPop()
	p.zeroOrNegativeStatus(p.regA)
}

func plp(p *Processor, _ Instruction) {
	p.regStatus = p.stackPop()
}

// stackPush 入栈一个byte，栈指针减小
func (p *Processor) stackPush(val byte) {
	p.sp -= 1
	p.writeMemUint8(uint16(p.sp)+StackBase, val)
}

// stackPop 弹出一个byte，栈指针增大
func (p *Processor) stackPop() byte {
	val := p.readMemUint8(uint16(p.sp) + StackBase)
	p.sp += 1
	return val
}

func (p *Processor) peek() byte {
	return p.readMemUint8(uint16(p.sp) + StackBase)
}
