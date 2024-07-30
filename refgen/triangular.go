package refgen

import (
    "math"
)

type Triangular struct {
    amp float64
    freq float64
    phi float64
    offset float64
}

func NewTriangular(amp, freq, phi, offset float64) Triangular {
    return Triangular{amp, freq, phi, offset}
}

func (tr Triangular) Compute(t float64) float64 {
    return tr.amp*(1+math.Asin(math.Sin(2*math.Pi*tr.freq*t + tr.phi))*2/math.Pi) + tr.offset
}
