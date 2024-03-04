package ppu

// AddrRegister ppu地址寄存器
// 通过修改0x2006寄存器来设置ppu地址
type AddrRegister struct {
	Low     byte
	High    byte
	HighPtr bool
}

func NewAddrRegister() AddrRegister {
	return AddrRegister{
		Low:     0,
		High:    0,
		HighPtr: true,
	}
}

func (ag *AddrRegister) set(addr uint16) {
	ag.High = byte(addr >> 8)
	ag.Low = byte(addr & 0xFF)
}

func (ag *AddrRegister) get() uint16 {
	return uint16(ag.High)<<8 + uint16(ag.Low)
}

func (ag *AddrRegister) inc(delta byte) {
	low := ag.Low
	ag.Low += delta
	if ag.Low < low {
		ag.High += 1
	}
	// 超过了可访问地址[,0x4000)
	if ag.get() > 0x3FFF {
		ag.set(ag.get() & 0x3FFF)
	}
}

// update 先修改高位再修改低位
func (ag *AddrRegister) update(data byte) {
	if ag.HighPtr {
		ag.High = data
	} else {
		ag.Low = data
	}
	// 超过了可访问地址[,0x4000)
	if ag.get() > 0x3FFF {
		ag.set(ag.get() & 0x3FFF)
	}
	ag.HighPtr = !ag.HighPtr
}

func (ag *AddrRegister) resetLatch() {
	ag.HighPtr = true
}
