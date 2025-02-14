package signal

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

type MultiKalmanFilter struct {
	filters []*KalmanFilter
}

func NewMultiKalmanFilter(processVariance, measurementVariance, initialErrorCovariance, initialEstimate float64, numSignals int) *MultiKalmanFilter {
	filters := make([]*KalmanFilter, numSignals)
	for i := 0; i < numSignals; i++ {
		filters[i] = NewKalmanFilter(processVariance, measurementVariance, initialErrorCovariance, initialEstimate)
	}
	return &MultiKalmanFilter{filters: filters}
}

func (mkf *MultiKalmanFilter) Compute(measurements []float64) []float64 {
	results := make([]float64, len(measurements))
	for i, measurement := range measurements {
		results[i] = mkf.filters[i].Compute(measurement)
	}
	return results
}
