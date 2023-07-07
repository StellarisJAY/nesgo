package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogicalAND(t *testing.T) {
	testCases := []struct {
		name    string
		program []byte
		expectA byte
	}{
		{"im", []byte{LDA_IM, 0b11011, AND_IM, 0b00001}, 1},
		{"zp", []byte{LDA_IM, 0b11011, STA_ZP, 0x2, AND_ZP, 0x2}, 0b11011},
		{"abs", []byte{LDA_IM, 0b11011, STA_ABS, 0x2, 0x30, AND_ABS, 0x2, 0x30}, 0b11011},
	}

	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectA, p.regA)
		})
	}
}

func TestLogicalEOR(t *testing.T) {
	testCases := []struct {
		name    string
		program []byte
		expectA byte
	}{
		{"im", []byte{LDA_IM, 0b11011, EOR_IM, 0b00100}, 0b11111},
		{"zp", []byte{LDA_IM, 0b11011, STA_ZP, 0x2, EOR_ZP, 0x2}, 0},
		{"abs", []byte{LDA_IM, 0b11011, STA_ABS, 0x2, 0x30, EOR_ABS, 0x2, 0x30}, 0},
	}

	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectA, p.regA)
		})
	}
}

func TestLogicalOR(t *testing.T) {
	testCases := []struct {
		name    string
		program []byte
		expectA byte
	}{
		{"im", []byte{LDA_IM, 0b11011, ORA_IM, 0b00100}, 0b11111},
		{"zp", []byte{LDA_IM, 0b11011, LDX_IM, 0b00100, STX_ZP, 0x2, EOR_ZP, 0x2}, 0b11111},
		{"abs", []byte{LDA_IM, 0b11011, LDX_IM, 0b00100, STX_ABS, 0x2, 0x30, EOR_ZP, 0x2, 0x30}, 0b11111},
	}

	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectA, p.regA)
		})
	}
}
