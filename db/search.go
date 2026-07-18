package db

import (
	"math"
)

// Search routes the query based on the engine's current lifecycle phase
func (e *VectorEngine) Search(query []float32) ([]float32, float32) {
	if e.Phase == PhaseIndexed {
		// 1. Find nearest frozen centroid
		cIdx := e.findNearestCentroid(query)

		// 2. Scan only that specific contiguous bucket (Cache friendly)
		return findNearestInList(query, e.Buckets[cIdx])
	}

	// Fallback: Brute force the raw list if still warming up
	return findNearestInList(query, e.RawVectors)
}

// findNearestCentroid is used for both routing inserts and routing searches
func (e *VectorEngine) findNearestCentroid(vec []float32) int {
	bestIdx := -1
	minDist := float32(math.MaxFloat32)
	for i, centroid := range e.Centroids {
		dist := euclideanSq(vec, centroid)
		if dist < minDist {
			minDist = dist
			bestIdx = i
		}
	}
	return bestIdx
}

// findNearestInList powers both the Phase 1 scan and the Phase 2 bucket scan
func findNearestInList(query []float32, list [][]float32) ([]float32, float32) {
	if len(list) == 0 {
		return nil, -1 // Edge case: empty list
	}

	var bestVec []float32
	minDist := float32(math.MaxFloat32)
	for _, v := range list {
		dist := euclideanSq(query, v)
		if dist < minDist {
			minDist = dist
			bestVec = v
		}
	}
	return bestVec, minDist
}

// euclideanSq calculates distance without the expensive square root operation
func euclideanSq(a, b []float32) float32 {
	var sum float32
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return sum
}
