package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCMP(t *testing.T) {
	cases := []struct {
		name    string
		program []byte
		pc      uint16
		x       byte
	}{
		{"im_eq", []byte{LDA_IM, 20, CMP_IM, 20, BEQ, 0x1, BRK, INX, INX}, 0x060a, 2},
		{"im_gt", []byte{LDA_IM, 20, CMP_IM, 10, BEQ, 0x3, BCS, 0x1, BRK, INX, INX}, 0x060c, 2},
		{"abs_eq", []byte{LDA_IM, 20, LDY_IM, 20, STY_ABS, 0x0, 0x1, CMP_ABS, 0x0, 0x1, BEQ, 0x1, BRK, INX, INX}, 0x0610, 2},
		{"abs_gt", []byte{LDA_IM, 20, LDY_IM, 10, STY_ABS, 0x0, 0x1, CMP_ABS, 0x0, 0x1, BEQ, 0x3, BCS, 0x1, BRK, INX, INX}, 0x0612, 2},
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
