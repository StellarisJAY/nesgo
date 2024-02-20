package apu

type pulse struct {
	enabled           bool
	lengthCounterHalt bool
	duty              byte
	constantVolume    bool
	volume            byte

	timerLow          byte
	timerHigh         byte
	lengthCounterLoad byte

	sweepEnabled bool
	sweepPeriod  byte
	sweepNegate  bool
	sweepShift   byte
}

func (p *pulse) write(val byte) {
	p.constantVolume = val&(1<<4) != 0
	p.volume = val & 0b1111
	p.lengthCounterHalt = val&(1<<5) != 0
	p.duty = val >> 6
}

func (p *pulse) writeSweep(val byte) {
	p.sweepEnabled = val&(1<<7) != 0
	p.sweepShift = val & 0b111
	p.sweepNegate = val&(1<<3) != 0
	p.sweepPeriod = (val >> 4) & 0b111
}

func (p *pulse) writeTimerLow(val byte) {
	p.timerLow = val
}

func (p *pulse) writeTimerHighAndLC(val byte) {
	p.timerHigh = val & 0b111
	p.lengthCounterLoad = val >> 3
}
