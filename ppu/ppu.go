package ppu

import (
	"fmt"
)

const (
	Vertical byte = iota
	Horizontal
	FourScreen
	OAMSize = 256
)

// PPU 图形处理器
type PPU struct {
	chrROM         []byte          // chrROM characterROM 保存Sprite静态数据
	paletteTable   []byte          // paletteTable 保存编号对应的颜色
	ram            []byte          // ram ppu RAM
	oamAddr        byte            // oamAddr
	oamData        []byte          // oamData
	mirroring      byte            // mirroring
	addrReg        AddrRegister    // addrReg 地址寄存器，因为ppu读取是异步的，需要寄存器记录读请求的地址
	ctrlReg        ControlRegister // ctrlReg ppu 控制寄存器
	maskReg        MaskRegister
	scrollReg      ScrollRegister
	internalBuffer byte           // internalBuffer 异步读取缓冲区
	statReg        StatusRegister // statReg 状态寄存器

	cycles       uint64 // cycles ppu 经过的时钟周期
	scanLines    uint16 // scanLines
	nmiInterrupt bool
	frame        *Frame
}

func NewPPU(chrROM []byte, mirroring byte) *PPU {
	return &PPU{
		chrROM:       chrROM,
		paletteTable: make([]byte, 32),
		ram:          make([]byte, 2048),
		oamAddr:      0,
		oamData:      make([]byte, 256),
		mirroring:    mirroring,
		addrReg:      NewAddrRegister(),
		ctrlReg:      NewControlRegister(),
		statReg:      NewStatusRegister(),
		frame:        NewFrame(),
		maskReg:      NewMaskRegister(),
		scrollReg:    NewScrollRegister(),
	}
}

func (p *PPU) incrementAddr() {
	p.addrReg.inc(p.ctrlReg.VRAMIncrement())
}

// Render 渲染当前的NameTable
func (p *PPU) Render() {
	var bank uint16
	if p.ctrlReg.get(BackgroundPattern) {
		bank = 1
	} else {
		bank = 0
	}
	var i uint16
	for i = 0; i < 960; i++ {
		x, y := i%32*8, i/32*8
		tileIndex := p.ram[i]
		p.renderTile(x, y, bank, tileIndex)
	}
}

// DisplayAllTiles 测试方法，在frame中渲染bank中的所有tiles
func (p *PPU) DisplayAllTiles() {
	var bank uint16 = 0
	var x, y uint16 = 0, 0
	var i byte = 0
	for bank <= 1 {
		if x >= 256 || x+8 >= 256 {
			x = 0
			y += 10
		}
		p.renderTile(x, y, bank, byte(i))
		x += 10
		if i == 255 {
			i = 0
			bank += 1
		} else {
			i++
		}
	}

}

// renderTile 在Frame的x，y位置渲染一个tile
func (p *PPU) renderTile(x, y uint16, bank uint16, tileIndex byte) {
	idx := uint16(tileIndex)
	bankBase := bank * 0x1000
	tile := p.chrROM[bankBase+idx*16 : bankBase+idx*16+16]
	var row uint16 = 0
	// 每个tile有8x8个像素
	for ; row < 8; row++ {
		// 一个像素是2bits，高位与低位分别在相距8字节的两个字节里面
		low := tile[row]
		high := tile[row+8]
		var col uint16 = 0
		for ; col < 8; col++ {
			colorId := (((high >> (8 - col)) & 1) << 1) | ((low >> (8 - col)) & 1)
			var color Color
			switch colorId {
			case 0:
				color = SystemPalette[0x01]
			case 1:
				color = SystemPalette[0x23]
			case 2:
				color = SystemPalette[0x27]
			case 3:
				color = SystemPalette[0x30]
			default:
				panic(fmt.Errorf("invalid color id: %d", colorId))
			}
			p.frame.setPixel(uint32(x+col), uint32(y+row), color)
		}
	}
}

// ReadData 返回上一个读取请求的结果，并将本次读取请求的结果放入buffer
// 因为NES的ppu读取和CPU指令是异步的，需要用buffer来模拟异步过程
func (p *PPU) ReadData() byte {
	addr := p.addrReg.get()
	p.incrementAddr()
	switch {
	case addr <= 0x1fff: // chr ROM
		result := p.internalBuffer
		p.internalBuffer = p.chrROM[addr]
		return result
	case addr <= 0x2fff: // ppu RAM
		result := p.internalBuffer
		p.internalBuffer = p.ram[p.mirrorVRAMAddr(addr)]
		return result
	case addr <= 0x3eff:
		panic("can't read memory between [0x3000, 0x3eff)")
	case addr <= 0x3fff:
		// 调色板的数据读取直接返回
		return p.paletteTable[addr-0x3f00]
	default:
		panic(fmt.Errorf("invalid ppu memory addr 0x%x", addr))
	}
	return 0
}

