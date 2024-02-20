package apu

type triangle struct {
	enabled           bool
	linearCounterLoad byte
	lengthCounterHalt bool
	linearControl     bool

	timerLow          byte
	timerHigh         byte
	lengthCounterLoad byte
}

func (t *triangle) write(val byte) {
	t.linearControl = val&(1<<7) != 0
	t.lengthCounterHalt = t.linearControl
	t.linearCounterLoad = val & 0b111_1111
}

func (t *triangle) writeTimerLow(val byte) {
	t.timerLow = val
}

func (t *triangle) writeTimerHighAndLengthCounter(val byte) {
	t.timerHigh = val & 0b111
	t.lengthCounterLoad = val >> 3
}
