package main

import (
	"errors"
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type SineWave struct {
	sampleFactor float64 // Just for ease of use so that we don't have to calculate every sample
	phase        float64
}

func (g *SineWave) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples { // increment = ((2 * PI) / SampleRate) * freq
		v := math.Sin(g.phase * 2.0 * math.Pi) // period of the wave is thus defined as: 2 * PI.
		samples[i][0] = v
		samples[i][1] = v
		_, g.phase = math.Modf(g.phase + g.sampleFactor)
	}

	return len(samples), true
}

func (*SineWave) Err() error {
	return nil
}

func SineTone(sr beep.SampleRate, freq float64) (beep.Streamer, error) {
	dt := freq / float64(sr)

	if dt >= 1.0/2.0 {
		return nil, errors.New("faiface sine tone generator: samplerate must be at least 2 times grater then frequency")
	}

	return &SineWave{dt, 0.1}, nil
}

func main() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10)) // sr.N(time.Second/10) = buffer size for duration 1/10 second
	sine, _ := SineTone(sr, 100)
	speaker.Play(sine)
	select {} // makes the program hang forever
}
