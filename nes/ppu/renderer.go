package ppu

import (
	"fmt"
)

type viewPort struct {
	x1 uint32
	y1 uint32
	x2 uint32
	y2 uint32
}

func newViewPort(x1, y1, x2, y2 uint32) viewPort {
	return viewPort{x1, y1, x2, y2}
}

// Render 渲染当前的frame
func (p *PPU) Render() {
	p.renderBackground()
	p.renderSprites()
}

// renderBackground 渲染background的960个tiles
func (p *PPU) renderBackground() {
	main, second := p.nameTables()
	scrollX, scrollY := uint32(p.scrollReg.X), uint32(p.scrollReg.Y)
	p.renderNameTable(main, newViewPort(scrollX, scrollY, 256, 240), -int32(scrollX), -int32(scrollY))
	if scrollX > 0 {
		p.renderNameTable(second, newViewPort(0, 0, scrollX, 240), 256-int32(scrollX), 0)
	} else if scrollY > 0 {
		p.renderNameTable(second, newViewPort(0, 0, 256, scrollY), 0, 240-int32(scrollY))
	}
}

// renderSprites 渲染oam中记录的所有sprites
func (p *PPU) renderSprites() {
	bank := p.ctrlReg.getSpritePattern()
	// oam数据记录64个sprite的位置和状态，每个sprite占用4字节
	for i := 0; i < len(p.oamData); i += 4 {
		// byte0 byte3 是 y和x
		y := uint16(p.oamData[i])
		x := uint16(p.oamData[i+3])
		index := p.oamData[i+1]
		attribute := p.oamData[i+2]
		// priority为0表示sprite在background后，跳过渲染
		if priority := attribute & (1 << 5); priority == 1 {
			continue
		}
		flipH := attribute&(1<<6) != 0
		flipV := attribute&(1<<7) != 0
		paletteIdx := attribute & 0b11
		if p.ctrlReg.isBigSprite() {
			// 渲染8x16 sprite
			pt := p.getChrBank(index & 1)
			// top和bottom两个tile的index
			topIdx := uint16(index & 0b11111110)
			bottomIdx := topIdx + 1
			p.renderSprite(x, y, flipH, flipV, pt[topIdx*16:topIdx*16+16], p.spritePalette(paletteIdx))
			p.renderSprite(x, y+8, flipH, flipV, pt[bottomIdx*16:bottomIdx*16+16], p.spritePalette(paletteIdx))
		} else {
			// 8x8 pixels small sprite
			patternTable := p.getChrBank(bank)
			idx := uint16(index)
			tile := patternTable[idx*16 : idx*16+16]
			p.renderSprite(x, y, flipH, flipV, tile, p.spritePalette(paletteIdx))
		}

	}
}

func (p *PPU) renderSprite(tileX, tileY uint16, flipH, flipV bool, tile []byte, palette [4]byte) {
	var y uint16 = 0
	// 每个tile有8 X height个像素
	for ; y < 8; y++ {
		// 一个像素是2bits，高位与低位分别在相距 height 字节的两个字节里面
		low := tile[y]
		high := tile[y+8]
		var x int16 = 7
		for ; x >= 0; x-- {
			// 像素颜色顺序按照大端序，从高位开始遍历
			colorId := ((high & 1) << 1) | (low & 1)
			low = low >> 1
			high = high >> 1
			// 将调色板的颜色编号映射到RGB颜色
			var color Color
			switch colorId {
			case 0:
				continue
			case 1:
				color = getRGBColor(palette[1])
			case 2:
				color = getRGBColor(palette[2])
			case 3:
				color = getRGBColor(palette[3])
			default:
				panic(fmt.Errorf("invalid color id: %d", colorId))
			}
			// flip
			x := uint16(x)
			switch {
			case flipH && flipV:
				p.frame.setPixel(uint32(tileX+7-x), uint32(tileY+7-y), color)
			case !flipH && !flipV:
				p.frame.setPixel(uint32(tileX+x), uint32(tileY+y), color)
			case !flipH && flipV:
				p.frame.setPixel(uint32(tileX+x), uint32(tileY+7-y), color)
			case flipH && !flipV:
				p.frame.setPixel(uint32(tileX+7-x), uint32(tileY+y), color)
			default:
			}
		}
	}
}

