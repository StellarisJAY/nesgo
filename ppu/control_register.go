package ppu

// ControlRegister ppu控制寄存器
type ControlRegister struct {
	val byte
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

func (cr *ControlRegister) VRAMIncrement() byte {
	if cr.get(AddrIncrement) {
		return 32
	} else {
		return 1
	}
}

func (cr *ControlRegister) Set(offset byte) {
	cr.val = cr.val | offset
}

func (cr *ControlRegister) get(offset byte) bool {
	return cr.val&offset != 0
}

func (cr *ControlRegister) Clear(offset byte) {
	cr.val = cr.val & (^offset)
}
