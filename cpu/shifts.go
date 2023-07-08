package cpu

func asl(p *Processor, op Instruction) {
	var val byte
	var addr uint16
	if op.addrMode == NoneAddressing {
		val = p.regA
	} else {
		addr := p.getMemoryAddress(op.addrMode)
		val = p.readMemUint8(addr)
	}
	if val&(1<<7) != 0 {
		p.regStatus |= CarryStatus
	} else {
		p.regStatus &= ^CarryStatus
	}
	val = val << 1
	p.zeroOrNegativeStatus(val)

	if op.addrMode == NoneAddressing {
		p.regA = val
	} else {
		p.writeMemUint8(addr, val)
	}
}

func lsr(p *Processor, op Instruction) {
	var val byte
	var addr uint16
	if op.addrMode == NoneAddressing {
		val = p.regA
	} else {
		addr := p.getMemoryAddress(op.addrMode)
		val = p.readMemUint8(addr)
	}
	if val&1 != 0 {
		p.regStatus |= CarryStatus
	} else {
		p.regStatus &= ^CarryStatus
	}
	val = val >> 1
	p.zeroOrNegativeStatus(val)

	if op.addrMode == NoneAddressing {
		p.regA = val
	} else {
		p.writeMemUint8(addr, val)
	}
}
