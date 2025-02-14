package ml

type Model interface {
    Compute(input []float64) ([]float32, error)
    Close()
}
