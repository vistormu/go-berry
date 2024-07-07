package refgen

import (
    "math"
)


type Sine struct {
    amp float64
    freq float64
    phi float64
    offset float64
}

// result += amp/2 * math.Sin(2*math.Pi*f*t - math.Pi/2) + amp/2

func NewSine(amp, freq, phi, offset float64) Sine {
    return Sine{amp, freq, phi, offset}
}

func (s Sine) Compute(t float64) float64 {
    return s.amp * math.Sin(2*math.Pi*s.freq*t + s.phi) + s.offset
}
