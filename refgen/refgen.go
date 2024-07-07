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

func (rg RefGen) Compute(t float64) float32 {
    var result float32
    for _, s := range rg.signals {
        result += float32(s.Compute(t))
    }

    return result
}
