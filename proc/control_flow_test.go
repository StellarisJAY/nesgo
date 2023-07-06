package proc

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
		y       byte
	}{
		{"jmp1", []byte{INX, JMP_A, 0x05, 0x80, BRK, INX}, 0x8007, 2, 0},
		{"jmp2", []byte{INY, JMP_A, 0x05, 0x80, BRK, INY}, 0x8007, 0, 2},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.pc, p.pc)
		})
	}
}
