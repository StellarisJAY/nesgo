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
