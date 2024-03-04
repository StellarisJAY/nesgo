package ppu

// ControlRegister ppu控制寄存器
type ControlRegister struct {
	Val byte
}

const (
	AddrIncrement     byte = 1 << 2
	SpritePattern     byte = 1 << 3
	BackgroundPattern byte = 1 << 4
	SpriteSize        byte = 1 << 5
	GenerateNMI       byte = 1 << 7
)

func NewControlRegister() ControlRegister {
	return ControlRegister{0}
}

func (cr *ControlRegister) nameTableAddr() uint16 {
	switch cr.Val & 0b11 {
	case 0:
		return 0x2000
	case 1:
		return 0x2400
	case 2:
		return 0x2800
	case 3:
		return 0x2c00
	default:
		panic("can't happen")
	}
}

func (cr *ControlRegister) VRAMIncrement() byte {
	if cr.get(AddrIncrement) {
		return 32
	} else {
		return 1
	}
}

func (cr *ControlRegister) Set(data byte) {
	cr.Val = data
}

func (cr *ControlRegister) get(offset byte) bool {
	return cr.Val&offset != 0
}

func (cr *ControlRegister) Clear(offset byte) {
	cr.Val = cr.Val & (^offset)
}

func (cr *ControlRegister) getBgPattern() byte {
	if cr.get(BackgroundPattern) {
		return 1
	} else {
		return 0
	}
}

func (cr *ControlRegister) getSpritePattern() byte {
	if cr.get(SpritePattern) {
		return 1
	} else {
		return 0
	}
}

func (cr *ControlRegister) isBigSprite() bool {
	return cr.get(SpriteSize)
}
