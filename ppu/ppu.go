package ppu

// ReadAndUpdateScreen 从内存读取每个cell，并在frame中修改cell的值，如果有更新则通知给渲染器
func ReadAndUpdateScreen(memory []byte, frame []byte) (updated bool) {
	idx := 0
	for _, cell := range memory {
		color := getRGBAColor(cell)
		r, g, b := color.R, color.G, color.B
		if frame[idx] != r || frame[idx+1] != g || frame[idx+2] != b {
			frame[idx] = r
			frame[idx+1] = g
			frame[idx+2] = b
			updated = true
		}
		idx += 3
	}
	return
}
