package cpu

func and(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	p.regA = p.regA & p.readMemUint8(addr)
	p.zeroOrNegativeStatus(p.regA)
}

func eor(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	eorWithA(p, p.readMemUint8(addr))
}

func ora(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	p.regA = p.regA | p.readMemUint8(addr)
	p.zeroOrNegativeStatus(p.regA)
}

func bit(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	res := p.regA & val
	if res == 0 {
		p.regStatus |= ZeroStatus
	} else {
		p.regStatus &= ^ZeroStatus
	}
	if val&(1<<6) != 0 {
		p.regStatus |= OverflowStatus
	} else {
		p.regStatus &= ^OverflowStatus
	}
	if val&(1<<7) != 0 {
		p.regStatus |= NegativeStatus
	} else {
		p.regStatus &= ^NegativeStatus
	}
}

func eorWithA(p *Processor, val byte) {
	p.regA = p.regA ^ val
	p.zeroOrNegativeStatus(p.regA)
}

func sre(p *Processor, op Instruction) {
	var val byte
	if op.AddrMode == NoneAddressing {
		val = p.regA
	} else {
		addr := p.getMemoryAddress(op.AddrMode)
		val = p.readMemUint8(addr)
	}
	val = lsrVal(p, val)
	eorWithA(p, val)
}
