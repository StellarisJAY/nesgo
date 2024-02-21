package apu

var (
	pulseTable = [32]float32{}
	tndTable   [203]float32
)

var lengthTable = []byte{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

type BasicAPU struct {
	muted            bool
	cycles           int
	frameCounterMode byte
	irqInhibitFlag   bool

	p1 *pulse
	p2 *pulse
	t  *triangle
	n  *noise
	d  *dmc

	outputChan       chan float32
	sampleRate       float64
	frameCounterRate int
}

func init() {
	for i := 0; i < 31; i++ {
		pulseTable[i] = 95.52 / (8128.0/float32(i) + 100)
	}
	for i := 0; i < 203; i++ {
		tndTable[i] = 163.67 / (24329.0/float32(i) + 100)
	}
}

func NewBasicAPU() *BasicAPU {
	return &BasicAPU{
		p1: &pulse{},
		p2: &pulse{},
		t:  &triangle{},
		n:  &noise{shiftRegister: 1}, //On power-up, the shift register is loaded with the value 1.
		d:  &dmc{},
	}
}

func (a *BasicAPU) SetMemReader(reader func(addr uint16) byte) {
	a.d.memReader = reader
}

func (a *BasicAPU) SetOutputChan(outChan chan float32) {
	a.outputChan = outChan
}

func (a *BasicAPU) SetRates(cpuFrequency int, sampleRate float64) {
	// cpu cycles per sample
	a.sampleRate = float64(cpuFrequency) / sampleRate
	a.frameCounterRate = cpuFrequency / 240
}

func (a *BasicAPU) Mute() {
	a.muted = true
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
		a.writeFrameCounter(data)
	}
}

func (a *BasicAPU) Tick() {
	if a.muted {
		return
	}
	oldCycles := a.cycles
	a.cycles++
	cycles := a.cycles
	a.stepTimer()
	f1, f2 := oldCycles/a.frameCounterRate, cycles/a.frameCounterRate
	if f1 != f2 {
		a.stepFrameCounter()
	}
	c1, c2 := int(float64(oldCycles)/a.sampleRate), int(float64(cycles)/a.sampleRate)
	if c1 != c2 {
		select {
		case a.outputChan <- a.Output():
		default:
		}
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
		frameVal := (a.cycles / a.frameCounterRate) % 4
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
		frameVal := (a.cycles / a.frameCounterRate) % 5
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

func (a *BasicAPU) stepTimer() {
	if a.cycles&1 == 0 {
		a.p1.stepTimer()
		a.p2.stepTimer()
	}
	// Unlike the pulse channels,
	// this timer ticks at the rate of the CPU clock rather than the APU (CPU/2) clock.
	a.t.stepTimer()
	a.n.stepTimer()
	a.d.stepTimer()
}

func (a *BasicAPU) stepEnvelope() {
	a.p1.stepEnvelope()
	a.p2.stepEnvelope()
	a.t.stepLinearCounter()
	a.n.stepEnvelope()
}

func (a *BasicAPU) stepSweep() {
	a.p1.stepSweep()
	a.p2.stepSweep()
}

func (a *BasicAPU) stepLengthCounter() {
	a.p1.stepLength()
	a.p2.stepLength()
	a.t.stepLength()
	a.n.stepLength()
}

func (a *BasicAPU) sendIRQ() {

}

func (a *BasicAPU) Output() float32 {
	pout1 := a.p1.output()
	pout2 := a.p2.output()
	pulseOut := pulseTable[pout1+pout2]
	t := a.t.output()
	d := a.d.output()
	n := a.n.output()
	tndOutput := tndTable[3*t+2*d+n]
	return pulseOut + tndOutput
}
