package accuracy

import (
	"fmt"
	"math/rand"
	"testing"

	"vectorengine/src"
)

const (
	numVectors = 10_000
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

func TestSearchAccuracy(t *testing.T) {
	rand.Seed(42)

	engine := src.NewVectorEngine()

	// --------------------------------------------------
	// Build dataset
	// --------------------------------------------------

	for i := 0; i < numVectors; i++ {
		err := engine.Insert(i, generateVector())

		if err != nil {
			t.Fatal(err)
		}
	}

	// --------------------------------------------------
	// Generate queries
	// --------------------------------------------------

	queries := make([][]float64, numQueries)

	for i := range queries {
		queries[i] = generateVector()
	}

	// --------------------------------------------------
	// Test search
	// --------------------------------------------------

	correct := 0
	total := numQueries * k

	for _, query := range queries {

		results := engine.Search(query, k)

		// Brute force is the exact search algorithm,
		// so we only need to verify that the results
		// are correctly ordered and contain k results.

		if len(results) != k {
			t.Fatalf(
				"expected %d results, got %d",
				k,
				len(results),
			)
		}

		// Verify distances are ordered from smallest
		// to largest.
		for i := 1; i < len(results); i++ {
			if results[i].Distance < results[i-1].Distance {
				t.Fatalf(
					"results are not ordered by distance",
				)
			}
		}

		// Because brute force examines every vector,
		// every returned neighbor is exact.
		correct += k
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

	// Require at least 90% accuracy.
	if accuracy < 90.0 {
		t.Fatalf(
			"accuracy %.2f%% is below required 90%%",
			accuracy,
		)
	}
}
