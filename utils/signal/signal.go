package signal

type Signal interface {
    Compute(a float64) float64
}

type MultiSignal interface {
    Compute(a []float64) []float64
}
