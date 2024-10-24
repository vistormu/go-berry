package algos

type KalmanFilter struct {
    q float64
    r float64
    xHat float64
    p float64
    f float64
    h float64
}
func NewKalmanFilter(processVariance, measurementVariance, initialErrorCovariance, initialEstimate float64) *KalmanFilter {
    return &KalmanFilter{
        q: processVariance,
        r: measurementVariance,
        p: initialErrorCovariance,
        f: 1.0,
        h: 1.0,
        xHat: initialEstimate,
    } 
}
func (kf *KalmanFilter) Compute(measurement float64) float64 {
    xHatPredict := kf.f * kf.xHat
    pPredict := kf.f * kf.p * kf.f + kf.q
    
    k := pPredict * kf.h / (kf.h * pPredict * kf.h + kf.r)
    kf.xHat = xHatPredict + k * (measurement - kf.h * xHatPredict)
    kf.p = (1 - k * kf.h) * pPredict

    return kf.xHat
}
