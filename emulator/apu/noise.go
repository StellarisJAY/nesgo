package apu

// Noise logic:
//
//	Timer --> Shift Register   Length Counter
//	               |                |
//	               v                v
//
// Envelope -------> Gate ----------> Gate --> (to mixer)
type noise struct {
	enabled           bool
	lengthCounterHalt bool
	constantVolume    bool
	volume            byte

	modeFlag          bool
	noisePeriod       byte
	lengthCounterLoad byte

	lengthValue     byte
	envelopeRestart bool
	envelopeVolume  byte
	envelopeValue   byte
	envelopePeriod  byte
	envelopeLoop    bool

	timerValue  uint16
	timerPeriod uint16

	shiftRegister uint16
}

var timerPeriodTable = [16]uint16{
	4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068,
}

func (n *noise) write(val byte) {
	n.lengthCounterHalt = val&(1<<5) != 0
	n.envelopeLoop = n.lengthCounterHalt
	n.constantVolume = val&(1<<4) != 0
	n.volume = val & 0b1111
	n.envelopePeriod = n.volume
}

func (n *noise) writeNoiseLoop(val byte) {
	n.modeFlag = val&(1<<7) != 0
	n.noisePeriod = val & 0b1111
	n.timerPeriod = timerPeriodTable[n.noisePeriod]
}

func (n *noise) writeLengthCounter(val byte) {
	n.lengthCounterLoad = val >> 3
	n.lengthValue = lengthTable[n.lengthCounterLoad]
	n.envelopeRestart = true
}

func (n *noise) stepTimer() {
	if n.timerValue == 0 {
		n.timerValue = n.timerPeriod
		// Feedback is calculated as the exclusive-OR of bit 0 and one other bit:
		// bit 6 if Mode flag is set, otherwise bit 1
		var feedBack uint16
		if n.modeFlag {
			feedBack = (n.shiftRegister & 1) ^ ((n.shiftRegister >> 6) & 1)
		} else {
			feedBack = (n.shiftRegister & 1) ^ ((n.shiftRegister >> 1) & 1)
		}
		// The shift register is shifted right by one bit.
		n.shiftRegister = n.shiftRegister >> 1
		// Bit 14, the leftmost bit, is set to the feedback calculated earlier.
		n.shiftRegister |= (feedBack & 1) << 14
	} else {
		n.timerValue--
	}
}

func (n *noise) stepEnvelope() {
	if n.envelopeRestart {
		n.envelopeVolume = 15
		n.envelopeValue = n.envelopePeriod
		n.envelopeRestart = false
	} else if n.envelopeValue > 0 {
		n.envelopeValue--
	} else {
		if n.envelopeVolume > 0 {
			n.envelopeVolume--
		} else if n.envelopeLoop {
			n.envelopeVolume = 15
		}
		n.envelopeValue = n.envelopePeriod
	}
}

func (n *noise) stepLength() {
	if !n.lengthCounterHalt && n.lengthValue > 0 {
		n.lengthValue--
	}
}

func (n *noise) output() byte {
	if !n.enabled {
		return 0
	}
	//The mixer receives the current envelope volume except when
	//1. Bit 0 of the shift register is set, or
	//2. The length counter is zero
	if n.lengthValue == 0 || n.shiftRegister&1 == 0 {
		return 0
	}
	if n.constantVolume {
		return n.volume
	}
	return n.envelopeVolume
}
