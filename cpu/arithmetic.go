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

func adc(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	val := p.readMemUint8(addr)
	addRegA(p, val)
}

func addRegA(p *Processor, val byte) {
	carry := uint16(p.regStatus&CarryStatus) >> 1
	res16 := uint16(p.regA) + uint16(val) + carry
	if res16 > 0xff {
		p.regStatus = p.regStatus | CarryStatus
	} else {
		p.regStatus = p.regStatus & (^CarryStatus)
	}
	res8 := byte(res16 & 0xFF)
	if (val^res8)&(p.regA^res8)&0x80 != 0 {
		p.regStatus |= OverflowStatus
	} else {
		p.regStatus &= ^OverflowStatus
	}
	p.regA = res8
	p.zeroOrNegativeStatus(p.regA)
}

func sbc(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	val := p.readMemUint8(addr)
	val = byte(int8(-val) - 1)
	addRegA(p, val)
}
