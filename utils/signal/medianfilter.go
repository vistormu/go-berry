package signal

import (
    "slices"
    "github.com/vistormu/go-berry/utils/num"
)

type MedianFilter[T num.Number] struct {
	windowSize int
	values     []T
}

func NewMedianFilter[T num.Number](windowSize int) *MedianFilter[T] {
	return &MedianFilter[T]{
		windowSize: windowSize,
		values:     make([]T, 0, windowSize),
	}
}

func (mf *MedianFilter[T]) Compute(value T) T {
	mf.values = append(mf.values, value)

	if len(mf.values) > mf.windowSize {
		mf.values = mf.values[1:]
	}

	sortedValues := make([]T, len(mf.values))
	copy(sortedValues, mf.values)

    slices.Sort(sortedValues)

	n := len(sortedValues)
	if n%2 == 1 {
		return sortedValues[n/2]
	}
	return (sortedValues[n/2-1] + sortedValues[n/2]) / 2.0
}

type MultiMedianFilter[T num.Number] struct {
    filters []*MedianFilter[T]
}

func NewMultiMedianFilter[T num.Number](windowSize int, numSignals int) *MultiMedianFilter[T] {
    filters := make([]*MedianFilter[T], numSignals)
    for i := 0; i < numSignals; i++ {
        filters[i] = NewMedianFilter[T](windowSize)
    }
    return &MultiMedianFilter[T]{filters: filters}
}

func (mmf *MultiMedianFilter[T]) Compute(values []T) []T {
    results := make([]T, len(values))
    for i, value := range values {
        results[i] = mmf.filters[i].Compute(value)
    }

    return results
}
