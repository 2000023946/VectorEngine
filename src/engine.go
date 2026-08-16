package src

import (
	"fmt"
	"math"
	"sort"
)

type Vector struct {
	ID     int
	Values []float64
}

type Result struct {
	ID       int
	Distance float64
}

type VectorEngine struct {
	vectors   []Vector
	dimension int
}

func NewVectorEngine() *VectorEngine {
	return &VectorEngine{
		vectors: make([]Vector, 0),
	}
}

func (ve *VectorEngine) Insert(id int, values []float64) error {
	// The first vector establishes the dimension
	// of the vector space.
	if ve.dimension == 0 {
		ve.dimension = len(values)
	} else if len(values) != ve.dimension {
		return fmt.Errorf(
			"dimension mismatch: expected %d, got %d",
			ve.dimension,
			len(values),
		)
	}

	vector := Vector{
		ID:     id,
		Values: values,
	}

	ve.vectors = append(ve.vectors, vector)

	return nil
}

func distance(a []float64, b []float64) float64 {
	var sum float64

	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

func (ve *VectorEngine) Search(query []float64, k int) []Result {
	if len(query) != ve.dimension {
		return nil
	}

	results := make([]Result, 0, len(ve.vectors))

	// Brute force:
	// compare the query against EVERY vector.
	for _, vector := range ve.vectors {
		d := distance(query, vector.Values)

		results = append(results, Result{
			ID:       vector.ID,
			Distance: d,
		})
	}

	// Closest vectors come first.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	if k > len(results) {
		k = len(results)
	}

	if k < 0 {
		k = 0
	}

	return results[:k]
}
