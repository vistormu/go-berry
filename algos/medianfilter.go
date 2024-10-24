package algos

import (
	"sort"
)

type MedianFilter struct {
	windowSize int
	values     []float64
}

func NewMedianFilter(windowSize int) *MedianFilter {
	return &MedianFilter{
		windowSize: windowSize,
		values:     make([]float64, 0, windowSize),
	}
}

func (mf *MedianFilter) Compute(value float64) float64 {
	mf.values = append(mf.values, value)

	if len(mf.values) > mf.windowSize {
		mf.values = mf.values[1:]
	}

	sortedValues := make([]float64, len(mf.values))
	copy(sortedValues, mf.values)
	sort.Float64s(sortedValues)

	n := len(sortedValues)
	if n%2 == 1 {
		return sortedValues[n/2]
	}
	return (sortedValues[n/2-1] + sortedValues[n/2]) / 2.0
}
