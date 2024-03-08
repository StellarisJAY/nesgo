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

	DEC_ZP  byte = 0xC6
	DEC_ZPX byte = 0xD6
	DEC_ABS byte = 0xCE
	DEC_ABX byte = 0xDE

	JMP_A byte = 0x4C
	JMP_I byte = 0x6C
	JSR   byte = 0x20
	RTS   byte = 0x60

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

	CLC byte = 0x18
	CLD byte = 0xD8
	CLI byte = 0x58
	CLV byte = 0xB8
	SEC byte = 0x38
	SED byte = 0xF8
	SEI byte = 0x78

	BCC byte = 0x90
	BCS byte = 0xB0
	BVC byte = 0x50
	BVS byte = 0x70
	BEQ byte = 0xF0
	BNE byte = 0xD0
	BMI byte = 0x30
	BPL byte = 0x10

	CMP_IM  byte = 0xC9
	CMP_ZP  byte = 0xC5
	CMP_ZPX byte = 0xD5
	CMP_ABS byte = 0xCD
	CMP_ABX byte = 0xDD
	CMP_ABY byte = 0xD9
	CMP_INX byte = 0xC1
	CMP_INY byte = 0xD1

	CPX_IM  byte = 0xE0
	CPX_ZP  byte = 0xE4
	CPX_ABS byte = 0xEC

	CPY_IM  byte = 0xC0
	CPY_ZP  byte = 0xC4
	CPY_ABS byte = 0xCC

	ADC_IM  byte = 0x69
	ADC_ZP  byte = 0x65
	ADC_ZPX byte = 0x75
	ADC_ABS byte = 0x6D
	ADC_ABX byte = 0x7D
	ADC_ABY byte = 0x79
	ADC_INX byte = 0x61
	ADC_INY byte = 0x71

	ASL     byte = 0x0A
	ASL_ZP  byte = 0x06
	ASL_ZPX byte = 0x16
	ASL_ABS byte = 0x0E
	ASL_ABX byte = 0x1E

	LSR     byte = 0x4A
	LSR_ZP  byte = 0x46
	LSR_ZPX byte = 0x56
	LSR_ABS byte = 0x4E
	LSR_ABX byte = 0x5E

	DEX byte = 0xCA
	DEY byte = 0x88

	SBC_IM  byte = 0xE9
	SBC_ZP  byte = 0xE5
	SBC_ZPX byte = 0xF5
	SBC_ABS byte = 0xED
	SBC_ABX byte = 0xFD
	SBC_ABY byte = 0xF9
	SBC_INX byte = 0xE1
	SBC_INY byte = 0xF1

	ROR     byte = 0x6A
	ROR_ZP  byte = 0x66
	ROR_ZPX byte = 0x76
	ROR_ABS byte = 0x6E
	ROR_ABX byte = 0x7E

	ROL     byte = 0x2A
	ROL_ZP  byte = 0x26
	ROL_ZPX byte = 0x36
	ROL_ABS byte = 0x2E
	ROL_ABX byte = 0x3E

	RTI byte = 0x40

	SRE_ZP  byte = 0x47
	SRE_ZPX byte = 0x57
	SRE_ABS byte = 0x4F
	SRE_ABX byte = 0x5F
	SRE_ABY byte = 0x5B
	SRE_INX byte = 0x43
	SRE_INY byte = 0x53
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
type InstructionHandler func(p *Processor, op *Instruction)

// Instruction CPU指令
type Instruction struct {
	Code     byte               // Code 指令编码
	Name     string             // Name 指令名称
	Length   byte               // Length 完整指令的字节数，包括参数
	Cycle    byte               // Cycle  执行指令所需CPU周期
	AddrMode AddressMode        // AddrMode 寻址模式
	handler  InstructionHandler // handler 指令处理器（可选）
}

