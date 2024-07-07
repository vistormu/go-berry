package refgen

type Step struct {
    amp float64
}

func NewStep(amp float64) Step {
    return Step{amp}
}

func (s Step) Compute(t float64) float64 {
    return s.amp
}
