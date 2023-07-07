package cpu

const (
	NOP byte = 0xEA
	BRK byte = 0x0
	TAX byte = 0xAA
	TAY byte = 0xA8
	TXA byte = 0x8A
	TYA byte = 0x98

	LDA_IM  byte = 0xA9
	LDA_ZP  byte = 0xA5
	LDA_ZPX byte = 0xB5
	LDA_ABS byte = 0xAD
	LDA_ABX byte = 0xBD
	LDA_ABY byte = 0xB9
	LDA_IX  byte = 0xA1
	LDA_IY  byte = 0xB1

	STA_ZP  byte = 0x85
	STA_ZPX byte = 0x95
	STA_ABS byte = 0x8D
	STA_ABX byte = 0x9D
	STA_ABY byte = 0x99
	STA_IX  byte = 0x81
	STA_IY  byte = 0x91

	STX_ZP  byte = 0x86
	STX_ZPY byte = 0x96
	STX_ABS byte = 0x8E

	STY_ZP  byte = 0x84
	STY_ZPX byte = 0x94
	STY_ABS byte = 0x8C

	LDX_IM  byte = 0xA2
	LDX_ZP  byte = 0xA6
	LDX_ZPY byte = 0xB6
	LDX_ABS byte = 0xAE
	LDX_ABY byte = 0xBE

	LDY_IM  byte = 0xA0
	LDY_ZP  byte = 0xA4
	LDY_ZPX byte = 0xB4
	LDY_ABS byte = 0xAC
	LDY_ABX byte = 0xBC

	INX     byte = 0xE8
	INY     byte = 0xC8
	INC_ZP  byte = 0xE6
	INC_ZPX byte = 0xF6
	INC_ABS byte = 0xEE
	INC_ABX byte = 0xFE

	JMP_A byte = 0x4C
	JMP_I byte = 0x6C

	TSX byte = 0xBA
	TXS byte = 0x9A
	PHA byte = 0x48
	PHP byte = 0x08
	PLA byte = 0x68
	PLP byte = 0x28

	AND_IM  byte = 0x29
	AND_ZP  byte = 0x25
	AND_ZPX byte = 0x35
	AND_ABS byte = 0x2D
	AND_ABX byte = 0x3D
	AND_ABY byte = 0x39
	AND_INX byte = 0x21
	AND_INY byte = 0x31

	EOR_IM  byte = 0x49
	EOR_ZP  byte = 0x45
	EOR_ZPX byte = 0x55
	EOR_ABS byte = 0x4D
	EOR_ABX byte = 0x5D
	EOR_ABY byte = 0x59
	EOR_INX byte = 0x41
	EOR_INY byte = 0x51

	ORA_IM  byte = 0x09
	ORA_ZP  byte = 0x05
	ORA_ZPX byte = 0x15
	ORA_ABS byte = 0x0D
	ORA_ABX byte = 0x1D
	ORA_ABY byte = 0x19
	ORA_INX byte = 0x01
	ORA_INY byte = 0x11

	BIT_ZP  byte = 0x24
	BIT_ABS byte = 0x2C

	ZeroStatus     byte = 1 << 1
	NegativeStatus byte = 1 << 7
)

// AddressMode 寻址模式
type AddressMode byte

const (
	Immediate      AddressMode = iota // 立即数
	ZeroPage                          // addr = mem8[pc]
	ZeroPageX                         // addr = mem8[pc] + regA
	ZeroPageY                         // addr = mem8[pc] + y
	Absolute                          // addr = mem16[pc]
	AbsoluteX                         // addr = mem16[pc] + regA
	AbsoluteY                         // addr = mem16[pc] + y
	IndirectX                         // addr = mem16[mem8[pc] + regA]
	IndirectY                         // addr = mem16[mem8[pc] + y]
	NoneAddressing                    // 指令不访问内存
)

// InstructionHandler 命令处理器
type InstructionHandler func(p *Processor, op Instruction)