func (p *PPU) renderNameTable(nameTable []byte, port viewPort, shiftX, shiftY int32) {
	attributeTable := nameTable[0x3c0:]
	bank := p.ctrlReg.getBgPattern()
	patternTable := p.getChrBank(bank)
	var i uint16
	// 32x30tiles，共960个，每行32个
	for i = 0; i < 960; i++ {
		tileX, tileY := i%32*8, i/32*8
		idx := uint16(nameTable[i])
		tile := patternTable[idx*16 : idx*16+16]
		palette := p.bgPalette(tileY/8, tileX/8, attributeTable)
		var y uint16 = 0
		// 每个tile有8x8个像素
		for ; y < 8; y++ {
			// 一个像素是2bits，高位与低位分别在相距8字节的两个字节里面
			low := tile[y]
			high := tile[y+8]
			var x int16 = 7
			for ; x >= 0; x-- {
				// 像素颜色顺序按照大端序，从高位开始遍历
				colorId := ((high & 1) << 1) | (low & 1)
				low = low >> 1
				high = high >> 1
				// 将调色板的颜色编号映射到RGB颜色
				var color Color
				switch colorId {
				case 0:
					color = getRGBColor(palette[0])
				case 1:
					color = getRGBColor(palette[1])
				case 2:
					color = getRGBColor(palette[2])
				case 3:
					color = getRGBColor(palette[3])
				default:
					panic(fmt.Errorf("invalid color id: %d", colorId))
				}
				pixelX, pixelY := uint32(tileX+uint16(x)), uint32(tileY+y)
				if pixelX >= port.x1 && pixelX < port.x2 && pixelY >= port.y1 && pixelY < port.y2 {
					p.frame.setPixel(uint32(shiftX+int32(pixelX)), uint32(shiftY+int32(pixelY)), color)
				}
			}
		}
	}
}

// spritePalette 获取sprite的调色板数据
func (p *PPU) spritePalette(idx byte) [4]byte {
	start := 0x11 + int(idx)*4
	return [4]byte{
		0,
		p.paletteTable[start],
		p.paletteTable[start+1],
		p.paletteTable[start+2],
	}
}

// bgPalette 获取row，col位置tile的调色板
// row和col是以tile为单位的坐标，不是像素坐标
func (p *PPU) bgPalette(row, col uint16, attributeTable []byte) [4]byte {
	// 2x2的tiles组成一个meta tile，4个meta tiles的调色板编号组成attributeTable的一个字节
	// 一行32个tiles，所以一行有8个meta tiles
	attrIdx := row/4*8 + col/4
	paletteByte := attributeTable[attrIdx]
	// [0] [2]
	// [4] [6]
	// 一个attr字节表示4个相邻的meta tiles，每个meta的偏移如上
	// 取余除以2，得到每个tile在attr中的那个meta tile
	y := row % 4 / 2
	x := col % 4 / 2
	offset := 2 * (y*2 + x)
	paletteIdx := (paletteByte >> offset) & 0b11
	first := 1 + int(paletteIdx)*4
	return [4]byte{
		p.paletteTable[0],
		p.paletteTable[first],
		p.paletteTable[first+1],
		p.paletteTable[first+2],
	}
}

// 获取当前的nameTable
func (p *PPU) nameTables() (main, second []byte) {
	// 虚拟的地址编号，0:0x2000,1:0x2400, 2:0x2800, 3:0x2c00
	addr := p.ctrlReg.nameTableAddr()
	mirror := p.getMirroring()
	// Vertical   Horizontal   OneScreen
	// [A] [B]    [A] [a]      [A] [A]
	// [a] [b]    [B] [b]      [A] [A]
	// main为ctrl中编号对应的A或B，second是另外一个
	// 返回映射到物理地址的nameTable数据
	// Horizontal的物理地址0x2400是B
	switch {
	case (mirror == Vertical && addr == 0x2000) || (mirror == Vertical && addr == 0x2800):
		return p.ram[0:0x400], p.ram[0x400:0x800]
	case (mirror == Horizontal && addr == 0x2000) || (mirror == Horizontal && addr == 0x2400):
		return p.ram[0:0x400], p.ram[0x400:0x800]
	case (mirror == Vertical && addr == 0x2400) || (mirror == Vertical && addr == 0x2c00):
		return p.ram[0x400:0x800], p.ram[0:0x400]
	case (mirror == Horizontal && addr == 0x2800) || (mirror == Horizontal && addr == 0x2c00):
		return p.ram[0x400:0x800], p.ram[0:0x400]
	case mirror == OneScreenLow:
		return p.ram[0:0x400], p.ram[0x400:0x800]
	case mirror == OneScreenHigh:
		return p.ram[0x400:0x800], p.ram[0:0x400]
	default:
	}
	return nil, nil
}
