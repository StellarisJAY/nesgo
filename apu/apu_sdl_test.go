package apu

import (
	"math"
	"testing"
	"time"
)

func TestAPU_AddWaveChannel(t *testing.T) {
	const sampleRate = 48000
	a, err := NewAPU(sampleRate)
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}

	a.AddSequencer("pulse", NewDefaultSequencer(sampleRate, func(x, freq, amp float64) float64 {
		return amp * (math.Floor(x*freq) - math.Floor(x*freq+0.5))
	}))

	// 一闪一闪亮晶晶
	a.AddSequence("pulse",
		&SequenceUnit{
			Duration:  500 * time.Millisecond,
			Freq:      261.63,
			Amplitude: 1.0,
		},
		&SequenceUnit{
			Duration:  500 * time.Millisecond,
			Freq:      261.63,
			Amplitude: 1.0,
		},
		&SequenceUnit{
			Duration:  500 * time.Millisecond,
			Freq:      392.00,
			Amplitude: 1.0,
		},
		&SequenceUnit{
			Duration:  500 * time.Millisecond,
			Freq:      392.00,
			Amplitude: 1.0,
		},
		&SequenceUnit{
			Duration:  500 * time.Millisecond,
			Freq:      440.00,
			Amplitude: 1.0,
		},
		&SequenceUnit{
			Duration:  500 * time.Millisecond,
			Freq:      440.00,
			Amplitude: 1.0,
		},
		&SequenceUnit{
			Duration:  500 * time.Millisecond,
			Freq:      392.00,
			Amplitude: 1.0,
		})

	a.Play()
	time.Sleep(1 * time.Second)
	// 测试Playing状态下添加新的unit
	a.AddSequence("pulse", &SequenceUnit{
		Duration:  1000 * time.Millisecond,
		Freq:      392.00,
		Amplitude: 1.0,
	})
	time.Sleep(4 * time.Second)
}