// Instruction CPU指令
type Instruction struct {
	code     byte               // code 指令编码
	name     string             // name 指令名称
	length   byte               // length 完整指令的字节数，包括参数
	cycle    byte               // cycle  执行指令所需CPU周期
	addrMode AddressMode        // addrMode 寻址模式
	handler  InstructionHandler // handler 指令处理器（可选）
}

var (
	Instructions = map[byte]Instruction{
		NOP:     {NOP, "NOP", 1, 2, NoneAddressing, nil},
		BRK:     {BRK, "BRK", 1, 1, NoneAddressing, nil},
		TAX:     {TAX, "TAX", 1, 1, NoneAddressing, tax},
		TAY:     {TAY, "TAY", 1, 2, NoneAddressing, tay},
		TXA:     {TXA, "TXA", 1, 2, NoneAddressing, txa},
		TYA:     {TYA, "TYA", 1, 2, NoneAddressing, tya},
		INX:     {INX, "INX", 1, 1, NoneAddressing, nil},
		INY:     {INY, "INY", 1, 2, NoneAddressing, nil},
		INC_ZP:  {INC_ZP, "INC", 2, 5, ZeroPage, inc},
		INC_ZPX: {INC_ZPX, "INC", 2, 6, ZeroPageX, inc},
		INC_ABS: {INC_ABS, "INC", 3, 6, Absolute, inc},
		INC_ABX: {INC_ABX, "INC", 3, 7, AbsoluteX, inc},
		// LDA with different addressing modes
		LDA_IM:  {LDA_IM, "LDA", 2, 2, Immediate, lda},
		LDA_ZP:  {LDA_ZP, "LDA", 2, 3, ZeroPage, lda},
		LDA_ZPX: {LDA_ZPX, "LDA", 2, 4, ZeroPageX, lda},
		LDA_ABS: {LDA_ABS, "LDA", 3, 4, Absolute, lda},
		LDA_ABX: {LDA_ABX, "LDA", 3, 4, AbsoluteX, lda},
		LDA_ABY: {LDA_ABY, "LDA", 3, 4, AbsoluteY, lda},
		LDA_IX:  {LDA_IX, "LDA", 2, 6, IndirectX, lda},
		LDA_IY:  {LDA_IY, "LDA", 2, 5, IndirectY, lda},
		// STA with different addressing modes
		STA_ZP:  {STA_ZP, "STA", 2, 3, ZeroPage, sta},
		STA_ZPX: {STA_ZPX, "STA", 2, 4, ZeroPageX, sta},
		STA_ABS: {STA_ABS, "STA", 3, 4, Absolute, sta},
		STA_ABX: {STA_ABX, "STA", 3, 5, AbsoluteX, sta},
		STA_ABY: {STA_ABY, "STA", 3, 5, AbsoluteY, sta},
		STA_IX:  {STA_IX, "STA", 2, 6, IndirectX, sta},
		STA_IY:  {STA_IY, "STA", 2, 6, IndirectY, sta},

		// STX
		STX_ZP:  {STX_ZP, "STX", 2, 3, ZeroPage, stx},
		STX_ZPY: {STX_ZPY, "STX", 2, 4, ZeroPageY, stx},
		STX_ABS: {STX_ABS, "STX", 3, 4, Absolute, stx},
		// STY
		STY_ZP:  {STY_ZP, "STY", 2, 3, ZeroPage, sty},
		STY_ZPX: {STY_ZPX, "STY", 2, 4, ZeroPageX, sty},
		STY_ABS: {STY_ABS, "STY", 3, 4, Absolute, sty},
		// LDX
		LDX_IM:  {LDX_IM, "LDX", 2, 2, Immediate, ldx},
		LDX_ZP:  {LDX_ZP, "LDX", 2, 2, ZeroPage, ldx},
		LDX_ZPY: {LDX_ZPY, "LDX", 2, 2, ZeroPageY, ldx},
		LDX_ABS: {LDX_ABS, "LDX", 2, 2, Absolute, ldx},
		LDX_ABY: {LDX_ABY, "LDX", 2, 2, AbsoluteY, ldx},
		// LDY
		LDY_IM:  {LDY_IM, "LDY", 2, 2, Immediate, ldy},
		LDY_ZP:  {LDY_ZP, "LDY", 2, 3, ZeroPage, ldy},
		LDY_ZPX: {LDY_ZPX, "LDY", 2, 4, ZeroPageX, ldy},
		LDY_ABS: {LDY_ABS, "LDY", 3, 4, Absolute, ldy},
		LDY_ABX: {LDY_ABX, "LDY", 3, 4, AbsoluteX, ldy},
		// JMP
		JMP_A: {JMP_A, "JMP", 2, 2, NoneAddressing, jmp},
		JMP_I: {JMP_I, "JMP", 2, 2, NoneAddressing, jmp},
		// Stack
		TSX: {TSX, "TSX", 1, 2, NoneAddressing, tsx},
		TXS: {TXS, "TXS", 1, 2, NoneAddressing, txs},
		PHA: {PHA, "PHA", 1, 2, NoneAddressing, pha},
		PHP: {PHP, "PHP", 1, 2, NoneAddressing, php},
		PLA: {PLA, "PLA", 1, 2, NoneAddressing, pla},
		PLP: {PLP, "PLP", 1, 2, NoneAddressing, plp},

		// Logical AND A
		AND_IM:  {AND_IM, "AND", 2, 2, Immediate, and},
		AND_ZP:  {AND_ZP, "AND", 2, 3, ZeroPage, and},
		AND_ZPX: {AND_ZPX, "AND", 2, 4, ZeroPageX, and},
		AND_ABS: {AND_ABS, "AND", 3, 4, Absolute, and},
		AND_ABX: {AND_ABX, "AND", 3, 4, AbsoluteX, and},
		AND_ABY: {AND_ABY, "AND", 3, 4, AbsoluteY, and},
		AND_INX: {AND_INX, "AND", 2, 6, IndirectX, and},
		AND_INY: {AND_INY, "AND", 2, 5, IndirectY, and},
		// Logical EOR A
		EOR_IM:  {EOR_IM, "EOR", 2, 2, Immediate, eor},
		EOR_ZP:  {EOR_ZP, "EOR", 2, 3, ZeroPage, eor},
		EOR_ZPX: {EOR_ZPX, "EOR", 2, 4, ZeroPageX, eor},
		EOR_ABS: {EOR_ABS, "EOR", 3, 4, Absolute, eor},
		EOR_ABX: {EOR_ABX, "EOR", 3, 4, AbsoluteX, eor},
		EOR_ABY: {EOR_ABY, "EOR", 3, 4, AbsoluteY, eor},
		EOR_INX: {EOR_INX, "EOR", 2, 6, IndirectX, eor},
		EOR_INY: {EOR_INY, "EOR", 2, 5, IndirectY, eor},
		// Logical OR A
		ORA_IM:  {ORA_IM, "ORA", 2, 2, Immediate, ora},
		ORA_ZP:  {ORA_ZP, "ORA", 2, 3, ZeroPage, ora},
		ORA_ZPX: {ORA_ZPX, "ORA", 2, 4, ZeroPageX, ora},
		ORA_ABS: {ORA_ABS, "ORA", 3, 4, Absolute, ora},
		ORA_ABX: {ORA_ABX, "ORA", 3, 4, AbsoluteX, ora},
		ORA_ABY: {ORA_ABY, "ORA", 3, 4, AbsoluteY, ora},
		ORA_INX: {ORA_INX, "ORA", 2, 6, IndirectX, ora},
		ORA_INY: {ORA_INY, "ORA", 2, 5, IndirectY, ora},
		// Logical Bit Test A
		BIT_ZP:  {BIT_ZP, "BIT", 2, 3, ZeroPage, bit},
		BIT_ABS: {BIT_ABS, "BIT", 3, 4, Absolute, bit},
	}
)

func (p *Processor) zeroOrNegativeStatus(value byte) {
	if value == 0 {
		p.regStatus |= ZeroStatus
	} else {
		p.regStatus &= ^ZeroStatus
	}
	if value&(1<<7) != 0 {
		p.regStatus |= NegativeStatus
	} else {
		p.regStatus &= ^NegativeStatus
	}
}
