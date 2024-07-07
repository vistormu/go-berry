package refgen

import (
    "math"
)

type Square struct {
    amp float64
    freq float64
    phi float64
    offset float64
}

func NewSquare(amp, freq, phi, offset float64) Square {
    return Square{amp, freq, phi, offset}
}

func (s Square) Compute(t float64) float64 {
    return s.amp*math.Copysign(1, math.Sin(2*math.Pi*s.freq*t+s.phi)) + s.offset
}
