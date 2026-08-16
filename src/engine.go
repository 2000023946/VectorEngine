package src

import (
	"fmt"
	"sort"
)

const (
	maxVectors = 2_000_000
	dimension  = 128
)

type Result struct {
	ID       int
	Distance float64
}

type VectorEngine struct {
	// --------------------------------------------------
	// Flat contiguous vector storage
	// --------------------------------------------------
	//
	// 2,000,000 vectors × 128 dimensions × 8 bytes
	// = 2.048 GB
	//
	// Vector i occupies:
	//
	// vectorData[i*dimension : (i+1)*dimension]
	// --------------------------------------------------
	vectorData []float64

	// External ID for each vector.
	// ids[i] corresponds to vectorData[i*dimension:...].
	ids []int

	dimension int
	count     int
}

func NewVectorEngine() *VectorEngine {
	return &VectorEngine{
		// One large ~2 GB contiguous allocation.
		vectorData: make([]float64, maxVectors*dimension),

		// Preallocate IDs for all possible vectors.
		ids: make([]int, maxVectors),

		dimension: dimension,
	}
}

func (ve *VectorEngine) Insert(id int, values []float64) error {
	// --------------------------------------------------
	// Validate dimension.
	// --------------------------------------------------

	if len(values) != ve.dimension {
		return fmt.Errorf(
			"dimension mismatch: expected %d, got %d",
			ve.dimension,
			len(values),
		)
	}

	// --------------------------------------------------
	// Check capacity.
	// --------------------------------------------------

	if ve.count >= maxVectors {
		return fmt.Errorf(
			"vector capacity exceeded: maximum %d vectors",
			maxVectors,
		)
	}

	// --------------------------------------------------
	// Find the next vector position.
	// --------------------------------------------------
	//
	// Vector 0 → offset 0
	// Vector 1 → offset 128
	// Vector 2 → offset 256
	// ...
	// --------------------------------------------------

	index := ve.count
	offset := index * ve.dimension

	// --------------------------------------------------
	// Copy values into the already allocated memory.
	//
	// This does NOT allocate a new []float64.
	// --------------------------------------------------

	copy(
		ve.vectorData[offset:offset+ve.dimension],
		values,
	)

	// Store the external ID.
	ve.ids[index] = id

	ve.count++

	return nil
}

// distanceToIndex calculates the Euclidean distance between
// the query and the vector stored at index.
//
// The vector is accessed directly from the flat array.
func (ve *VectorEngine) squaredDistance(
	query []float64,
	index int,
) float64 {
	offset := index * ve.dimension

	var sum float64

	for i := 0; i < ve.dimension; i++ {
		diff := query[i] - ve.vectorData[offset+i]
		sum += diff * diff
	}

	return sum
}

func (ve *VectorEngine) Search(query []float64, k int) []Result {
	if len(query) != ve.dimension {
		return nil
	}

	if ve.count == 0 || k <= 0 {
		return []Result{}
	}

	if k > ve.count {
		k = ve.count
	}

	// Allocate only enough space for the requested top-k.
	results := make([]Result, 0, k)

	for i := 0; i < ve.count; i++ {
		distance := ve.squaredDistance(query, i)

		result := Result{
			ID:       ve.ids[i],
			Distance: distance,
		}

		if len(results) < k {
			results = append(results, result)
			continue
		}

		// Find the current worst result.
		worst := 0

		for j := 1; j < k; j++ {
			if results[j].Distance > results[worst].Distance {
				worst = j
			}
		}

		// Replace it if this vector is better.
		if distance < results[worst].Distance {
			results[worst] = result
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	return results
}

// Reset clears the logical contents of the engine without
// reallocating its backing memory.
//
// The vectorData and ids buffers are intentionally reused.
// Insert will overwrite the old values.
func (ve *VectorEngine) Reset() {
	ve.count = 0
}
