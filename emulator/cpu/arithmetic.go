package cpu

func compare(p *Processor, val byte, compareWith byte) {
	if val <= compareWith {
		p.regStatus |= CarryStatus
	} else {
		p.regStatus &= ^CarryStatus
	}
	p.zeroOrNegativeStatus(compareWith - val)
}

func cmp(p *Processor, op *Instruction) {
	addr, cross := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	compare(p, val, p.regA)
	if cross {
		p.bus.Tick(1)
	}
}

func cpx(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	compare(p, val, p.regX)
}

func cpy(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	compare(p, val, p.regY)
}

func adc(p *Processor, op *Instruction) {
	addr, cross := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	addRegA(p, val)
	if cross {
		p.bus.Tick(1)
	}
}

func addRegA(p *Processor, val byte) {
	var carry uint16
	if p.regStatus&CarryStatus != 0 {
		carry = 1
	} else {
		carry = 0
	}
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

func sbc(p *Processor, op *Instruction) {
	addr, cross := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	val = byte(int8(-val) - 1)
	addRegA(p, val)
	if cross {
		p.bus.Tick(1)
	}
}

func subRegA(p *Processor, val byte) {
	delta := byte(int8(-val) - 1)
	addRegA(p, delta)
}

func dcp(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	val -= 1
	p.writeMemUint8(addr, val)
	if val <= p.regA {
		p.regStatus |= CarryStatus
	}
	p.zeroOrNegativeStatus(p.regA - val)
}

func isc(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	val += 1
	p.writeMemUint8(addr, val)
	subRegA(p, val)
}