var (
	Instructions = map[byte]*Instruction{
		NOP: {NOP, "NOP", 1, 2, NoneAddressing, nil},
		BRK: {BRK, "BRK", 1, 7, NoneAddressing, nil},
		TAX: {TAX, "TAX", 1, 2, NoneAddressing, tax},
		TAY: {TAY, "TAY", 1, 2, NoneAddressing, tay},
		TXA: {TXA, "TXA", 1, 2, NoneAddressing, txa},
		TYA: {TYA, "TYA", 1, 2, NoneAddressing, tya},
		INX: {INX, "INX", 1, 2, NoneAddressing, nil},
		INY: {INY, "INY", 1, 2, NoneAddressing, nil},
		DEX: {DEX, "DEX", 1, 2, NoneAddressing, dex},
		DEY: {DEY, "DEY", 1, 2, NoneAddressing, dey},

		INC_ZP:  {INC_ZP, "INC", 2, 5, ZeroPage, inc},
		INC_ZPX: {INC_ZPX, "INC", 2, 6, ZeroPageX, inc},
		INC_ABS: {INC_ABS, "INC", 3, 6, Absolute, inc},
		INC_ABX: {INC_ABX, "INC", 3, 7, AbsoluteX, inc},

		DEC_ZP:  {DEC_ZP, "DEC", 2, 5, ZeroPage, dec},
		DEC_ZPX: {DEC_ZPX, "DEC", 2, 6, ZeroPageX, dec},
		DEC_ABS: {DEC_ABS, "DEC", 3, 6, Absolute, dec},
		DEC_ABX: {DEC_ABX, "DEC", 3, 7, AbsoluteX, dec},
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
		LDX_ZP:  {LDX_ZP, "LDX", 2, 3, ZeroPage, ldx},
		LDX_ZPY: {LDX_ZPY, "LDX", 2, 4, ZeroPageY, ldx},
		LDX_ABS: {LDX_ABS, "LDX", 3, 4, Absolute, ldx},
		LDX_ABY: {LDX_ABY, "LDX", 3, 4, AbsoluteY, ldx},
		// LDY
		LDY_IM:  {LDY_IM, "LDY", 2, 2, Immediate, ldy},
		LDY_ZP:  {LDY_ZP, "LDY", 2, 3, ZeroPage, ldy},
		LDY_ZPX: {LDY_ZPX, "LDY", 2, 4, ZeroPageX, ldy},
		LDY_ABS: {LDY_ABS, "LDY", 3, 4, Absolute, ldy},
		LDY_ABX: {LDY_ABX, "LDY", 3, 4, AbsoluteX, ldy},
		// JMP
		JMP_A: {JMP_A, "JMP", 3, 3, Absolute, jmp},
		JMP_I: {JMP_I, "JMP", 3, 5, IndirectY, jmpIndirect},
		JSR:   {JSR, "JSR", 3, 6, Absolute, jsr},
		RTS:   {RTS, "RTS", 1, 6, NoneAddressing, rts},
		// Stack
		TSX: {TSX, "TSX", 1, 2, NoneAddressing, tsx},
		TXS: {TXS, "TXS", 1, 2, NoneAddressing, txs},
		PHA: {PHA, "PHA", 1, 3, NoneAddressing, pha},
		PHP: {PHP, "PHP", 1, 3, NoneAddressing, php},
		PLA: {PLA, "PLA", 1, 4, NoneAddressing, pla},
		PLP: {PLP, "PLP", 1, 4, NoneAddressing, plp},

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
		// Set and clear Status Register
		CLC: {CLC, "CLC", 1, 2, NoneAddressing, clc},
		CLD: {CLD, "CLD", 1, 2, NoneAddressing, cld},
		CLI: {CLI, "CLI", 1, 2, NoneAddressing, cli},
		CLV: {CLV, "CLV", 1, 2, NoneAddressing, clv},
		SEC: {SEC, "SEC", 1, 2, NoneAddressing, sec},
		SED: {SED, "SED", 1, 2, NoneAddressing, sed},
		SEI: {SEI, "SEI", 1, 2, NoneAddressing, sei},
		// Branch
		BCC: {BCC, "BCC", 2, 2, NoneAddressing, bcc},
		BCS: {BCS, "BCS", 2, 2, NoneAddressing, bcs},
		BVC: {BVC, "BVC", 2, 2, NoneAddressing, bvc},
		BVS: {BVS, "BVS", 2, 2, NoneAddressing, bvs},
		BEQ: {BEQ, "BEQ", 2, 2, NoneAddressing, beq},
		BNE: {BNE, "BNE", 2, 2, NoneAddressing, bne},
		BMI: {BMI, "BMI", 2, 2, NoneAddressing, bmi},
		BPL: {BPL, "BPL", 2, 2, NoneAddressing, bpl},
		// CMP
		CMP_IM:  {CMP_IM, "CMP", 2, 2, Immediate, cmp},
		CMP_ZP:  {CMP_ZP, "CMP", 2, 3, ZeroPage, cmp},
		CMP_ZPX: {CMP_ZPX, "CMP", 2, 4, ZeroPageX, cmp},
		CMP_ABS: {CMP_ABS, "CMP", 3, 4, Absolute, cmp},
		CMP_ABX: {CMP_ABX, "CMP", 3, 4, AbsoluteX, cmp},
		CMP_ABY: {CMP_ABY, "CMP", 3, 4, AbsoluteY, cmp},
		CMP_INX: {CMP_INX, "CMP", 2, 6, IndirectX, cmp},
		CMP_INY: {CMP_INY, "CMP", 2, 5, IndirectY, cmp},
		// CPX CPY
		CPX_IM:  {CPX_IM, "CPX", 2, 2, Immediate, cpx},
		CPX_ZP:  {CPX_ZP, "CPX", 2, 3, ZeroPage, cpx},
		CPX_ABS: {CPX_ABS, "CPX", 3, 4, Absolute, cpx},
		CPY_IM:  {CPY_IM, "CPY", 2, 2, Immediate, cpy},
		CPY_ZP:  {CPY_ZP, "CPY", 2, 3, ZeroPage, cpy},
		CPY_ABS: {CPY_ABS, "CPY", 3, 4, Absolute, cpy},
		// ADC
		ADC_IM:  {ADC_IM, "ADC", 2, 2, Immediate, adc},
		ADC_ZP:  {ADC_ZP, "ADC", 2, 3, ZeroPage, adc},
		ADC_ZPX: {ADC_ZPX, "ADC", 2, 4, ZeroPageX, adc},
		ADC_ABS: {ADC_ABS, "ADC", 3, 4, Absolute, adc},
		ADC_ABX: {ADC_ABX, "ADC", 3, 4, AbsoluteX, adc},
		ADC_ABY: {ADC_ABY, "ADC", 3, 4, AbsoluteY, adc},
		ADC_INX: {ADC_INX, "ADC", 2, 6, IndirectX, adc},
		ADC_INY: {ADC_INY, "ADC", 2, 5, IndirectY, adc},

		ASL:     {ASL, "ASL", 1, 2, NoneAddressing, asl},
		ASL_ZP:  {ASL_ZP, "ASL", 2, 5, ZeroPage, asl},
		ASL_ZPX: {ASL_ZPX, "ASL", 2, 6, ZeroPageX, asl},
		ASL_ABS: {ASL_ABS, "ASL", 3, 6, Absolute, asl},
		ASL_ABX: {ASL_ABX, "ASL", 3, 7, AbsoluteX, asl},

		LSR:     {LSR, "LSR", 1, 2, NoneAddressing, lsr},
		LSR_ZP:  {LSR_ZP, "LSR", 2, 5, ZeroPage, lsr},
		LSR_ZPX: {LSR_ZPX, "LSR", 2, 6, ZeroPageX, lsr},
		LSR_ABS: {LSR_ABS, "LSR", 3, 6, Absolute, lsr},
		LSR_ABX: {LSR_ABX, "LSR", 3, 7, AbsoluteX, lsr},

		SBC_IM:  {SBC_IM, "SBC", 2, 2, Immediate, sbc},
		SBC_ZP:  {SBC_ZP, "SBC", 2, 3, ZeroPage, sbc},
		SBC_ZPX: {SBC_ZPX, "SBC", 2, 4, ZeroPageX, sbc},
		SBC_ABS: {SBC_ABS, "SBC", 3, 4, Absolute, sbc},
		SBC_ABX: {SBC_ABX, "SBC", 3, 4, AbsoluteX, sbc},
		SBC_ABY: {SBC_ABY, "SBC", 3, 4, AbsoluteY, sbc},
		SBC_INX: {SBC_INX, "SBC", 2, 6, IndirectX, sbc},
		SBC_INY: {SBC_INY, "SBC", 2, 5, IndirectY, sbc},

		ROR:     {ROR, "ROR", 1, 2, NoneAddressing, ror},
		ROR_ZP:  {ROR_ZP, "ROR", 2, 5, ZeroPage, ror},
		ROR_ZPX: {ROR_ZPX, "ROR", 2, 6, ZeroPageX, ror},
		ROR_ABS: {ROR_ABS, "ROR", 3, 6, Absolute, ror},
		ROR_ABX: {ROR_ABX, "ROR", 3, 7, AbsoluteX, ror},

		ROL:     {ROL, "ROL", 1, 2, NoneAddressing, rol},
		ROL_ZP:  {ROL_ZP, "ROL", 2, 5, ZeroPage, rol},
		ROL_ZPX: {ROL_ZPX, "ROL", 2, 6, ZeroPageX, rol},
		ROL_ABS: {ROL_ABS, "ROL", 3, 6, Absolute, rol},
		ROL_ABX: {ROL_ABX, "ROL", 3, 7, AbsoluteX, rol},

		RTI: {RTI, "RTI", 1, 6, NoneAddressing, rti},

		// Unofficial ops
		SRE_ZP:  {SRE_ZP, "SRE", 2, 5, ZeroPage, sre},
		SRE_ZPX: {SRE_ZPX, "SRE", 2, 6, ZeroPageX, sre},
		SRE_ABS: {SRE_ABS, "SRE", 3, 6, Absolute, sre},
		SRE_ABX: {SRE_ABX, "SRE", 3, 7, AbsoluteX, sre},
		SRE_ABY: {SRE_ABY, "SRE", 3, 7, AbsoluteY, sre},
		SRE_INX: {SRE_INX, "SRE", 2, 8, IndirectX, sre},
		SRE_INY: {SRE_INY, "SRE", 2, 8, IndirectY, sre},
		// Unofficial NOPs
		0x04: {0x04, "NOP", 2, 3, ZeroPage, nopRead},
		0x44: {0x44, "NOP", 2, 3, ZeroPage, nopRead},
		0x64: {0x64, "NOP", 2, 3, ZeroPage, nopRead},
		0x14: {0x14, "NOP", 2, 4, ZeroPageX, nopRead},
		0x34: {0x34, "NOP", 2, 4, ZeroPageX, nopRead},
		0x54: {0x54, "NOP", 2, 4, ZeroPageX, nopRead},
		0x74: {0x74, "NOP", 2, 4, ZeroPageX, nopRead},
		0xd4: {0xd4, "NOP", 2, 4, ZeroPageX, nopRead},
		0xf4: {0xf4, "NOP", 2, 4, ZeroPageX, nopRead},
		0x0c: {0x0c, "NOP", 3, 4, Absolute, nopRead},
		0x1c: {0x1c, "NOP", 3, 4, AbsoluteX, nopRead},
		0x3c: {0x3c, "NOP", 3, 4, AbsoluteX, nopRead},
		0x5c: {0x5c, "NOP", 3, 4, AbsoluteX, nopRead},
		0x7c: {0x7c, "NOP", 3, 4, AbsoluteX, nopRead},
		0xdc: {0xdc, "NOP", 3, 4, AbsoluteX, nopRead},
		0xfc: {0xfc, "NOP", 3, 4, AbsoluteX, nopRead},
		0x1a: {0x1a, "NOP", 1, 2, NoneAddressing, nop},
		0x3a: {0x3a, "NOP", 1, 2, NoneAddressing, nop},
		0x5a: {0x5a, "NOP", 1, 2, NoneAddressing, nop},
		0x7a: {0x7a, "NOP", 1, 2, NoneAddressing, nop},
		0xda: {0xda, "NOP", 1, 2, NoneAddressing, nop},
		0xfa: {0xfa, "NOP", 1, 2, NoneAddressing, nop},

		0x80: {0x80, "NOP", 2, 2, Immediate, nopRead},
		0x82: {0x82, "NOP", 2, 2, Immediate, nopRead},
		0x89: {0x89, "NOP", 2, 2, Immediate, nopRead},
		0xc2: {0xc2, "NOP", 2, 2, Immediate, nopRead},
		0xe2: {0xe2, "NOP", 2, 2, Immediate, nopRead},

		0x02: {0x02, "NOP", 1, 2, NoneAddressing, nopRead},
		0x12: {0x12, "NOP", 1, 2, NoneAddressing, nopRead},
		0x22: {0x22, "NOP", 1, 2, NoneAddressing, nopRead},
		0x32: {0x32, "NOP", 1, 2, NoneAddressing, nopRead},
		0x42: {0x42, "NOP", 1, 2, NoneAddressing, nopRead},
		0x52: {0x52, "NOP", 1, 2, NoneAddressing, nopRead},
		0x62: {0x62, "NOP", 1, 2, NoneAddressing, nopRead},
		0x72: {0x72, "NOP", 1, 2, NoneAddressing, nopRead},
		0x92: {0x92, "NOP", 1, 2, NoneAddressing, nopRead},
		0xb2: {0xb2, "NOP", 1, 2, NoneAddressing, nopRead},
		0xd2: {0xd2, "NOP", 1, 2, NoneAddressing, nopRead},
		0xf2: {0xf2, "NOP", 1, 2, NoneAddressing, nopRead},

		0xa7: {0xa7, "LAX", 2, 3, ZeroPage, lax},
		0xb7: {0xb7, "LAX", 2, 4, ZeroPageY, lax},
		0xaf: {0xaf, "LAX", 3, 4, Absolute, lax},
		0xbf: {0xbf, "LAX", 3, 4, AbsoluteY, lax},
		0xa3: {0xa3, "LAX", 2, 6, IndirectX, lax},
		0xb3: {0xb3, "LAX", 2, 5, IndirectY, lax},

		0x87: {0x87, "SAX", 2, 3, ZeroPage, sax},
		0x97: {0x97, "SAX", 2, 4, ZeroPageY, sax},
		0x8f: {0x8f, "SAX", 3, 4, Absolute, sax},
		0x83: {0x83, "SAX", 2, 6, IndirectX, sax},

		0xeb: {0xeb, "SBC", 2, 2, Immediate, sbc},

		0xc7: {0xc7, "DCP", 2, 5, ZeroPage, dcp},
		0xd7: {0xd7, "DCP", 2, 6, ZeroPageX, dcp},
		0xcf: {0xcf, "DCP", 3, 6, Absolute, dcp},
		0xdf: {0xdf, "DCP", 3, 7, AbsoluteX, dcp},
		0xdb: {0xdb, "DCP", 3, 7, AbsoluteY, dcp},
		0xc3: {0xd3, "DCP", 2, 8, IndirectX, dcp},
		0xd3: {0xc3, "DCP", 2, 8, IndirectY, dcp},

		0xe7: {0xe7, "ISC", 2, 5, ZeroPage, isc},
		0xf7: {0xf7, "ISC", 2, 6, ZeroPageX, isc},
		0xef: {0xef, "ISC", 3, 6, Absolute, isc},
		0xff: {0xff, "ISC", 3, 7, AbsoluteX, isc},
		0xfb: {0xfb, "ISC", 3, 7, AbsoluteY, isc},
		0xe3: {0xf3, "ISC", 2, 8, IndirectX, isc},
		0xf3: {0xe3, "ISC", 2, 8, IndirectY, isc},

		0x07: {0x07, "SLO", 2, 5, ZeroPage, slo},
		0x17: {0x17, "SLO", 2, 6, ZeroPageX, slo},
		0x0f: {0x0f, "SLO", 3, 6, Absolute, slo},
		0x1f: {0x1f, "SLO", 3, 7, AbsoluteX, slo},
		0x1b: {0x1b, "SLO", 3, 7, AbsoluteY, slo},
		0x03: {0x13, "SLO", 2, 8, IndirectX, slo},
		0x13: {0x03, "SLO", 2, 8, IndirectY, slo},

		0x27: {0x27, "RLA", 2, 5, ZeroPage, rla},
		0x37: {0x37, "RLA", 2, 6, ZeroPageX, rla},
		0x2f: {0x2f, "RLA", 3, 6, Absolute, rla},
		0x3f: {0x3f, "RLA", 3, 7, AbsoluteX, rla},
		0x3b: {0x3b, "RLA", 3, 7, AbsoluteY, rla},
		0x23: {0x33, "RLA", 2, 8, IndirectX, rla},
		0x33: {0x23, "RLA", 2, 8, IndirectY, rla},

		0x67: {0x67, "RRA", 2, 5, ZeroPage, rra},
		0x77: {0x77, "RRA", 2, 6, ZeroPageX, rra},
		0x6f: {0x6f, "RRA", 3, 6, Absolute, rra},
		0x7f: {0x7f, "RRA", 3, 7, AbsoluteX, rra},
		0x7b: {0x7b, "RRA", 3, 7, AbsoluteY, rra},
		0x63: {0x73, "RRA", 2, 8, IndirectX, rra},
		0x73: {0x63, "RRA", 2, 8, IndirectY, rra},

		// unofficial and unstable
		0xab: {0xab, "LXA", 2, 3, Immediate, lxa},
		0x8b: {0x8b, "XAA", 2, 3, Immediate, xaa},
		0xbb: {0xbb, "LAS", 3, 2, AbsoluteY, las},
		0x9b: {0x9b, "TAS", 3, 2, AbsoluteY, tas},
		0x93: {0x93, "AHX", 2, 8, IndirectY, ahx},
		0x9f: {0x9f, "AHX", 3, 4, AbsoluteY, ahx},
		0x9e: {0x9e, "SHX", 3, 4, AbsoluteY, shx},
		0x9c: {0x9c, "SHY", 3, 4, AbsoluteX, shy},

		0x0b: {0x0b, "ANC", 2, 2, Immediate, anc},
		0x2b: {0x2b, "ANC", 2, 2, Immediate, anc},
		0x4b: {0x4b, "ALR", 2, 2, Immediate, alr},
		0x6b: {0x6b, "ARR", 2, 2, Immediate, arr},
		0xcb: {0xcb, "AXS", 2, 2, Immediate, axs},
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

func nopRead(p *Processor, op *Instruction) {
	addr, cross := p.getMemoryAddress(op.AddrMode)
	_ = p.readMemUint8(addr)
	if cross {
		p.bus.Tick(1)
	}
}

func nop(_ *Processor, _ *Instruction) {}
