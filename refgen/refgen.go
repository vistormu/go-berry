package refgen

type Signal interface {
    Compute(t float64) float64
}

type RefGen struct {
    signals []Signal
}

func NewRefGen(signals []Signal) RefGen {
    return RefGen{signals}
}

func (rg RefGen) Compute(t float64) float64 {
    result := 0.0
    for _, s := range rg.signals {
        result += s.Compute(t)
    }

    return result
}
