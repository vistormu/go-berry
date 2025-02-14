package num

import (
    "math"
)


func Mean[T Number](values []T) T {
    var sum T
    length := len(values)

    for _, v := range values {
        sum += v
    }

    return sum / T(length)
}

func MultiMean[T Number](values [][]T) []T {
	if len(values) == 0 {
		return nil
	}

	// Get the number of indices in each inner slice.
	numIndices := len(values[0])
	means := make([]T, numIndices)

	// Iterate over each index position.
	for i := 0; i < numIndices; i++ {
		var sum T
		// Sum the i-th element from each inner slice.
		for _, row := range values {
			if len(row) != numIndices {
				panic("all inner slices must have the same length")
			}
			sum += row[i]
		}
		// Compute the mean for index i.
		means[i] = sum / T(len(values))
	}

	return means
}

func StdDev[T Number](values []T) T {
	if len(values) == 0 {
		return 0
	}

	var sum T
	for _, v := range values {
		sum += v
	}

	mean := sum / T(len(values))

	var varianceSum T
	for _, v := range values {
		diff := v - mean
		varianceSum += diff * diff
	}

	variance := varianceSum / T(len(values))
	return T(math.Sqrt(float64(variance)))
}
