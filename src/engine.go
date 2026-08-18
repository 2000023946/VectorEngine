package src

import (
	"fmt"
	"math"
	"sort"
)

const (
	maxVectors = 2_000_000
	dimension  = 128

	// quantizeScale controls precision vs. overflow headroom.
	// float value v is stored as int32(round(v * quantizeScale)).
	quantizeScale = 10000.0
)

type Result struct {
	ID       int
	Distance float64
}

type VectorEngine struct {
	// --------------------------------------------------
	// Flat contiguous QUANTIZED vector storage
	// --------------------------------------------------
	//
	// 2,000,000 vectors × 128 dimensions × 4 bytes (int32)
	// = 1.024 GB  (half the size of the float64 version)
	//
	// Vector i occupies:
	//
	// vectorData[i*dimension : (i+1)*dimension]
	// --------------------------------------------------
	vectorData []int32

	// External ID for each vector.
	ids []int

	dimension int
	count     int
}

func NewVectorEngine() *VectorEngine {
	return &VectorEngine{
		vectorData: make([]int32, maxVectors*dimension),
		ids:        make([]int, maxVectors),
		dimension:  dimension,
	}
}

// quantize converts a single float64 to its int32 quantized form.
func quantize(v float64) int32 {
	return int32(math.Round(v * quantizeScale))
}

// quantizeVector quantizes a whole vector, allocating a new []int32.
func quantizeVector(values []float64) []int32 {
	out := make([]int32, len(values))
	for i, v := range values {
		out[i] = quantize(v)
	}
	return out
}

func (ve *VectorEngine) Insert(id int, values []float64) error {
	if len(values) != ve.dimension {
		return fmt.Errorf(
			"dimension mismatch: expected %d, got %d",
			ve.dimension,
			len(values),
		)
	}

	if ve.count >= maxVectors {
		return fmt.Errorf(
			"vector capacity exceeded: maximum %d vectors",
			maxVectors,
		)
	}

	index := ve.count
	offset := index * ve.dimension

	// Quantize once, here, at insert time — this cost is O(1)
	// per vector and never repeated on the (much hotter) search path.
	for i, v := range values {
		ve.vectorData[offset+i] = quantize(v)
	}

	ve.ids[index] = id
	ve.count++

	return nil
}

// squaredDistance calls the NEON assembly routine against the vector
// stored at index, using an already-quantized query.
// func (ve *VectorEngine) squaredDistance(
// 	quantizedQuery []int32,
// 	index int,
// ) int64 {
// 	offset := index * ve.dimension
// 	stored := ve.vectorData[offset : offset+ve.dimension]

// 	return SquaredDistanceNEON(quantizedQuery, stored)
// }

func (ve *VectorEngine) squaredDistance(
	quantizedQuery []int32,
	index int,
) int64 {
	offset := index * ve.dimension
	stored := ve.vectorData[offset : offset+ve.dimension]

	var distance int64

	for i := 0; i < ve.dimension; i++ {
		diff := int64(quantizedQuery[i]) - int64(stored[i])
		distance += diff * diff
	}

	return distance
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

	// Quantize the query exactly once per search — not once per comparison.
	quantizedQuery := quantizeVector(query)

	results := make([]Result, 0, k)

	for i := 0; i < ve.count; i++ {
		distance := ve.squaredDistance(quantizedQuery, i)

		result := Result{
			ID: ve.ids[i],
			// Convert back to float64 for the external API. Relative
			// ordering is preserved even though absolute precision
			// was reduced by quantization.
			Distance: float64(distance),
		}

		if len(results) < k {
			results = append(results, result)
			continue
		}

		worst := 0
		for j := 1; j < k; j++ {
			if results[j].Distance > results[worst].Distance {
				worst = j
			}
		}

		if result.Distance < results[worst].Distance {
			results[worst] = result
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	return results
}

func (ve *VectorEngine) Reset() {
	ve.count = 0
}
