package utils


type Number interface {
    ~float64 | ~float32 | ~int
}

func Clip[T Number](value, min, max T) T {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

type Window struct {
    Data []float64
}

func NewWindow(capacity int) *Window {
    return &Window{
        Data: make([]float64, capacity),
    }
}

func (w *Window) Append(element float64) {
    copy(w.Data, w.Data[1:])
    w.Data[len(w.Data)-1] = element
}

type KalmanFilter struct {
    q float64
    r float64
    xHat float64
    p float64
    f float64
    h float64
}

func NewKalmanFilter(processVariance, measurementVariance, initialErrorCovariance float64) *KalmanFilter {
    return &KalmanFilter{
        q: processVariance,
        r: measurementVariance,
        p: initialErrorCovariance,
        f: 1.0,
        h: 1.0,
        xHat: 0.0,
    } 
}

func (kf *KalmanFilter) SetInitialEstimate(value float64) {
    kf.xHat = value
}

func (kf *KalmanFilter) Compute(measurement float64) float64 {
    xHatPredict := kf.f * kf.xHat
    pPredict := kf.f * kf.p * kf.f + kf.q
    
    k := pPredict * kf.h / (kf.h * pPredict * kf.h + kf.r)
    kf.xHat = xHatPredict + k * (measurement - kf.h * xHatPredict)
    kf.p = (1 - k * kf.h) * pPredict

    return kf.xHat
}
