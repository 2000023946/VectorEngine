package db

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// ==========================================
// ACCURACY TESTS
// ==========================================

func TestVectorEngine_InsertAndSearch(t *testing.T) {
	engine := NewVectorEngine(3, 2)

	pointA := []float32{1.0, 1.0, 1.0}
	pointB := []float32{5.0, 5.0, 5.0}
	pointC := []float32{10.0, 10.0, 10.0}

	engine.Insert(pointA)
	engine.Insert(pointB)
	engine.Insert(pointC)

	for {
		if engine.state.Load().Phase == PhaseIndexed {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	query := []float32{4.5, 4.5, 4.5}
	result := engine.Search(query)

	if !reflect.DeepEqual(result, pointB) {
		t.Errorf("Expected closest match to be %v, got %v", pointB, result)
	}
}

// ==========================================
// BENCHMARK TESTS
// ==========================================

// generateDummyEngine fills the DB with random vectors
func generateDummyEngine(size int, dim int) *VectorEngine {
	// Optimal K is the square root of N
	k := int(math.Sqrt(float64(size)))
	if k < 2 {
		k = 2
	}

	engine := NewVectorEngine(size, k)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < size; i++ {
		data := make([]float32, dim)
		for j := 0; j < dim; j++ {
			data[j] = rng.Float32() * 100.0
		}
		engine.Insert(data)
	}

	for {
		if engine.state.Load().Phase == PhaseIndexed {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	return engine
}

func BenchmarkEngine_Search(b *testing.B) {
	dim := 30
	sizes := []int{1000, 10000, 50000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Dataset_%d", size), func(b *testing.B) {
			engine := generateDummyEngine(size, dim)
			query := make([]float32, dim)
			for j := 0; j < dim; j++ {
				query[j] = 50.0
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				engine.Search(query)
			}
		})
	}
}

func BenchmarkSearchParallel(b *testing.B) {
	// N = 50,000, K = sqrt(50,000) = 223
	numVectors := 50000
	k := 223
	vectorsPerBucket := numVectors / k
	dim := 30 // Aligned with the sequential test

	engine := NewVectorEngine(numVectors, k)

	centroids := make([][]float32, k)
	buckets := make([][][]float32, k)

	for i := 0; i < k; i++ {
		centroids[i] = make([]float32, dim)
		for j := 0; j < vectorsPerBucket; j++ {
			buckets[i] = append(buckets[i], make([]float32, dim))
		}
	}

	newState := &IndexState{
		Phase:     PhaseIndexed,
		Centroids: centroids,
		Buckets:   buckets,
	}
	engine.state.Store(newState)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		query := make([]float32, dim)
		for pb.Next() {
			engine.Search(query)
		}
	})
}
