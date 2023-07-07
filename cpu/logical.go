package cpu

func and(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	p.regA = p.regA & p.readMemUint8(addr)
	p.zeroOrNegativeStatus(p.regA)
}

func eor(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	p.regA = p.regA ^ p.readMemUint8(addr)
	p.zeroOrNegativeStatus(p.regA)
}

func ora(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	p.regA = p.regA | p.readMemUint8(addr)
	p.zeroOrNegativeStatus(p.regA)
}

func bit(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	val := p.regA & p.readMemUint8(addr)
	p.zeroOrNegativeStatus(val)
}
