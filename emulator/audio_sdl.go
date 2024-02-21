//go:build sdl

package emulator

import (
	"github.com/gordonklaus/portaudio"
	"log"
)

type Audio struct {
	sampleChan     chan float32
	stream         *portaudio.Stream
	outputChannels int
	sampleRate     float64
}

func NewAudio() *Audio {
	sampleChan := make(chan float32, 44100)
	return &Audio{
		sampleChan: sampleChan,
		stream:     nil,
	}
}

func (a *Audio) Start() error {
	api, err := portaudio.DefaultHostApi()
	if err != nil {
		return err
	}
	parameters := portaudio.HighLatencyParameters(nil, api.DefaultOutputDevice)
	log.Println(parameters.SampleRate)
	stream, err := portaudio.OpenStream(parameters, a.callback)
	if err != nil {
		return err
	}
	if err := stream.Start(); err != nil {
		return err
	}
	a.outputChannels = parameters.Output.Channels
	a.sampleRate = parameters.SampleRate
	a.stream = stream
	return nil
}

func (a *Audio) Stop() error {
	return a.stream.Close()
}

func (a *Audio) callback(out []float32) {
	var output float32
	for i := range out {
		// only output to channel0
		if i%a.outputChannels == 0 {
			select {
			case sample := <-a.sampleChan:
				output = sample
			default:
				output = 0
			}
		}
		out[i] = output
	}
}
