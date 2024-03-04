package cpu

func and(p *Processor, op *Instruction) {
	addr, cross := p.getMemoryAddress(op.AddrMode)
	p.regA = p.regA & p.readMemUint8(addr)
	p.zeroOrNegativeStatus(p.regA)
	if cross {
		p.bus.Tick(1)
	}
}

func eor(p *Processor, op *Instruction) {
	addr, cross := p.getMemoryAddress(op.AddrMode)
	eorWithA(p, p.readMemUint8(addr))
	if cross {
		p.bus.Tick(1)
	}
}

func ora(p *Processor, op *Instruction) {
	addr, cross := p.getMemoryAddress(op.AddrMode)
	oraWithA(p, p.readMemUint8(addr))
	if cross {
		p.bus.Tick(1)
	}
}

func bit(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
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

func oraWithA(p *Processor, val byte) {
	p.regA = p.regA | val
	p.zeroOrNegativeStatus(p.regA)
}

func andWithA(p *Processor, val byte) {
	p.regA = p.regA & val
	p.zeroOrNegativeStatus(p.regA)
}

func sre(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	val = lsrVal(p, val)
	p.writeMemUint8(addr, val)
	eorWithA(p, val)
}

func slo(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	val = aslVal(p, val)
	p.writeMemUint8(addr, val)
	oraWithA(p, val)
}

func rla(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	val = rolVal(p, val)
	p.writeMemUint8(addr, val)
	andWithA(p, val)
}

func rra(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	val = rorVal(p, val)
	p.writeMemUint8(addr, val)
	addRegA(p, val)
}

func xaa(p *Processor, op *Instruction) {
	p.regA = p.regX
	p.zeroOrNegativeStatus(p.regA)
	addr, _ := p.getMemoryAddress(op.AddrMode)
	andWithA(p, p.readMemUint8(addr))
}

func anc(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	andWithA(p, val)
	if p.regStatus&NegativeStatus != 0 {
		p.regStatus |= CarryStatus
	} else {
		p.regStatus &= ^CarryStatus
	}
}

func alr(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	andWithA(p, val)
	p.regA = lsrVal(p, p.regA)
}

func arr(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	andWithA(p, val)
	p.regA = rorVal(p, p.regA)
}

func axs(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	ax := p.regX & p.regA
	res := ax - val
	if val <= res {
		p.regStatus |= CarryStatus
	}
	p.zeroOrNegativeStatus(res)
	p.regX = res
}
