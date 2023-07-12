package ppu

// AddrRegister ppu地址寄存器
// 通过修改0x2006寄存器来设置ppu地址
type AddrRegister struct {
	low     byte
	high    byte
	highPtr bool
}

func NewAddrRegister() AddrRegister {
	return AddrRegister{
		low:     0,
		high:    0,
		highPtr: true,
	}
}

func (ag *AddrRegister) set(addr uint16) {
	ag.high = byte(addr >> 8)
	ag.low = byte(addr & 0xFF)
}

func (ag *AddrRegister) get() uint16 {
	return uint16(ag.high)<<8 + uint16(ag.low)
}

func (ag *AddrRegister) inc(delta byte) {
	low := ag.low
	ag.low += delta
	if ag.low < low {
		ag.high += 1
	}
	// 超过了可访问地址[,0x4000)
	if ag.get() > 0x3FFF {
		ag.set(ag.get() & 0x3FFF)
	}
}

// update 先修改高位再修改低位
func (ag *AddrRegister) update(data byte) {
	if ag.highPtr {
		ag.high = data
	} else {
		ag.low = data
	}
	// 超过了可访问地址[,0x4000)
	if ag.get() > 0x3FFF {
		ag.set(ag.get() & 0x3FFF)
	}
	ag.highPtr = !ag.highPtr
}

func (ag *AddrRegister) resetLatch() {
	ag.highPtr = true
}
