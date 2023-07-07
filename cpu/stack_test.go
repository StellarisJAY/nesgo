package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStackPush(t *testing.T) {
	testCases := []struct {
		name       string
		program    []byte
		expectSp   byte
		expectPeek byte
	}{
		{"pha", []byte{LDA_IM, 10, PHA}, 0xFF - 1, 10},
		{"pha2", []byte{LDA_IM, 10, PHA, LDA_IM, 12, PHA}, 0xFF - 2, 12},
		{"php", []byte{LDA_IM, 0xff, PHP}, 0xFF - 1, NegativeStatus},
		{"php2", []byte{LDA_IM, 0xff, PHP, LDA_IM, 0, PHP}, 0xFF - 2, ZeroStatus},
	}

	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectSp, p.sp)
			assert.Equal(t, tc.expectPeek, p.peek())
		})
	}
}

func TestStackPop(t *testing.T) {
	testCases := []struct {
		name       string
		program    []byte
		expectSp   byte
		expectA    byte
		expectStat byte
	}{
		{"pla", []byte{LDA_IM, 10, PHA, PLA}, 0xFF, 10, 0},
		{"pla2", []byte{LDA_IM, 10, PHA, LDA_IM, 0, PHA, PLA}, 0xFF - 1, 0, ZeroStatus},
		{"plp", []byte{LDA_IM, 0xff, PHP, PLP}, 0xFF, 0xff, NegativeStatus},
		{"plp2", []byte{LDA_IM, 0xff, PHP, LDA_IM, 0, PHP, PLP}, 0xFF - 1, 0, ZeroStatus},
	}

	for _, c := range testCases {
		tc := c
		t.Run(tc.name, func(t *testing.T) {
			p := NewProcessor()
			p.LoadAndRun(tc.program)
			assert.Equal(t, tc.expectSp, p.sp)
			assert.Equal(t, tc.expectA, p.regA)
			assert.Equal(t, tc.expectStat, p.regStatus)
		})
	}
}
