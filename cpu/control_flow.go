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
