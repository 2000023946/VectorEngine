package accuracy

import (
	"fmt"
	"math/rand"
	"testing"

	"vectorengine/src"
)

const (
	numQueries = 1_000
	dimension  = 128
	k          = 10
)

func generateVector() []float64 {
	vector := make([]float64, dimension)

	for i := range vector {
		vector[i] = rand.Float64()
	}

	return vector
}

func distance(a []float64, b []float64) float64 {
	var sum float64

	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return sum
}

func bruteForce(
	vectors [][]float64,
	ids []int,
	query []float64,
	k int,
) map[int]bool {

	type candidate struct {
		id       int
		distance float64
	}

	results := make([]candidate, 0, len(vectors))

	for i, vector := range vectors {
		results = append(results, candidate{
			id:       ids[i],
			distance: distance(query, vector),
		})
	}

	// Simple exact sort.
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].distance < results[i].distance {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	groundTruth := make(map[int]bool, k)

	for i := 0; i < k && i < len(results); i++ {
		groundTruth[results[i].id] = true
	}

	return groundTruth
}

func runAccuracyTest(
	t *testing.T,
	numVectors int,
) {
	t.Helper()

	rand.Seed(42)

	engine := src.NewVectorEngine()

	vectors := make([][]float64, numVectors)
	ids := make([]int, numVectors)

	// ----------------------------------------------
	// Build dataset
	// ----------------------------------------------

	for i := 0; i < numVectors; i++ {
		vector := generateVector()

		vectors[i] = vector
		ids[i] = i

		err := engine.Insert(i, vector)

		if err != nil {
			t.Fatal(err)
		}
	}

	// ----------------------------------------------
	// Generate queries
	// ----------------------------------------------

	queries := make([][]float64, numQueries)

	for i := range queries {
		queries[i] = generateVector()
	}

	// ----------------------------------------------
	// Compare IVF against exact brute force
	// ----------------------------------------------

	correct := 0
	total := numQueries * k

	for _, query := range queries {

		groundTruth := bruteForce(
			vectors,
			ids,
			query,
			k,
		)

		results := engine.Search(query, k)

		if len(results) != k {
			t.Fatalf(
				"expected %d results, got %d",
				k,
				len(results),
			)
		}

		for _, result := range results {
			if groundTruth[result.ID] {
				correct++
			}
		}
	}

	accuracy := float64(correct) / float64(total) * 100

	fmt.Printf("\n")
	fmt.Printf("=====================================\n")
	fmt.Printf("VectorEngine Accuracy Test\n")
	fmt.Printf("=====================================\n")
	fmt.Printf("Vectors:   %d\n", numVectors)
	fmt.Printf("Queries:   %d\n", numQueries)
	fmt.Printf("Dimension: %d\n", dimension)
	fmt.Printf("K:         %d\n", k)
	fmt.Printf("Correct:   %d / %d\n", correct, total)
	fmt.Printf("Accuracy:  %.2f%%\n", accuracy)
	fmt.Printf("=====================================\n")

	if accuracy < 90.0 {
		t.Fatalf(
			"accuracy %.2f%% is below required 90%%",
			accuracy,
		)
	}
}

func TestSearchAccuracy10K(t *testing.T) {
	runAccuracyTest(t, 10_000)
}

// func TestSearchAccuracy100K(t *testing.T) {
// 	runAccuracyTest(t, 100_000)
// }

// func TestSearchAccuracy1M(t *testing.T) {
// 	runAccuracyTest(t, 1_000_000)
// }
