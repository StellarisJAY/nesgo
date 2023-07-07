package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJMP(t *testing.T) {
	cases := []struct {
		name    string
		program []byte
		pc      uint16
		x       byte
	}{
		{"jmp1", []byte{INX, JMP_A, 0x05, 0x06, BRK, INX}, 0x0607, 2},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.pc, p.pc)
			assert.Equal(t, tc.x, p.regX)
		})
	}
}

func TestSubroutine(t *testing.T) {
	cases := []struct {
		name    string
		program []byte
		pc      uint16
		x       byte
	}{
		// 调用0x0604的函数，函数在两次INX后返回到BRK位置
		{"call1", []byte{JSR, 0x04, 0x06, BRK, INX, INX, RTS}, 0x0604, 2},
		// 调用0x0606的函数，函数赋值X=10并INX，然后返回到BRK前的INX
		{"call2", []byte{INX, JSR, 0x06, 0x06, INX, BRK, LDX_IM, 10, INX, RTS}, 0x0606, 12},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.pc, p.pc)
			assert.Equal(t, tc.x, p.regX)
		})
	}
}

func TestBranch(t *testing.T) {
	cases := []struct {
		name    string
		program []byte
		pc      uint16
		x       byte
	}{
		{"bcc", []byte{INX, BCC, 0x2, BRK, INX}, 0x0606, 2},
		{"bcs", []byte{SEC, BCS, 0x2, BRK, INX}, 0x0606, 1},
		{"beq", []byte{LDA_IM, 0, BEQ, 0x2, BRK, INX}, 0x0607, 1},
		{"bne", []byte{BNE, 0x2, BRK, INX}, 0x0605, 1},
		{"bmi", []byte{LDA_IM, 0xff, BMI, 0x2, BRK, INX}, 0x0607, 1},
		{"bpl", []byte{BPL, 0x2, BRK, INX}, 0x0605, 1},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.pc, p.pc)
			assert.Equal(t, tc.x, p.regX)
		})
	}
}
