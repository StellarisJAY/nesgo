package cpu

func (p *Processor) inx() {
	p.regX = wrappingAddOne(p.regX)
	p.zeroOrNegativeStatus(p.regX)
}

func (p *Processor) iny() {
	p.regY = wrappingAddOne(p.regY)
	p.zeroOrNegativeStatus(p.regY)
}

func inc(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	value := p.readMemUint8(addr)
	p.writeMemUint8(addr, wrappingAddOne(value))
}

func wrappingAddOne(val byte) byte {
	if val == 0xff {
		return 0
	} else {
		return val + 1
	}
}

func dex(p *Processor, _ Instruction) {
	p.regX = p.regX - 1
	p.zeroOrNegativeStatus(p.regX)
}

func dey(p *Processor, _ Instruction) {
	p.regY = p.regY - 1
	p.zeroOrNegativeStatus(p.regY)
}
