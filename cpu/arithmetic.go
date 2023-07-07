package cpu

func cmp(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	val := p.readMemUint8(addr)
	if p.regA >= val {
		p.regStatus |= CarryStatus
		if val == p.regA {
			p.regStatus |= ZeroStatus
		}
	}
}

func cpx(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	val := p.readMemUint8(addr)
	if val <= p.regX {
		p.regStatus |= CarryStatus
		if val == p.regX {
			p.regStatus |= ZeroStatus
		}
	}
}

func cpy(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	val := p.readMemUint8(addr)
	if val <= p.regY {
		p.regStatus |= CarryStatus
		if val == p.regY {
			p.regStatus |= ZeroStatus
		}
	}
}
