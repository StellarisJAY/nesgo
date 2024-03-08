package cpu

func asl(p *Processor, op *Instruction) {
	var val byte
	var addr uint16
	var cross bool
	if op.AddrMode == NoneAddressing {
		val = p.regA
	} else {
		addr, cross = p.getMemoryAddress(op.AddrMode)
		val = p.readMemUint8(addr)
	}
	val = aslVal(p, val)

	if op.AddrMode == NoneAddressing {
		p.regA = val
	} else {
		p.writeMemUint8(addr, val)
	}
	if cross {
		p.bus.Tick(1)
	}
}

func lsr(p *Processor, op *Instruction) {
	var val byte
	var addr uint16
	var cross bool
	if op.AddrMode == NoneAddressing {
		val = p.regA
	} else {
		addr, cross = p.getMemoryAddress(op.AddrMode)
		val = p.readMemUint8(addr)
	}
	val = lsrVal(p, val)

	if op.AddrMode == NoneAddressing {
		p.regA = val
	} else {
		p.writeMemUint8(addr, val)
	}
	if cross {
		p.bus.Tick(1)
	}
}

func aslVal(p *Processor, val byte) byte {
	if val&(1<<7) != 0 {
		p.regStatus |= CarryStatus
	} else {
		p.regStatus &= ^CarryStatus
	}
	val = val << 1
	p.zeroOrNegativeStatus(val)
	return val
}

func lsrVal(p *Processor, val byte) byte {
	if val&1 != 0 {
		p.regStatus |= CarryStatus
	} else {
		p.regStatus &= ^CarryStatus
	}
	val = val >> 1
	p.zeroOrNegativeStatus(val)
	return val
}

func rolVal(p *Processor, val byte) byte {
	oldCarry := p.regStatus&CarryStatus != 0
	if val&(1<<7) != 0 {
		p.regStatus = p.regStatus | CarryStatus
	} else {
		p.regStatus = p.regStatus & (^CarryStatus)
	}
	val = val << 1
	if oldCarry {
		val = val | 0b1
	}
	p.zeroOrNegativeStatus(val)
	return val
}

func rorVal(p *Processor, val byte) byte {
	oldCarry := p.regStatus&CarryStatus != 0
	if val&1 != 0 {
		p.regStatus = p.regStatus | CarryStatus
	} else {
		p.regStatus = p.regStatus & (^CarryStatus)
	}
	val = val >> 1
	if oldCarry {
		val = val | 0b10000000
	}
	p.zeroOrNegativeStatus(val)
	return val
}

func ror(p *Processor, op *Instruction) {
	var val byte
	var addr uint16
	var cross bool
	if op.AddrMode == NoneAddressing {
		val = p.regA
	} else {
		addr, cross = p.getMemoryAddress(op.AddrMode)
		val = p.readMemUint8(addr)
	}
	val = rorVal(p, val)
	if op.AddrMode == NoneAddressing {
		p.regA = val
	} else {
		p.writeMemUint8(addr, val)
	}
	if cross {
		p.bus.Tick(1)
	}
}

func rol(p *Processor, op *Instruction) {
	var val byte
	var addr uint16
	var cross bool
	if op.AddrMode == NoneAddressing {
		val = p.regA
	} else {
		addr, cross = p.getMemoryAddress(op.AddrMode)
		val = p.readMemUint8(addr)
	}
	val = rolVal(p, val)
	if op.AddrMode == NoneAddressing {
		p.regA = val
	} else {
		p.writeMemUint8(addr, val)
	}
	if cross {
		p.bus.Tick(1)
	}
}
