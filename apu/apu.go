package apu

const (
	frameCounterRate = 1790000 / 240
)

var (
	pulseTable = [32]float32{}
)

type BasicAPU struct {
	cycles           int
	frameCounterMode byte
	irqInhibitFlag   bool

	p1 *pulse
	p2 *pulse
	t  *triangle
	n  *noise
	d  *dmc
}

func NewBasicAPU() *BasicAPU {
	return &BasicAPU{
		p1: &pulse{},
		p2: &pulse{},
		t:  &triangle{},
		n:  &noise{},
		d:  &dmc{},
	}
}

func (a *BasicAPU) Write(addr uint16, data byte) {
	switch addr {
	case 0x4000: // pulse1 ddlcvvvv
		a.p1.write(data)
	case 0x4004: // pulse2 ddlcvvvv
		a.p2.write(data)
	case 0x4001:
		a.p1.writeSweep(data)
	case 0x4005:
		a.p2.writeSweep(data)
	case 0x4002: // pulse1 timer low
		a.p1.writeTimerLow(data)
	case 0x4006: // pulse2 timer low
		a.p2.writeTimerLow(data)
	case 0x4003: // pulse1 timer high and lcl
		a.p1.writeTimerHighAndLC(data)
	case 0x4007: // pulse2 timer high and lcl
		a.p2.writeTimerHighAndLC(data)
	case 0x4008:
		a.t.write(data)
	case 0x4009: // unused
	case 0x400A:
		a.t.writeTimerLow(data)
	case 0x400B:
		a.t.writeTimerHighAndLengthCounter(data)
	case 0x400C:
		a.n.write(data)
	case 0x400D: // unused
	case 0x400E:
		a.n.writeNoiseLoop(data)
	case 0x400F:
		a.n.writeLengthCounter(data)
	case 0x4010:
		a.d.write(data)
	case 0x4011:
		a.d.writeLoadCounter(data)
	case 0x4012:
		a.d.writeSampleAddress(data)
	case 0x4013:
		a.d.writeSampleLength(data)
	case 0x4015:
		a.writeStatus(data)
	case 0x4017:

	}
}

func (a *BasicAPU) Tick() {
	oldCycles := a.cycles
	a.cycles++
	cycles := a.cycles
	f1, f2 := oldCycles/frameCounterRate, cycles/frameCounterRate
	if f1 != f2 {
		a.stepFrameCounter()
	}
}

func (a *BasicAPU) writeStatus(val byte) {
	a.p1.enabled = val&1 != 0
	a.p2.enabled = val&2 != 0
	a.t.enabled = val&4 != 0
	a.n.enabled = val&8 != 0
	a.d.enabled = val&16 != 0
}

func (a *BasicAPU) writeFrameCounter(val byte) {
	if val&(1<<7) != 0 {
		a.frameCounterMode = 5
	} else {
		a.frameCounterMode = 4
	}
	a.irqInhibitFlag = val&(1<<6) != 0
}

// mode 4 或 mode 5 的 frameCounter
// See: https://www.nesdev.org/wiki/APU_Frame_Counter
func (a *BasicAPU) stepFrameCounter() {
	if a.frameCounterMode == 4 {
		frameVal := (a.cycles / frameCounterRate) % 4
		switch frameVal {
		case 0:
			a.stepEnvelope()
		case 1:
			a.stepEnvelope()
			a.stepSweep()
			a.stepLengthCounter()
		case 2:
			a.stepEnvelope()
		case 3:
			a.stepEnvelope()
			a.stepSweep()
			a.stepLengthCounter()
			a.sendIRQ()
		}
		return
	}
	if a.frameCounterMode == 5 {
		frameVal := (a.cycles / frameCounterRate) % 5
		switch frameVal {
		case 0:
			a.stepEnvelope()
		case 1:
			a.stepEnvelope()
			a.stepSweep()
			a.stepLengthCounter()
		case 2:
			a.stepEnvelope()
		case 3:
		case 4:
			a.stepEnvelope()
			a.stepSweep()
			a.stepLengthCounter()
		}
		return
	}
}

func (a *BasicAPU) stepEnvelope() {
	// step pulse1, pulse2, noise, triangle linear counter
}

func (a *BasicAPU) stepSweep() {

}

func (a *BasicAPU) stepLengthCounter() {

}

func (a *BasicAPU) sendIRQ() {

}

func (a *BasicAPU) Output() {

}
