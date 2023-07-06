package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncrementRegs(t *testing.T) {
	cases := []struct {
		name    string
		program []byte
		x       byte
		y       byte
	}{
		{"inx", []byte{LDX_IM, 10, INX, INX}, 12, 0},
		{"iny", []byte{LDY_IM, 10, INY, INY}, 0, 12},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.x, p.regX)
			assert.Equal(t, tc.y, p.regY)
		})
	}
}

func TestIncrementMem(t *testing.T) {
	cases := []struct {
		name    string
		program []byte
		addr    uint16
		val     byte
	}{
		{"zp", []byte{LDA_IM, 10, STA_ZP, 0x20, INC_ZP, 0x20, INC_ZP, 0x20, INC_ZP, 0x20}, 0x20, 13},
		{"zpx", []byte{LDA_IM, 10, STA_ZP, 0x30, LDX_IM, 0x10, INC_ZPX, 0x20, INC_ZPX, 0x20, INC_ZPX, 0x20}, 0x30, 13},
		{"abs", []byte{LDA_IM, 10, STA_ZP, 0x20, INC_ABS, 0x20, 0x0, INC_ABS, 0x20, 0x0, INC_ABS, 0x20, 0x0}, 0x20, 13},
		{"abx", []byte{LDA_IM, 10, STA_ZP, 0x30, LDX_IM, 0x10, INC_ABX, 0x20, 0x0, INC_ABX, 0x20, 0x0}, 0x30, 12},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.val, p.readMemUint8(tc.addr))
		})
	}
}
