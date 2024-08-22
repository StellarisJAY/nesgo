package cpu

func jmp(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	p.pc = addr
}

func jmpIndirect(p *Processor, _ *Instruction) {
	// 先从参数取到地址
	addr := p.readMemUint16(p.pc)
	var target uint16
	if addr&0xFF == 0xFF {
		low := p.readMemUint8(addr)
		high := p.readMemUint8(addr & 0xFF00)
		target = uint16(high)<<8 + uint16(low)
	} else {
		target = p.readMemUint16(addr)
	}
	p.pc = target
}

// jsr 返回地址入栈，跳转到目标地址
func jsr(p *Processor, _ *Instruction) {
	// 将返回地址：pc+2-1 （跳过16位立即数）入栈
	ra := p.pc + 2 - 1
	p.stackPush(byte(ra >> 8))
	p.stackPush(byte(ra & 0xFF))
	// 跳转到目标地址
	target, _ := p.getMemoryAddress(Absolute)
	p.pc = target
}

// rts，返回地址出栈，跳转到返回地址
func rts(p *Processor, _ *Instruction) {
	low := p.stackPop()
	high := p.stackPop()
	ra := uint16(high)<<8 | uint16(low)
	p.pc = ra + 1
}

func rti(p *Processor, _ *Instruction) {
	status := p.stackPop()
	p.regStatus = status
	p.regStatus &= ^BreakStatus
	p.regStatus |= Break2Status
	low := p.stackPop()
	high := p.stackPop()
	ra := uint16(high)<<8 | uint16(low)
	p.pc = ra
}

func branch(p *Processor, condition bool) {
	if condition {
		// branch 成功会+1 cycle
		p.bus.Tick(1)
		addr, _ := p.getMemoryAddress(Immediate)
		offset := int8(p.readMemUint8(addr))
		target := p.pc + 1 + uint16(offset)
		if pageCross(p.pc+1, target) {
			p.bus.Tick(1)
		}
		p.pc = target
	}
}

func bcc(p *Processor, _ *Instruction) {
	branch(p, p.regStatus&CarryStatus == 0)
}

func bcs(p *Processor, _ *Instruction) {
	branch(p, p.regStatus&CarryStatus != 0)
}

func beq(p *Processor, _ *Instruction) {
	branch(p, p.regStatus&ZeroStatus != 0)
}

func bne(p *Processor, _ *Instruction) {
	branch(p, p.regStatus&ZeroStatus == 0)
}

func bmi(p *Processor, _ *Instruction) {
	branch(p, p.regStatus&NegativeStatus != 0)
}

func bpl(p *Processor, _ *Instruction) {
	branch(p, p.regStatus&NegativeStatus == 0)
}

func bvc(p *Processor, _ *Instruction) {
	branch(p, p.regStatus&OverflowStatus == 0)
}

func bvs(p *Processor, _ *Instruction) {
	branch(p, p.regStatus&OverflowStatus != 0)
}
