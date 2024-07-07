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
    Data []float32
}

func NewWindow(capacity int) *Window {
    return &Window{
        Data: make([]float32, capacity),
    }
}

func (w *Window) Append(element float32) {
    copy(w.Data, w.Data[1:])
    w.Data[len(w.Data)-1] = element
}

type KalmanFilter struct {
    q float32
    r float32
    xHat float32
    p float32
    f float32
    h float32
}

func NewKalmanFilter(processVariance, measurementVariance, initialErrorCovariance float32) *KalmanFilter {
    return &KalmanFilter{
        q: processVariance,
        r: measurementVariance,
        p: initialErrorCovariance,
        f: 1.0,
        h: 1.0,
        xHat: 0.0,
    } 
}

func (kf *KalmanFilter) SetInitialEstimate(value float32) {
    kf.xHat = value
}

func (kf *KalmanFilter) Compute(measurement float32) float32 {
    xHatPredict := kf.f * kf.xHat
    pPredict := kf.f * kf.p * kf.f + kf.q
    
    k := pPredict * kf.h / (kf.h * pPredict * kf.h + kf.r)
    kf.xHat = xHatPredict + k * (measurement - kf.h * xHatPredict)
    kf.p = (1 - k * kf.h) * pPredict

    return kf.xHat
}
