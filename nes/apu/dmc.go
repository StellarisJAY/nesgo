package apu

// DMC Logic:
//
//	Timer
//	 |
//	 v
//
// Reader ---> Buffer ---> Shifter ---> Output level ---> (to the mixer)
type dmc struct {
	enabled       bool
	irqEnable     bool
	loop          bool
	rateIndex     byte
	sampleAddress uint16
	sampleLength  uint16

	sampleBuffer   byte
	sampleCounter  byte
	addressCounter uint16
	bytesRemaining uint16

	timerPeriod uint16
	timerValue  uint16

	bitsRemaining byte
	shiftRegister byte

	value byte

	memReader func(addr uint16) byte
}

var dmcPeriodTable = [16]uint16{428, 380, 340, 320, 286, 254, 226, 214, 190, 160, 142, 128, 106, 84, 72, 54}

func (d *dmc) write(val byte) {
	d.irqEnable = val&(1<<7) != 0
	d.loop = val&(1<<6) != 0
	d.rateIndex = val & 0b1111
	d.timerPeriod = dmcPeriodTable[d.rateIndex]
}

func (d *dmc) writeLoadCounter(val byte) {
	d.value = val & 0b111_1111
}

func (d *dmc) writeSampleAddress(val byte) {
	// Sample address =  $C000 + (A * 64)
	d.sampleAddress = 0xC000 + uint16(val)<<6
}

func (d *dmc) writeSampleLength(val byte) {
	// Sample length = (L * 16) + 1 bytes
	d.sampleLength = uint16(val)<<4 + 1
}

func (d *dmc) stepTimer() {
	if !d.enabled {
		return
	}
	d.stepReader()
	if d.timerValue == 0 {
		d.timerValue = d.timerPeriod
		d.stepShifter()
	} else {
		d.timerValue--
	}
}

func (d *dmc) restart() {
	d.addressCounter = d.sampleAddress
	d.bytesRemaining = d.sampleLength
}

func (d *dmc) stepReader() {
	// shiftRegister already consumed all bits
	// read a new byte and start a new shift period
	if d.bytesRemaining > 0 && d.bitsRemaining == 0 {
		d.bitsRemaining = 8
		d.shiftRegister = d.memReader(d.addressCounter)
		d.addressCounter++
		// wrapped around to 0x8000
		if d.addressCounter == 0 {
			d.addressCounter = 0x8000
		}
		d.bytesRemaining--
		if d.bytesRemaining == 0 && d.loop {
			d.restart()
		}
	}
}

// shiftRegister bit 0 determines output level +2 or -2
// right shift register on every step
func (d *dmc) stepShifter() {
	if d.bitsRemaining == 0 {
		return
	}
	if d.shiftRegister&1 == 1 {
		// if adding or subtracting 2 would cause the output level to leave the 0-127 range,
		// leave the output level unchanged.
		if d.value < 125 {
			d.value += 2
		}
	} else {
		if d.value > 2 {
			d.value -= 2
		}
	}
	d.shiftRegister >>= 1
	d.bitsRemaining--
}

func (d *dmc) output() byte {
	return d.value
}
