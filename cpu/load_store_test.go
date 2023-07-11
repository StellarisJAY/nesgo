package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSTA(t *testing.T) {
	cases := []struct {
		name      string
		program   []byte
		memAddr   uint16
		expectVal byte
	}{
		{"zp", []byte{LDA_IM, 10, STA_ZP, 0x20}, 0x20, 10},
		{"abs", []byte{LDA_IM, 10, STA_ABS, 0x0}, 0x0, 10},
		{"abx", []byte{LDA_IM, 0x60, TAX, LDA_IM, 10, STA_ABX, 0x0}, 0x60, 10},
		{"ix", []byte{LDA_IM, 0x30, STA_ZP, 0x20, LDA_IM, 10, STA_IX, 0x20}, 0x30, 10},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectVal, p.readMemUint8(tc.memAddr))
		})
	}
}

func TestSTX(t *testing.T) {
	cases := []struct {
		name      string
		program   []byte
		memAddr   uint16
		expectVal byte
	}{
		{"zp", []byte{LDA_IM, 10, TAX, STX_ZP, 0x20}, 0x20, 10},
		{"zpy", []byte{LDA_IM, 10, TAX, STX_ZPY, 0x20}, 0x20, 10},
		{"abs", []byte{LDA_IM, 10, TAX, STX_ABS, 0x30}, 0x30, 10},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectVal, p.readMemUint8(tc.memAddr))
		})
	}
}

func TestSTY(t *testing.T) {
	cases := []struct {
		name      string
		program   []byte
		memAddr   uint16
		expectVal byte
	}{
		{"zp", []byte{LDA_IM, 10, TAY, STY_ZP, 0x20}, 0x20, 10},
		{"zpy", []byte{LDA_IM, 10, TAY, STY_ZPX, 0x20}, 0x20, 10},
		{"abs", []byte{LDA_IM, 10, TAY, STY_ABS, 0x30}, 0x30, 10},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectVal, p.readMemUint8(tc.memAddr))
		})
	}
}

func TestLDX(t *testing.T) {
	cases := []struct {
		name      string
		program   []byte
		expectVal byte
	}{
		{"im", []byte{LDX_IM, 10}, 10},
		{"zp", []byte{LDA_IM, 10, STA_ZP, 0x20, LDX_ZP, 0x20}, 10},
		{"zpy", []byte{LDA_IM, 10, STA_ZP, 0x20, LDX_ZPY, 0x20}, 10},
		{"abs", []byte{LDA_IM, 10, STA_ZP, 0x30, LDX_ABS, 0x30}, 10},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectVal, p.regX)
		})
	}
}

func TestLDY(t *testing.T) {
	cases := []struct {
		name      string
		program   []byte
		expectVal byte
	}{
		{"im", []byte{LDY_IM, 10}, 10},
		{"zp", []byte{LDA_IM, 10, STA_ZP, 0x20, LDY_ZP, 0x20}, 10},
		{"zpy", []byte{LDA_IM, 10, STA_ZP, 0x20, LDY_ZPX, 0x20}, 10},
		{"abs", []byte{LDA_IM, 10, STA_ZP, 0x30, LDY_ABS, 0x30}, 10},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectVal, p.regY)
		})
	}
}

func TestLDA(t *testing.T) {
	cases := []struct {
		name         string
		program      []byte
		expectA      byte
		expectStatus byte
	}{
		{"normal", []byte{LDA_IM, 10}, 10, BreakStatus},
		{"zero", []byte{LDA_IM, 10, LDA_IM, 0}, 0, ZeroStatus},
		{"negative", []byte{LDA_IM, 10, LDA_IM, 0xff}, 0xFF, NegativeStatus},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectA, p.regA)
			assert.True(t, p.regStatus&tc.expectStatus != 0)
		})
	}
}
