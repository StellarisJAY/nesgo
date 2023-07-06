package cpu

func jmp(p *Processor, op Instruction) {
	addr := p.getMemoryAddress(op.addrMode)
	p.pc = p.readMemUint16(addr)
}
