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

func jmpOffset(p *Processor) {
	addr := p.getMemoryAddress(Immediate)
	offset := p.readMemUint8(addr)
	// offset为负，pc减小
	if offset&(1<<7) != 0 {
		rev := int8(offset)
		p.pc = p.pc + 1 - uint16(-int16(rev))
	} else {
		p.pc = p.pc + 1 + uint16(offset)
	}
}

func bcc(p *Processor, _ Instruction) {
	if p.regStatus&CarryStatus == 0 {
		jmpOffset(p)
	}
}

func bcs(p *Processor, _ Instruction) {
	if p.regStatus&CarryStatus != 0 {
		jmpOffset(p)
	}
}

func beq(p *Processor, _ Instruction) {
	if p.regStatus&ZeroStatus != 0 {
		jmpOffset(p)
	}
}

func bne(p *Processor, _ Instruction) {
	if p.regStatus&ZeroStatus == 0 {
		jmpOffset(p)
	}
}

func bmi(p *Processor, _ Instruction) {
	if p.regStatus&NegativeStatus != 0 {
		jmpOffset(p)
	}
}

func bpl(p *Processor, _ Instruction) {
	if p.regStatus&NegativeStatus == 0 {
		jmpOffset(p)
	}
}

func bvc(p *Processor, _ Instruction) {
	if p.regStatus&OverflowStatus == 0 {
		jmpOffset(p)
	}
}

func bvs(p *Processor, _ Instruction) {
	if p.regStatus&OverflowStatus != 0 {
		jmpOffset(p)
	}
}
