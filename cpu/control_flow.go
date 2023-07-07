package cpu

func jmp(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	p.pc = addr
}

func jmpIndirect(p *Processor, _ Instruction) {
	// 先从参数取到地址
	addr := p.getMemoryAddress(Absolute)
	// 在地址取到跳转目标
	target := p.readMemUint8(addr)
	p.pc = uint16(target)
}

// jsr 返回地址入栈，跳转到目标地址
func jsr(p *Processor, _ Instruction) {
	// 将返回地址：pc+2 （跳过16位立即数）入栈
	ra := p.pc + 2
	p.stackPush(byte(ra & 0xFF))
	p.stackPush(byte(ra >> 8))
	// 跳转到目标地址
	target := p.getMemoryAddress(Absolute)
	p.pc = target
}

// rts，返回地址出栈，跳转到返回地址
func rts(p *Processor, _ Instruction) {
	high := p.stackPop()
	low := p.stackPop()
	ra := uint16(high)<<8 | uint16(low)
	p.pc = ra
}
