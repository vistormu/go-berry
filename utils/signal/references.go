package signal

import (
    "math"
)

type BaseWave struct {
    amp    float64
    freq   float64
    phi    float64
    offset float64
}
func NewBaseWave(amp, freq, phi, offset float64) BaseWave {
    return BaseWave{amp, freq, phi, offset}
}

type Sine struct {
    BaseWave
}
func NewSine(amp, freq, phi, offset float64) Sine {
    return Sine{NewBaseWave(amp, freq, phi, offset)}
}
func (s Sine) Compute(t float64) float64 {
    return s.amp*math.Sin(2*math.Pi*s.freq*t+s.phi) + s.offset
}

type Square struct {
    BaseWave
}
func NewSquare(amp, freq, phi, offset float64) Square {
    return Square{NewBaseWave(amp, freq, phi, offset)}
}
func (s Square) Compute(t float64) float64 {
    return s.amp*math.Copysign(1, math.Sin(2*math.Pi*s.freq*t+s.phi)) + s.offset
}

type Triangular struct {
    BaseWave
}
func NewTriangular(amp, freq, phi, offset float64) Triangular {
    return Triangular{NewBaseWave(amp, freq, phi, offset)}
}
func (tr Triangular) Compute(t float64) float64 {
    return tr.amp*(1+math.Asin(math.Sin(2*math.Pi*tr.freq*t+tr.phi))*2/math.Pi) + tr.offset
}
