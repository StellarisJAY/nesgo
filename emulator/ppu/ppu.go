package ppu

import (
	"fmt"
	"slices"
)

const (
	Vertical byte = iota
	Horizontal
	FourScreen
	OneScreenLow
	OneScreenHigh
	OAMSize = 256
)

// PPU 图形处理器
type PPU struct {
	getChrBank     func(byte) []byte // getChrBank 动态获取chr的patternTable，mapper可以动态切换patternTable
	getMirroring   func() byte       // getMirroring 动态获取mirroring
	writeCHR       func(uint16, byte)
	paletteTable   []byte          // paletteTable 保存编号对应的颜色
	ram            []byte          // ram ppu RAM
	oamAddr        byte            // oamAddr 当前oam写地址
	oamData        []byte          // oamData sprite数据
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

type Snapshot struct {
	PaletteTable   []byte
	Ram            []byte
	OamAddr        byte
	OamData        []byte
	AddrReg        AddrRegister
	CtrlReg        ControlRegister
	MashReg        MaskRegister
	ScrollReg      ScrollRegister
	InternalBuffer byte
	StatReg        StatusRegister
	Cycles         uint64
	ScanLines      uint16
	NMIInterrupt   bool
	Frame          []byte
}

func NewPPU(getChrBank func(byte) []byte, getMirroring func() byte, writeCHR func(uint16, byte)) *PPU {
	return &PPU{
		getChrBank:   getChrBank,
		getMirroring: getMirroring,
		writeCHR:     writeCHR,
		paletteTable: make([]byte, 32),
		ram:          make([]byte, 2048),
		oamAddr:      0,
		oamData:      make([]byte, OAMSize),
		addrReg:      NewAddrRegister(),
		ctrlReg:      NewControlRegister(),
		statReg:      NewStatusRegister(),
		frame:        NewFrame(),
		maskReg:      NewMaskRegister(),
		scrollReg:    NewScrollRegister(),
	}
}

func (p *PPU) MakeSnapshot() Snapshot {
	return Snapshot{
		PaletteTable:   slices.Clone(p.paletteTable),
		Ram:            slices.Clone(p.ram),
		OamAddr:        p.oamAddr,
		OamData:        slices.Clone(p.oamData),
		AddrReg:        p.addrReg,
		CtrlReg:        p.ctrlReg,
		MashReg:        p.maskReg,
		ScrollReg:      p.scrollReg,
		InternalBuffer: p.internalBuffer,
		StatReg:        p.statReg,
		Cycles:         p.cycles,
		ScanLines:      p.scanLines,
		NMIInterrupt:   p.nmiInterrupt,
	}
}

func (p *PPU) Reverse(s Snapshot) []byte {
	rev := PPU{
		getChrBank:     p.getChrBank,
		getMirroring:   p.getMirroring,
		writeCHR:       p.writeCHR,
		paletteTable:   s.PaletteTable,
		ram:            s.Ram,
		oamAddr:        s.OamAddr,
		oamData:        s.OamData,
		addrReg:        s.AddrReg,
		ctrlReg:        s.CtrlReg,
		maskReg:        s.MashReg,
		scrollReg:      s.ScrollReg,
		internalBuffer: s.InternalBuffer,
		statReg:        s.StatReg,
		cycles:         s.Cycles,
		scanLines:      s.ScanLines,
		nmiInterrupt:   s.NMIInterrupt,
		frame:          p.frame,
	}
	*p = rev
	return p.FrameData()
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
		p.internalBuffer = p.getChrBank(byte(addr / 0x1000))[addr%0x1000]
		return result
	case addr >= 0x2000 && addr <= 0x2fff: // ppu RAM
		result := p.internalBuffer
		mirrorAddr := p.mirrorVRAMAddr(addr)
		p.internalBuffer = p.ram[mirrorAddr]
		return result
	case addr <= 0x3eff: // 0x3000~0x3eff映射到0x2000~0x2eff
		result := p.internalBuffer
		p.internalBuffer = p.ram[p.mirrorVRAMAddr(addr-0x1000)]
		return result
	case addr <= 0x3fff:
		// mirror down to 32
		addrMirror := (addr - 0x3f00) % 32
		// 调色板的数据读取直接返回
		return p.paletteTable[addrMirror]
	default:
		panic(fmt.Errorf("invalid ppu memory addr 0x%X", addr))
	}
	return 0
}

