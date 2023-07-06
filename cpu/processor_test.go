package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransfer(t *testing.T) {
	cases := []struct {
		name    string
		program []byte
		expectX byte
		expectY byte
		expectA byte
	}{
		{"tax", []byte{LDA_IM, 10, TAX}, 10, 0, 10},
		{"tay", []byte{LDA_IM, 10, TAY}, 0, 10, 10},
		{"txa", []byte{LDX_IM, 10, TXA}, 10, 0, 10},
		{"tya", []byte{LDY_IM, 10, TYA}, 0, 10, 10},
	}

	for _, c := range cases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectX, p.regX)
			assert.Equal(t, tc.expectY, p.regY)
			assert.Equal(t, tc.expectA, p.regA)

		})
	}
}
