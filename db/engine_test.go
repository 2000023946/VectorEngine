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
	// 1. Initialize our engine with a threshold of 3, and a queue size of 10
	engine := NewVectorEngine(3, 10)

	// 2. Insert test data
	pointA := []float32{1.0, 1.0, 1.0}
	pointB := []float32{5.0, 5.0, 5.0}
	pointC := []float32{10.0, 10.0, 10.0}

	engine.Insert(pointA)
	engine.Insert(pointB)
	engine.Insert(pointC)

	// WAIT FOR GOROUTINE: Give the background worker time to pull
	// from the queue and run buildIndex()
	time.Sleep(100 * time.Millisecond)

	// Verify the state machine actually transitioned
	engine.mu.RLock()
	phase := engine.Phase
	engine.mu.RUnlock()

	if phase != PhaseIndexed {
		t.Fatalf("Engine failed to transition to PhaseIndexed. Current phase: %v", phase)
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
	// Provide a queue large enough to hold the entire batch without blocking
	engine := NewVectorEngine(size, size+100)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < size; i++ {
		data := make([]float32, dim)
		for j := 0; j < dim; j++ {
			data[j] = rng.Float32() * 100.0
		}
		engine.Insert(data)
	}

	// WAIT FOR GOROUTINE: Give the background worker a moment to execute
	// Lloyd's algorithm on the full dataset before the benchmark starts
	time.Sleep(500 * time.Millisecond)

	return engine
}

// BenchmarkEngine_Search measures the latency of your cache-friendly bucket scan
func BenchmarkEngine_Search(b *testing.B) {
	// Setup: 30-dimensional vectors
	dim := 30

	// Test against different dataset sizes
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
				engine.Search(query)
			}
		})
	}
}
