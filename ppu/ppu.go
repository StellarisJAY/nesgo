package ppu

import "fmt"

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
	oamData        []byte          // oamData
	mirroring      byte            // mirroring
	addrReg        AddrRegister    // addrReg 地址寄存器，因为ppu读取是异步的，需要寄存器记录读请求的地址
	ctrlReg        ControlRegister // ctrlReg ppu 控制寄存器
	internalBuffer byte            // internalBuffer 异步读取缓冲区
	statReg        StatusRegister  // statReg 状态寄存器

	cycles    uint64 // cycles ppu 经过的时钟周期
	scanLines uint16 // scanLines

	nmiInterrupt bool
}

func NewPPU(chrROM []byte, mirroring byte) *PPU {
	return &PPU{
		chrROM:       chrROM,
		paletteTable: make([]byte, 32),
		ram:          make([]byte, 2048),
		oamData:      make([]byte, 256),
		mirroring:    mirroring,
		addrReg:      NewAddrRegister(),
		ctrlReg:      NewControlRegister(),
		statReg:      NewStatusRegister(),
	}
}

func (p *PPU) incrementAddr() {
	p.addrReg.inc(p.ctrlReg.VRAMIncrement())
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
		return vramAddr
	}
	switch p.mirroring {
	case Vertical:
		if nameTable == 2 || nameTable == 3 {
			return vramAddr - 0x800
		}
	case Horizontal:
		if nameTable == 1 || nameTable == 2 {
			return vramAddr - 0x400
		} else {
			return vramAddr - 0x800
		}
	default:
	}
	return vramAddr
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
	return status
}

func (p *PPU) IsInterrupt() bool {
	return p.nmiInterrupt
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
