package apu

type dmc struct {
	enabled   bool
	irqEnable bool
	loop      bool
	frequency byte

	loadCounter   byte
	sampleAddress byte
	sampleLength  byte
}

func (d *dmc) write(val byte) {
	d.irqEnable = val&(1<<7) != 0
	d.loop = val&(1<<6) != 0
	d.frequency = val & 0b1111
}

func (d *dmc) writeLoadCounter(val byte) {
	d.loadCounter = val & 0b111_1111
}

func (d *dmc) writeSampleAddress(val byte) {
	d.sampleAddress = val
}

func (d *dmc) writeSampleLength(val byte) {
	d.sampleLength = val
}
