package cpu

func (p *Processor) inx() {
	p.regX = wrappingAddOne(p.regX)
	p.zeroOrNegativeStatus(p.regX)
}

func (p *Processor) iny() {
	p.regY = wrappingAddOne(p.regY)
	p.zeroOrNegativeStatus(p.regY)
}

func inc(p *Processor, op *Instruction) {
	addr, _ := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	value := wrappingAddOne(val)
	p.zeroOrNegativeStatus(value)
	p.writeMemUint8(addr, value)
}

func wrappingAddOne(val byte) byte {
	if val == 0xff {
		return 0
	} else {
		return val + 1
	}
}

func wrappingMinusOne(val byte) byte {
	if val == 0 {
		val = 0xff
	} else {
		val = val - 1
	}
	return val
}

func dex(p *Processor, _ *Instruction) {
	p.regX = wrappingMinusOne(p.regX)
	p.zeroOrNegativeStatus(p.regX)
}

func dey(p *Processor, _ *Instruction) {
	p.regY = wrappingMinusOne(p.regY)
	p.zeroOrNegativeStatus(p.regY)
}

func dec(p *Processor, op *Instruction) {
	addr, cross := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	val = wrappingMinusOne(val)
	p.zeroOrNegativeStatus(val)
	p.writeMemUint8(addr, val)
	if cross {
		p.bus.Tick(1)
	}
}