func (p *PPU) WriteData(val byte) {
	addr := p.addrReg.get()
	p.incrementAddr()
	switch {
	case addr <= 0x1fff: // chr ROM
		panic(fmt.Errorf("can't write chr ROM addr: 0x%x", addr))
	case addr <= 0x2fff: // ppu RAM
		p.ram[p.mirrorVRAMAddr(addr)] = val
	case addr <= 0x3eff:
		panic("can't read memory between [0x3000, 0x3eff)")
	case addr == 0x3f10 || addr == 0x3f14 || addr == 0x3f18 || addr == 0x3f1c: // mirroring to palette
		addr = addr - 0x10
		p.paletteTable[addr-0x3f00] = val
	case addr <= 0x3fff:
		// 调色板的数据
		p.paletteTable[addr-0x3f00] = val
	default:
		panic(fmt.Errorf("invalid ppu memory addr 0x%x", addr))
	}
}

// 0x2000到0x3fff一共4KiB空间，其中一个32x32的nameTable为1KiB，所以空间被划分为了4份
// Horizontal, 空间被划分成A,B两份，其中0x2000~0x23ff和0x2400~0x27ff是A
// [A] [a]
// [B] [b]
// Vertical，A：0x2000~0x23ff和0x2800~0x2bff
// [A] [B]
// [a] [b]
// FourScreen
// [A] [B]
// [C] [D]
// SingleScreen
// [A] [a]
// [a] [a]
func (p *PPU) mirrorVRAMAddr(addr uint16) uint16 {
	vramAddr := addr & 0x3fff
	idx := vramAddr - 0x2000
	nameTable := idx / 0x0400
	if nameTable == 0 {
		return idx
	}
	switch p.mirroring {
	case Vertical:
		if nameTable == 2 || nameTable == 3 {
			return idx - 0x800
		}
	case Horizontal:
		if nameTable == 1 || nameTable == 2 {
			return idx - 0x400
		} else {
			return idx - 0x800
		}
	default:
	}
	return idx
}

func (p *PPU) Tick(cycles uint64) bool {
	p.cycles += cycles
	if p.cycles >= 341 {
		p.cycles -= 341
		p.scanLines += 1
		if p.scanLines == 241 {
			p.statReg.setVBlankStarted()
			p.statReg.resetSprite0Hit()
			if p.ctrlReg.get(GenerateNMI) {
				p.nmiInterrupt = true
			}
		}
		if p.scanLines >= 262 {
			p.scanLines = 0
			p.nmiInterrupt = false
			p.statReg.resetVBlankStarted()
			p.statReg.resetSprite0Hit()
			return true
		}
	}
	return false
}

func (p *PPU) WriteControl(val byte) {
	p.ctrlReg.Set(val)
}

func (p *PPU) WriteAddrReg(val byte) {
	p.addrReg.update(val)
}

func (p *PPU) ReadStatus() byte {
	status := p.statReg.val
	// 读取状态会导致reset vblan、addr
	p.statReg.resetVBlankStarted()
	p.addrReg.resetLatch()
	p.scrollReg.resetLatch()
	return status
}

func (p *PPU) PollInterrupt() bool {
	if interrupt := p.nmiInterrupt; interrupt {
		p.nmiInterrupt = false
		return interrupt
	}
	return false
}

func (p *PPU) PeekInterrupt() bool {
	return p.nmiInterrupt
}

func (p *PPU) FrameData() []byte {
	return p.frame.data
}

func (p *PPU) WriteMask(val byte) {
	p.maskReg.set(val)
}

func (p *PPU) WriteScroll(val byte) {
	p.scrollReg.write(val)
}

func (p *PPU) WriteOamAddr(val byte) {
	p.oamAddr = val
}

func (p *PPU) WriteOam(data byte) {
	p.oamData[p.oamAddr] = data
	p.oamAddr += 1
}

func (p *PPU) WriteOamDMA(data []byte) {
	for i := 0; i < 256; i++ {
		p.oamData[p.oamAddr] = data[i]
		p.oamAddr += 1
	}
}

func (p *PPU) ReadOam() byte {
	return p.oamData[p.oamAddr]
}

// ReadAndUpdateScreen 从内存读取每个cell，并在frame中修改cell的值，如果有更新则通知给渲染器
func ReadAndUpdateScreen(memory []byte, frame []byte) (updated bool) {
	idx := 0
	for _, cell := range memory {
		color := getRGBAColor(cell)
		r, g, b := color.R, color.G, color.B
		if frame[idx] != r || frame[idx+1] != g || frame[idx+2] != b {
			frame[idx] = r
			frame[idx+1] = g
			frame[idx+2] = b
			updated = true
		}
		idx += 3
	}
	return
}
