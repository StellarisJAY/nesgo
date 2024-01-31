//go:build sdl

package apu

import (
	"github.com/ebitengine/oto/v3"
	"io"
	"math"
	"sync"
	"time"
)

// WaveFunc 波形函数，输入x、频率、振幅，输出y
// example：正弦波 y = A * sin(freq * 2PI *x)
type WaveFunc func(x, freq, amp float64) float64

// SequenceUnit 输出的音频序列单元，声明一个波形在某段时间内的频率和振幅
type SequenceUnit struct {
	Duration  time.Duration
	Freq      float64 // 频率
	Amplitude float64 // 振幅

	maxSamples int // 该单元的最大sample数量，用于记录单元的结束位置
}

// Sequencer 音频序列器，输出一种波形，外界可以通过Reader读取波形数据
type Sequencer interface {
	// AddSequence 添加需要输出的音频序列
	AddSequence(units ...*SequenceUnit)
	io.Reader
}

type APU struct {
	playerContext *oto.Context
	waves         map[string]Sequencer
	players       []*oto.Player
}

type DefaultSequencer struct {
	sequences  []*SequenceUnit
	sampleRate int
	dt         float64

	sampleNum int
	curSeq    int
	waveFunc  WaveFunc
	m         *sync.Mutex
}

func NewAPU(sampleRate int) (*APU, error) {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 1,
		Format:       oto.FormatFloat32LE,
	})
	if err != nil {
		return nil, err
	}
	<-ready
	return &APU{
		playerContext: ctx,
		players:       make([]*oto.Player, 0),
		waves:         make(map[string]Sequencer),
	}, nil
}

func (a *APU) AddSequencer(name string, w Sequencer) {
	a.waves[name] = w
	player := a.playerContext.NewPlayer(w)
	a.players = append(a.players, player)
}

func (a *APU) Play() {
	go func() {
		for _, player := range a.players {
			player.Play()
		}
	}()
}

func (a *APU) AddSequence(channel string, units ...*SequenceUnit) {
	if ch, ok := a.waves[channel]; ok {
		ch.AddSequence(units...)
	}
}

func NewDefaultSequencer(sampleRate int, waveFunc WaveFunc) *DefaultSequencer {
	dt := 1.0 / float64(sampleRate)
	return &DefaultSequencer{
		sequences:  make([]*SequenceUnit, 0),
		sampleRate: sampleRate,
		dt:         dt,
		waveFunc:   waveFunc,
		m:          &sync.Mutex{},
	}
}

func (w *DefaultSequencer) Read(buf []byte) (int, error) {
	needSamples := len(buf) / 4
	w.m.Lock()
	defer w.m.Unlock()
	// Read EOF
	// 当前没有音频序列需要输出
	if len(w.sequences) == 0 {
		return 0, io.EOF
	}
	x := float64(w.sampleNum) * w.dt
	i := 0
	for ; i < needSamples; i++ {
		// 当前序列单元已结束，切换下一个单元
		if w.sequences[w.curSeq].maxSamples == w.sampleNum {
			if w.curSeq == len(w.sequences)-1 {
				break
			}
			w.curSeq++
		}
		currentSeq := w.sequences[w.curSeq]
		y := w.waveFunc(x, currentSeq.Freq, currentSeq.Amplitude)
		bits := math.Float32bits(float32(y))
		buf[i*4] = byte(bits & 0xff)
		buf[i*4+1] = byte(bits >> 8 & 0xff)
		buf[i*4+2] = byte(bits >> 16 & 0xff)
		buf[i*4+3] = byte(bits >> 24 & 0xff)
		w.sampleNum++
		x += w.dt
	}
	if i == 0 {
		return 0, io.EOF
	}
	return i * 4, nil
}

func (w *DefaultSequencer) AddSequence(units ...*SequenceUnit) {
	w.m.Lock()
	defer w.m.Unlock()
	for _, unit := range units {
		lastSamples := 0
		if len(w.sequences) > 0 {
			lastSamples = w.sequences[len(w.sequences)-1].maxSamples
		}
		w.sequences = append(w.sequences, unit)
		unit.maxSamples = lastSamples + int(unit.Duration.Seconds()*float64(w.sampleRate))
	}
}
