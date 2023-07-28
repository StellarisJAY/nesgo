package cpu

func lda(p *Processor, op *Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	p.regA = val
	p.zeroOrNegativeStatus(p.regA)
}

func ldx(p *Processor, op *Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	p.regX = val
	p.zeroOrNegativeStatus(p.regX)
}

func ldy(p *Processor, op *Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	val := p.readMemUint8(addr)
	p.regY = val
	p.zeroOrNegativeStatus(p.regY)
}

func sta(p *Processor, op *Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	p.writeMemUint8(addr, p.regA)
}

func stx(p *Processor, op *Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	p.writeMemUint8(addr, p.regX)
}

func sty(p *Processor, op *Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	p.writeMemUint8(addr, p.regY)
}

func tax(p *Processor, _ *Instruction) {
	p.regX = p.regA
	p.zeroOrNegativeStatus(p.regX)
}

func tay(p *Processor, _ *Instruction) {
	p.regY = p.regA
	p.zeroOrNegativeStatus(p.regY)
}

func txa(p *Processor, _ *Instruction) {
	p.regA = p.regX
	p.zeroOrNegativeStatus(p.regA)
}

func tya(p *Processor, _ *Instruction) {
	p.regA = p.regY
	p.zeroOrNegativeStatus(p.regA)
}

func lax(p *Processor, op *Instruction) {
	p.regA = p.readMemUint8(p.getMemoryAddress(op.AddrMode))
	p.zeroOrNegativeStatus(p.regA)
	p.regX = p.regA
}

func lxa(p *Processor, op *Instruction) {
	lda(p, op)
	tax(p, op)
}

func sax(p *Processor, op *Instruction) {
	res := p.regA & p.regX
	p.writeMemUint8(p.getMemoryAddress(op.AddrMode), res)
}

func las(p *Processor, op *Instruction) {
	val := p.readMemUint8(p.getMemoryAddress(op.AddrMode))
	val = val & p.sp
	p.regA, p.regX, p.sp = val, val, val
	p.zeroOrNegativeStatus(val)
}
func tas(p *Processor, op *Instruction) {
	val := p.regA & p.regX
	p.sp = val
	addr := p.readMemUint16(p.pc) + uint16(p.regY)
	data := (byte(addr>>8) + 1) & p.sp
	p.writeMemUint8(addr, data)
}

func ahx(p *Processor, op *Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	data := p.regA & p.regX & byte(addr>>8)
	p.writeMemUint8(addr, data)
}

func shx(p *Processor, op *Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	data := p.regX & (byte(addr>>8) + 1)
	p.writeMemUint8(addr, data)
}

func shy(p *Processor, op *Instruction) {
	addr := p.getMemoryAddress(op.AddrMode)
	data := p.regY & (byte(addr>>8) + 1)
	p.writeMemUint8(addr, data)
}
