package cpu

func asl(p *Processor, op Instruction) {
	var val byte
	var addr uint16
	if op.AddrMode == NoneAddressing {
		val = p.regA
	} else {
		addr := p.getMemoryAddress(op.AddrMode)
		val = p.readMemUint8(addr)
	}
	if val&(1<<7) != 0 {
		p.regStatus |= CarryStatus
	} else {
		p.regStatus &= ^CarryStatus
	}
	val = val << 1
	p.zeroOrNegativeStatus(val)

	if op.AddrMode == NoneAddressing {
		p.regA = val
	} else {
		p.writeMemUint8(addr, val)
	}
}

func lsr(p *Processor, op Instruction) {
	var val byte
	var addr uint16
	if op.AddrMode == NoneAddressing {
		val = p.regA
	} else {
		addr := p.getMemoryAddress(op.AddrMode)
		val = p.readMemUint8(addr)
	}
	if val&1 != 0 {
		p.regStatus |= CarryStatus
	} else {
		p.regStatus &= ^CarryStatus
	}
	val = val >> 1
	p.zeroOrNegativeStatus(val)

	if op.AddrMode == NoneAddressing {
		p.regA = val
	} else {
		p.writeMemUint8(addr, val)
	}
}

func ror(p *Processor, op Instruction) {
	var val byte
	var addr uint16
	if op.AddrMode == NoneAddressing {
		val = p.regA
	} else {
		addr := p.getMemoryAddress(op.AddrMode)
		val = p.readMemUint8(addr)
	}
	oldCarry := p.regStatus&CarryStatus != 0
	if val&1 != 0 {
		p.regStatus = p.regStatus | CarryStatus
	} else {
		p.regStatus = p.regStatus & (^CarryStatus)
	}
	val = val >> 1
	if oldCarry {
		val = val & 0b10000000
	}
	p.zeroOrNegativeStatus(val)
	if op.AddrMode == NoneAddressing {
		p.regA = val
	} else {
		p.writeMemUint8(addr, val)
	}
}

func rol(p *Processor, op Instruction) {
	var val byte
	var addr uint16
	if op.AddrMode == NoneAddressing {
		val = p.regA
	} else {
		addr := p.getMemoryAddress(op.AddrMode)
		val = p.readMemUint8(addr)
	}
	oldCarry := p.regStatus&CarryStatus != 0
	if val&(1<<7) != 0 {
		p.regStatus = p.regStatus | CarryStatus
	} else {
		p.regStatus = p.regStatus & (^CarryStatus)
	}
	val = val << 1
	if oldCarry {
		val = val & 0b1
	}
	p.zeroOrNegativeStatus(val)
	if op.AddrMode == NoneAddressing {
		p.regA = val
	} else {
		p.writeMemUint8(addr, val)
	}
}
