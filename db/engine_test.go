package db

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// ==========================================
// ACCURACY TESTS
// ==========================================

func TestEngine_InsertAndSearch(t *testing.T) {
	// 1. Initialize an engine for 3D vectors (e.g., simple X,Y,Z for testing)
	engine := NewEngine(10, 3)

	// 2. Insert test data
	testData := []Vector{
		{ID: "point_A", Data: []float32{1.0, 1.0, 1.0}},
		{ID: "point_B", Data: []float32{5.0, 5.0, 5.0}},
		{ID: "point_C", Data: []float32{10.0, 10.0, 10.0}},
	}

	for _, v := range testData {
		err := engine.Insert(v)
		if err != nil {
			t.Fatalf("Failed to insert vector %s: %v", v.ID, err)
		}
	}

	// Verify insertion count
	if engine.Count() != 3 {
		t.Errorf("Expected 3 vectors in engine, got %d", engine.Count())
	}

	// 3. Perform a Search query
	// This point is closest to point_B (5,5,5)
	query := []float32{4.5, 4.5, 4.5}

	// Ask for the top 2 closest matches
	results, err := engine.Search(query, 2)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// 4. Validate Accuracy
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if results[0].Vector.ID != "point_B" {
		t.Errorf("Expected closest match to be point_B, got %s", results[0].Vector.ID)
	}

	if results[1].Vector.ID != "point_A" {
		t.Errorf("Expected second match to be point_A, got %s", results[1].Vector.ID)
	}
}

// ==========================================
// BENCHMARK TESTS
// ==========================================

// generateDummyEngine fills the DB with random vectors to simulate hardware data
func generateDummyEngine(size int, dim int) *VectorEngine {
	engine := NewEngine(size, dim)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < size; i++ {
		data := make([]float32, dim)
		for j := 0; j < dim; j++ {
			// Random float between 0 and 100
			data[j] = rng.Float32() * 100.0
		}

		engine.Insert(Vector{
			ID:   fmt.Sprintf("vec_%d", i),
			Data: data,
		})
	}
	return engine
}

// BenchmarkEngine_Search_BruteForce measures the latency of your O(N) scan
func BenchmarkEngine_Search_BruteForce(b *testing.B) {
	// Setup: 30-dimensional vectors (e.g., 10 historical readings of X,Y,Z)
	dim := 30

	// Test against different dataset sizes to prove O(N) scaling
	sizes := []int{1000, 10000, 50000}

	for _, size := range sizes {
		// b.Run creates sub-benchmarks for each dataset size
		b.Run(fmt.Sprintf("Dataset_%d", size), func(b *testing.B) {

			// 1. Pre-load the database before starting the timer
			engine := generateDummyEngine(size, dim)

			// 2. Create a stable query vector
			query := make([]float32, dim)
			for j := 0; j < dim; j++ {
				query[j] = 50.0
			}

			// 3. Reset the timer so ingestion time isn't counted in the benchmark
			b.ResetTimer()

			// 4. The actual benchmark loop [1]
			for i := 0; i < b.N; i++ {
				// We search for the Top 10 closest matches
				_, err := engine.Search(query, 10)
				if err != nil {
					b.Fatalf("Search failed during benchmark: %v", err)
				}
			}
		})
	}
}
