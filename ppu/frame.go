package ppu

type Frame struct {
	data           []byte
	compressBuffer []byte
}

const (
	WIDTH  = 256
	HEIGHT = 240
)

func NewFrame() *Frame {
	return &Frame{
		make([]byte, WIDTH*HEIGHT*3),
		make([]byte, 0, 62208),
	}
}

func (f *Frame) setPixel(x, y uint32, color Color) {
	first := y*3*WIDTH + x*3
	if first+2 < uint32(len(f.data)) {
		f.data[first] = color.R
		f.data[first+1] = color.G
		f.data[first+2] = color.B
	}
}

// 压缩帧数据，原本每个像素为3字节，压缩后可变为1字节
//
//		[w*h pixels][colors]
//		每个像素用一个字节表示颜色编号，一帧画面可以有256种颜色
//		每种颜色以三字节记录在frame header
//		原帧大小 = width*height*3 = 184320
//		最大压缩帧大小 = 256*3 + width*height = 62208
//	    最小压缩率 = 33.7%
func (f *Frame) compressedFrameData() []byte {
	colors := make(map[uint32]byte)
	colorVals := make([]byte, 3*256)
	var nextColorId byte = 0
	for i := 0; i < len(f.data); i += 3 {
		rgb := (uint32(f.data[i]) << 16) | (uint32(f.data[i+1]) << 8) | uint32(f.data[i+2])
		if id, ok := colors[rgb]; !ok {
			colors[rgb] = nextColorId
			f.compressBuffer = append(f.compressBuffer, nextColorId)
			colorVals[nextColorId*3] = f.data[i]
			colorVals[nextColorId*3+1] = f.data[i+1]
			colorVals[nextColorId*3+2] = f.data[i+2]
			nextColorId++
		} else {
			f.compressBuffer = append(f.compressBuffer, id)
		}
	}
	bufLen := WIDTH*HEIGHT + len(colors)*3
	copy(f.compressBuffer[WIDTH*HEIGHT:bufLen], colorVals[:len(colors)*3])
	result := f.compressBuffer[:bufLen]
	f.compressBuffer = f.compressBuffer[:0]
	return result
}

func decompressFrame(data []byte) []byte {
	pixels := data[:WIDTH*HEIGHT]
	colors := data[WIDTH*HEIGHT:]
	frame := make([]byte, WIDTH*HEIGHT*3)
	for i, p := range pixels {
		frame[i*3] = colors[p*3]
		frame[i*3+1] = colors[p*3+1]
		frame[i*3+2] = colors[p*3+2]
	}
	return frame
}

func (f *Frame) Data() []byte {
	return f.data
}
