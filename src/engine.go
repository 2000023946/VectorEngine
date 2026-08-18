package src

import (
	"fmt"
	"math"
	"sort"
)

const (
	maxVectors = 2_000_000
	dimension  = 128

	// int8 range used for quantized vectors.
	// We use [-127, 127] instead of [-128, 127] so the
	// quantization is symmetric around zero.
	quantizeScale = 127.0
)

type Result struct {
	ID       int
	Distance float64
}

type VectorEngine struct {
	// --------------------------------------------------
	// Flat contiguous INT8 vector storage
	// --------------------------------------------------
	//
	// 2,000,000 vectors × 128 dimensions × 1 byte
	// = 256 MB
	//
	// Vector i occupies:
	//
	// vectorData[i*dimension : (i+1)*dimension]
	//
	// This is 4x smaller than the previous int32 layout.
	// --------------------------------------------------
	vectorData []int8

	// External ID for each vector.
	ids []int

	dimension int
	count     int
}

func NewVectorEngine() *VectorEngine {
	return &VectorEngine{
		vectorData: make([]int8, maxVectors*dimension),
		ids:        make([]int, maxVectors),
		dimension:  dimension,
	}
}

// quantize converts a single float64 to its int8 quantized form.
//
// Expected input range:
//
//	[-1.0, 1.0]
//
// Mapping:
//
//	-1.0 -> -127
//	 0.0 ->    0
//	+1.0 -> +127
//
// Values outside [-1, 1] are clamped to prevent int8 overflow.
func quantize(v float64) int8 {
	if v > 1.0 {
		v = 1.0
	} else if v < -1.0 {
		v = -1.0
	}

	return int8(math.Round(v * quantizeScale))
}

// quantizeVector quantizes a whole vector.
//
// This allocates only once per query/insert operation and is
// deliberately kept off the hot comparison loop.
func quantizeVector(values []float64) []int8 {
	out := make([]int8, len(values))

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

	// Quantize once at insertion time.
	//
	// The stored representation is only 1 byte per dimension,
	// reducing the vector memory footprint by 4x compared
	// with int32 storage.
	for i, v := range values {
		ve.vectorData[offset+i] = quantize(v)
	}

	ve.ids[index] = id
	ve.count++

	return nil
}

// squaredDistance calls the INT8 NEON assembly routine against
// the vector stored at index, using an already-quantized query.
func (ve *VectorEngine) squaredDistance(
	quantizedQuery []int8,
	index int,
) int64 {
	offset := index * ve.dimension

	stored := ve.vectorData[offset : offset+ve.dimension]

	return SquaredDistanceNEON(quantizedQuery, stored)
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

	// Quantize the query exactly once.
	//
	// We do NOT quantize once per comparison.
	quantizedQuery := quantizeVector(query)

	results := make([]Result, 0, k)

	for i := 0; i < ve.count; i++ {
		distance := ve.squaredDistance(
			quantizedQuery,
			i,
		)

		result := Result{
			ID: ve.ids[i],

			// Distance is returned as float64 to preserve
			// the existing external API.
			Distance: float64(distance),
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

		// Replace it if the new vector is closer.
		if result.Distance < results[worst].Distance {
			results[worst] = result
		}
	}

	// Return nearest neighbors first.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	return results
}

func (ve *VectorEngine) Reset() {
	ve.count = 0
}
