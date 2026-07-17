package db

import (
	"errors"
	"sort"
)

// SearchResult holds the matched vector and its distance score
type SearchResult struct {
	Vector   Vector
	Distance float32
}

// Search performs a brute-force linear scan to find the top K nearest neighbors
func (e *VectorEngine) Search(query []float32, k int) ([]SearchResult, error) {
	if len(query) != e.dim {
		return nil, errors.New("query dimension mismatch")
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	// 1. Calculate distances to all vectors
	results := make([]SearchResult, len(e.vectors))
	for i, v := range e.vectors {
		results[i] = SearchResult{
			Vector:   v,
			Distance: squaredEuclidean(query, v.Data),
		}
	}

	// 2. Sort by distance (ascending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	// 3. Return top K
	if k > len(results) {
		k = len(results)
	}
	return results[:k], nil
}

// squaredEuclidean computes the distance without the math.Sqrt() overhead.
// Preserves the correct sorting order while saving CPU cycles.
func squaredEuclidean(a, b []float32) float32 {
	var sum float32 = 0.0
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return sum
}
