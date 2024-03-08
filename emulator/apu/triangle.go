package apu

// triangle logic:
// Linear Counter   Length Counter
//
//	|                |
//	v                v
//
// Timer ---> Gate ----------> Gate ---> Sequencer ---> (to mixer)
type triangle struct {
	enabled           bool
	linearCounterLoad byte
	lengthCounterHalt bool
	linearControl     bool

	linearCounterReload bool
	linearValue         byte

	lengthValue byte

	timerValue uint16
	dutyValue  byte

	timerLow          byte
	timerHigh         byte
	lengthCounterLoad byte
}

// triangle sends the following looping 32-step sequence of values to the mixer:
var triangleTable = []byte{
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
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
	t.lengthValue = lengthTable[t.lengthCounterLoad]
	t.linearCounterReload = true
}

func (t *triangle) stepTimer() {
	if t.timerValue > 0 {
		t.timerValue--
	} else {
		t.timerValue = t.timerPeriod()
		if t.linearValue > 0 && t.lengthValue > 0 {
			t.dutyValue = (t.dutyValue + 1) % 32
		}
	}
}

func (t *triangle) stepLength() {
	if !t.lengthCounterHalt && t.lengthValue > 0 {
		t.lengthValue--
	}
}

func (t *triangle) stepLinearCounter() {
	// If the linear counter reload flag is set, the linear counter is reloaded with the counter reload value,
	// otherwise if the linear counter is non-zero, it is decremented.

	if t.linearCounterReload {
		t.linearValue = t.linearCounterLoad
	} else if t.linearValue > 0 {
		t.linearValue--
	}
	// If the control flag is clear, the linear counter reload flag is cleared.
	if !t.linearControl {
		t.linearCounterReload = false
	}
}

func (t *triangle) output() byte {
	if !t.enabled {
		return 0
	}
	if t.timerValue < 3 {
		return 0
	}
	// The sequencer is clocked by the timer as long as both the linear counter and the length counter are nonzero.
	if t.lengthValue == 0 || t.linearValue == 0 {
		return 0
	}
	return triangleTable[t.dutyValue]
}

func (t *triangle) timerPeriod() uint16 {
	return uint16(t.timerLow) | uint16(t.timerHigh)>>8
}
