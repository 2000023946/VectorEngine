package db

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// ==========================================
// ACCURACY TESTS
// ==========================================

func TestVectorEngine_InsertAndSearch(t *testing.T) {
	// 1. Initialize our engine with a threshold of 3
	// This ensures the index builds on the 3rd insert!
	engine := NewVectorEngine(3)

	// 2. Insert test data
	pointA := []float32{1.0, 1.0, 1.0}
	pointB := []float32{5.0, 5.0, 5.0}
	pointC := []float32{10.0, 10.0, 10.0}

	engine.Insert(pointA)
	engine.Insert(pointB)
	engine.Insert(pointC)

	// Verify the state machine actually transitioned
	if engine.Phase != PhaseIndexed {
		t.Fatalf("Engine failed to transition to PhaseIndexed. Current phase: %v", engine.Phase)
	}

	// 3. Perform a Search query
	// This point is closest to pointB (5,5,5)
	query := []float32{4.5, 4.5, 4.5}

	// Ask for the closest match
	result, dist := engine.Search(query)

	// 4. Validate Accuracy
	if !reflect.DeepEqual(result, pointB) {
		t.Errorf("Expected closest match to be %v, got %v (Distance: %f)", pointB, result, dist)
	}
}

// ==========================================
// BENCHMARK TESTS
// ==========================================

// generateDummyEngine fills the DB with random vectors
func generateDummyEngine(size int, dim int) *VectorEngine {
	// By setting threshold exactly equal to size, the very last
	// insert in this loop will trigger buildIndex() automatically!
	engine := NewVectorEngine(size)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < size; i++ {
		data := make([]float32, dim)
		for j := 0; j < dim; j++ {
			// Random float between 0 and 100
			data[j] = rng.Float32() * 100.0
		}
		engine.Insert(data)
	}

	return engine
}

// BenchmarkEngine_Search measures the latency of your cache-friendly bucket scan
func BenchmarkEngine_Search(b *testing.B) {
	// Setup: 30-dimensional vectors (e.g., 10 historical readings of X,Y,Z)
	dim := 30

	// Test against different dataset sizes to prove O(1) routing + bucket scanning
	sizes := []int{1000, 10000, 50000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Dataset_%d", size), func(b *testing.B) {

			// 1. Pre-load the database before starting the timer
			engine := generateDummyEngine(size, dim)

			// 2. Create a stable query vector
			query := make([]float32, dim)
			for j := 0; j < dim; j++ {
				query[j] = 50.0
			}

			// 3. Reset the timer so ingestion/indexing time isn't counted
			b.ResetTimer()

			// 4. The actual benchmark loop
			for i := 0; i < b.N; i++ {
				// We search for the closest match
				engine.Search(query)
			}
		})
	}
}