func (p *PPU) WriteData(val byte) {
	addr := p.addrReg.get()
	switch {
	case addr <= 0x1fff: // chr ROM
		p.writeCHR(addr, val)
	case addr <= 0x2fff: // ppu RAM
		p.ram[p.mirrorVRAMAddr(addr)] = val
	case addr <= 0x3eff: // 0x3000~0x3eff映射到0x2000~0x2eff
		p.ram[p.mirrorVRAMAddr(addr-0x1000)] = val
	case addr == 0x3f10 || addr == 0x3f14 || addr == 0x3f18 || addr == 0x3f1c: // mirroring to palette
		addr = addr - 0x10
		p.paletteTable[addr-0x3f00] = val
	case addr <= 0x3fff:
		// mirror down to 32
		addrMirror := (addr - 0x3f00) % 32
		// 调色板的数据
		p.paletteTable[addrMirror] = val
	default:
		panic(fmt.Errorf("invalid ppu memory addr 0x%X", addr))
	}
	p.incrementAddr()
}

// 0x2000到0x3fff一共4KiB虚拟空间，其中一个32x32的nameTable为1KiB，所以空间被划分为了4份
// 实际内存只有0x2000~0x7fff，所以需要将0x2800以后的虚拟地址映射
func (p *PPU) mirrorVRAMAddr(addr uint16) uint16 {
	idx := addr - 0x2000
	nameTable := idx / 0x0400
	if nameTable == 0 {
		return idx
	}
	switch p.getMirroring() {
	// Vertical，A：0x2000~0x23ff和0x2800~0x2bff
	// [A] [B]
	// [a] [b]
	case Vertical:
		if nameTable == 2 || nameTable == 3 {
			return idx - 0x800
		}
	case Horizontal:
		// Horizontal, A: 0x2000 a: 0x2400，B: 0x2800 b: 0x2c00
		// 因为ram只有2KiB空间，0x2400实际上是B的数据
		// 所以，B和a减去0x400，b减去0x800
		if nameTable == 3 {
			return idx - 0x800
		} else {
			return idx - 0x400
		}
	case OneScreenLow:
		// OneScreenLow: [0x2000,0x2400)
		return idx % 0x400
	case OneScreenHigh:
		// OneScreenHigh: [0x2400, 0x2800)
		return (idx % 0x400) + 0x400
	default:
	}
	return idx
}

func (p *PPU) Tick(cycles uint64) bool {
	p.cycles += cycles
	if p.cycles >= 341 {
		if p.isSprite0Hit(p.cycles) {
			p.statReg.setSprite0Hit()
		}
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
	before := p.ctrlReg.get(GenerateNMI)
	p.ctrlReg.Set(val)
	if !before && p.ctrlReg.get(GenerateNMI) && p.statReg.isVBlank() {
		p.nmiInterrupt = true
	}
}

func (p *PPU) WriteAddrReg(val byte) {
	p.addrReg.update(val)
	p.scrollReg.reset()
}

func (p *PPU) ReadStatus() byte {
	status := p.statReg.Val
	// 读取状态会导致reset vblank addr
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
	return p.frame.Data()
}

func (p *PPU) Frame() *Frame {
	return p.frame
}

func (p *PPU) WriteMask(val byte) {
	p.maskReg.set(val)
}

func (p *PPU) WriteScroll(val byte) {
	p.scrollReg.write(val)
}

func (p *PPU) WriteStatus(val byte) {
	p.statReg.Val = val
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
	result := p.oamData[p.oamAddr]
	return result
}

func (p *PPU) isSprite0Hit(cycles uint64) bool {
	y := uint64(p.oamData[0])
	x := uint64(p.oamData[3])
	return (y == uint64(p.scanLines)) && (x <= cycles) && p.maskReg.getBit(ShowSprite)
}
