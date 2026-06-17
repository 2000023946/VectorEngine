package vectorengine

import "fmt"

func EuclideanDistance(a []float32, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors must have same dimension")
	}

	var sum float32

	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return sum, nil
}
