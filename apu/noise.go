package apu

type noise struct {
	enabled           bool
	lengthCounterHalt bool
	constantVolume    bool
	volume            byte

	loopNoise         bool
	noisePeriod       byte
	lengthCounterLoad byte
}

func (n *noise) write(val byte) {
	n.lengthCounterHalt = val&(1<<5) != 0
	n.constantVolume = val&(1<<4) != 0
	n.volume = val & 0b1111
}

func (n *noise) writeNoiseLoop(val byte) {
	n.loopNoise = val&(1<<7) != 0
	n.noisePeriod = val & 0b1111
}

func (n *noise) writeLengthCounter(val byte) {
	n.lengthCounterLoad = val >> 3
}
