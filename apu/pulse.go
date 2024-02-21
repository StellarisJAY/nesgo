package apu

// Pulse logic:
//
//	                 Sweep -----> Timer
//	                   |            |
//	                   |            |
//	                   |            v
//	                   |        Sequencer   Length Counter
//	                   |            |             |
//	                   |            |             |
//	                   v            v             v
//	Envelope -------> Gate -----> Gate -------> Gate --->(to mixer)
type pulse struct {
	enabled           bool
	lengthCounterHalt bool
	duty              byte
	dutyValue         byte
	constantVolume    bool
	volume            byte

	startEnvelope  bool
	envelopeValue  byte
	envelopeVolume byte
	envelopePeriod byte
	envelopeLoop   bool

	timerLow    byte
	timerHigh   byte
	timerValue  uint16
	lengthValue byte

	sweepReload  bool
	sweepEnabled bool
	sweepPeriod  byte
	sweepNegate  bool
	sweepShift   byte
	sweepValue   byte
}

var dutyTable = [4][8]byte{
	{0, 1, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 0, 0, 0, 0, 0},
	{0, 1, 1, 1, 1, 0, 0, 0},
	{1, 0, 0, 1, 1, 1, 1, 1},
}

func (p *pulse) write(val byte) {
	p.constantVolume = val&(1<<4) != 0
	p.volume = val & 0b1111
	p.envelopePeriod = p.volume
	p.lengthCounterHalt = val&(1<<5) != 0
	p.envelopeLoop = p.lengthCounterHalt
	p.duty = val >> 6
	p.startEnvelope = true
}

func (p *pulse) writeSweep(val byte) {
	p.sweepEnabled = val&(1<<7) != 0
	p.sweepShift = val & 0b111
	p.sweepNegate = val&(1<<3) != 0
	p.sweepPeriod = (val >> 4) & 0b111
	p.sweepReload = true
}

func (p *pulse) writeTimerLow(val byte) {
	p.timerLow = val
}

func (p *pulse) writeTimerHighAndLC(val byte) {
	p.timerHigh = val & 0b111
	p.lengthValue = lengthTable[val>>3]
	p.dutyValue = 0
	p.startEnvelope = true
}

func (p *pulse) stepEnvelope() {
	if p.startEnvelope { // envelope周期开始， 初始化volume
		p.envelopeVolume = 15
		p.envelopeValue = p.envelopePeriod
		p.startEnvelope = false
	} else if p.envelopeValue > 0 { // 每次step减小value直到0
		p.envelopeValue--
	} else { // 周期结束
		if p.envelopeVolume > 0 { // 减小volume
			p.envelopeVolume--
		} else if p.envelopeLoop {
			p.envelopeVolume = 15
		}
		// 下一个周期
		p.envelopeValue = p.envelopePeriod
	}
}

func (p *pulse) stepSweep() {
	if p.sweepReload {
		p.sweepValue = p.sweepPeriod
		p.sweepReload = false
	} else if p.sweepValue > 0 {
		p.sweepValue--
	} else {
		if p.sweepReload {
			if p.sweepEnabled {
				p.sweep()
			}
		}
		p.sweepValue = p.sweepPeriod
	}
}

func (p *pulse) stepLength() {
	if !p.lengthCounterHalt && p.lengthValue > 0 {
		p.lengthValue--
	}
}

func (p *pulse) stepTimer() {
	if p.timerValue == 0 {
		p.timerValue = p.timerPeriod()
		p.dutyValue = (p.dutyValue + 1) % 8
	} else {
		p.timerValue--
	}
}

func (p *pulse) sweep() {
	period := p.timerPeriod()
	change := period >> p.sweepShift
	if p.sweepNegate {
		p.setTimerPeriod(period - change)
	} else {
		p.setTimerPeriod(period + change)
	}
}

func (p *pulse) timerPeriod() uint16 {
	return uint16(p.timerLow) | uint16(p.timerHigh)>>8
}

func (p *pulse) setTimerPeriod(val uint16) {
	p.timerLow = byte(val & 0xff)
	p.timerHigh = byte(val >> 8)
}

func (p *pulse) output() byte {
	if !p.enabled {
		return 0
	}
	if p.lengthValue == 0 {
		return 0
	}
	if dutyTable[p.duty][p.dutyValue] == 0 {
		return 0
	}
	if p.constantVolume {
		return p.volume
	} else {
		return p.envelopeVolume
	}
}
