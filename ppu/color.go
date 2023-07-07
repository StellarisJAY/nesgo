package ppu

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

// getRGBAColor 将一个内存cell的值转换成RGBA颜色
func getRGBAColor(cell byte) Color {
	switch cell {
	case 0:
		return Color{0x0, 0x0, 0x0, 0xff}
	case 1:
		return Color{0xff, 0xff, 0xff, 0xff}
	default:
		return Color{0xff, 0xff, 0xff, 0x0}
	}
}

func (c Color) Uint32() uint32 {
	return uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
}
