package signal

import (
    "github.com/vistormu/go-berry/utils/num"
)

type KalmanFilter[T num.Number] struct {
    q float64
    r float64
    xHat float64
    p float64
    f float64
    h float64
}

func NewKalmanFilter[T num.Number](processVariance, measurementVariance, initialErrorCovariance float64, initialEstimate T) *KalmanFilter[T] {
    return &KalmanFilter[T]{
        q: processVariance,
        r: measurementVariance,
        p: initialErrorCovariance,
        f: 1.0,
        h: 1.0,
        xHat: float64(initialEstimate),
    } 
}

func (kf *KalmanFilter[T]) Compute(measurement T) T {
    xHatPredict := kf.f * kf.xHat
    pPredict := kf.f * kf.p * kf.f + kf.q
    
    k := pPredict * kf.h / (kf.h * pPredict * kf.h + kf.r)
    kf.xHat = xHatPredict + k * (float64(measurement) - kf.h * xHatPredict)
    kf.p = (1 - k * kf.h) * pPredict

    return T(kf.xHat)
}

type MultiKalmanFilter[T num.Number] struct {
	filters []*KalmanFilter[T]
}

func NewMultiKalmanFilter[T num.Number](processVariance, measurementVariance, initialErrorCovariance float64, initialEstimates []T) *MultiKalmanFilter[T] {
    numSignals := len(initialEstimates)
	filters := make([]*KalmanFilter[T], numSignals)
	for i := 0; i < numSignals; i++ {
		filters[i] = NewKalmanFilter(processVariance, measurementVariance, initialErrorCovariance, initialEstimates[i])
	}
	return &MultiKalmanFilter[T]{filters: filters}
}

func (mkf *MultiKalmanFilter[T]) Compute(measurements []T) []T {
	results := make([]T, len(measurements))
	for i, measurement := range measurements {
		results[i] = mkf.filters[i].Compute(measurement)
	}
	return results
}
