package bus

import "log"

type JoyPadButton byte

// JoyPad joy pad
type JoyPad struct {
	strobe       bool         // strobe 该模式下只返回A键的状态
	buttonIdx    byte         // buttonIdx 当前遍历到的button idx
	buttonStatus JoyPadButton // JoyPad 状态，记录哪些键被按下
}

const (
	ButtonA JoyPadButton = 1 << iota
	ButtonB
	Select
	Start
	Up
	Down
	Left
	Right
)

func NewJoyPad() *JoyPad {
	return &JoyPad{false, 0, 0}
}

func (j *JoyPad) write(val byte) {
	log.Println("write joy pad: ", val)
	j.strobe = val&1 == 1
	if j.strobe {
		j.buttonIdx = 0
	}
}

func (j *JoyPad) read() byte {
	res := (j.buttonStatus & (1 << j.buttonIdx)) >> j.buttonIdx
	// strobe模式不遍历idx
	if !j.strobe && j.buttonIdx <= 7 {
		j.buttonIdx += 1
	}
	log.Println("read joy pad: ", res)
	return byte(res)
}

func (j *JoyPad) SetButtonPressed(button JoyPadButton, pressed bool) {
	if pressed {
		j.buttonStatus |= button
	} else {
		j.buttonStatus &= ^button
	}
}
