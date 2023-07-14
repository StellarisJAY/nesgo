package cpu

func lda(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	p.regA = val
	p.zeroOrNegativeStatus(p.regA)
}

func ldx(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	p.regX = val
	p.zeroOrNegativeStatus(p.regX)
}

func ldy(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	p.regY = val
	p.zeroOrNegativeStatus(p.regY)
}

func sta(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	p.writeMemUint8(addr, p.regA)
}

func stx(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	p.writeMemUint8(addr, p.regX)
}

func sty(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	p.writeMemUint8(addr, p.regY)
}

func tax(p *Processor, _ Instruction) {
	p.regX = p.regA
	p.zeroOrNegativeStatus(p.regX)
}

func tay(p *Processor, _ Instruction) {
	p.regY = p.regA
	p.zeroOrNegativeStatus(p.regY)
}

func txa(p *Processor, _ Instruction) {
	p.regA = p.regX
	p.zeroOrNegativeStatus(p.regA)
}

func tya(p *Processor, _ Instruction) {
	p.regA = p.regY
	p.zeroOrNegativeStatus(p.regA)
}

func lax(p *Processor, op Instruction) {
	p.regA = p.readMemUint8(p.getMemoryAddress(op.AddrMode))
	p.zeroOrNegativeStatus(p.regA)
	p.regX = p.regA
}

func sax(p *Processor, op Instruction) {
	res := p.regA & p.regX
	p.writeMemUint8(p.getMemoryAddress(op.AddrMode), res)
}
