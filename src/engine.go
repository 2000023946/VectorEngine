package src

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
)

const (
	maxVectors = 2_000_000
	dimension  = 128

	// Number of independent search partitions.
	// Each partition is processed by one goroutine.
	searchWorkers = 6

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
	// = 1.024 GB
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
	mu        sync.Mutex
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
	return int32(v * quantizeScale)
}

// quantizeVector quantizes a whole vector.
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

	// Quantize once at insert time.
	// This cost is paid once and never repeated during search.
	for i, v := range values {
		ve.vectorData[offset+i] = quantize(v)
	}

	ve.ids[index] = id
	ve.count++

	return nil
}

// squaredDistance computes the distance between the quantized query
// and the vector stored at index using the NEON SIMD kernel.
func (ve *VectorEngine) squaredDistance(
	quantizedQuery []int32,
	index int,
) int64 {
	offset := index * ve.dimension

	stored := ve.vectorData[offset : offset+ve.dimension]

	return SquaredDistanceNEON(quantizedQuery, stored)
}

// searchRange searches one independent partition of the dataset.
//
// Each goroutine gets its own range and its own local top-k result.
// No shared result state is modified during the search.
func (ve *VectorEngine) searchRange(
	quantizedQuery []int32,
	start int,
	end int,
	k int,
) []Result {

	results := make([]Result, 0, k)

	for i := start; i < end; i++ {
		distance := ve.squaredDistance(quantizedQuery, i)

		result := Result{
			ID:       ve.ids[i],
			Distance: float64(distance),
		}

		if len(results) < k {
			results = append(results, result)
			continue
		}

		// Find the worst result in this partition.
		worst := 0

		for j := 1; j < k; j++ {
			if results[j].Distance > results[worst].Distance {
				worst = j
			}
		}

		// Replace the local worst result if this vector is closer.
		if result.Distance < results[worst].Distance {
			results[worst] = result
		}
	}

	return results
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
	// Every goroutine receives the same read-only quantized query.
	quantizedQuery := quantizeVector(query)

	// --------------------------------------------------
	// Determine number of partitions.
	// --------------------------------------------------
	//
	// Never create more workers than vectors.
	//
	// runtime.GOMAXPROCS(0) gives the number of logical
	// processors available to the Go runtime.
	// --------------------------------------------------

	workers := searchWorkers

	if workers > runtime.GOMAXPROCS(0) {
		workers = runtime.GOMAXPROCS(0)
	}

	if workers > ve.count {
		workers = ve.count
	}

	// --------------------------------------------------
	// Split the dataset into M independent ranges.
	// --------------------------------------------------

	resultsPerWorker := make([][]Result, workers)

	chunkSize := (ve.count + workers - 1) / workers

	done := make(chan int, workers)

	for worker := 0; worker < workers; worker++ {
		start := worker * chunkSize
		end := start + chunkSize

		if end > ve.count {
			end = ve.count
		}

		if start >= ve.count {
			done <- worker
			continue
		}

		go func(worker, start, end int) {
			resultsPerWorker[worker] = ve.searchRange(
				quantizedQuery,
				start,
				end,
				k,
			)

			done <- worker
		}(worker, start, end)
	}

	// Wait for every partition to finish.
	for i := 0; i < workers; i++ {
		<-done
	}

	// --------------------------------------------------
	// Merge local top-k results.
	//
	// Each worker produced at most K results.
	//
	// M workers × K results
	//             ↓
	//        final global K
	// --------------------------------------------------

	results := make([]Result, 0, k)

	for worker := 0; worker < workers; worker++ {
		for _, result := range resultsPerWorker[worker] {

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
	}

	// Return results sorted from nearest to farthest.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	return results
}

func (ve *VectorEngine) Reset() {
	ve.count = 0
}
