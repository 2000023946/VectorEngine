package performance

import (
	"math/rand"
	"testing"

	"vectorengine/src"
)

const (
	dimension = 128
	k         = 10
)

func generateVector() []float64 {
	vector := make([]float64, dimension)

	for i := range vector {
		vector[i] = rand.Float64()
	}

	return vector
}

func createEngine(numVectors int) *src.VectorEngine {
	engine := src.NewVectorEngine()

	for i := 0; i < numVectors; i++ {
		err := engine.Insert(i, generateVector())

		if err != nil {
			panic(err)
		}
	}

	return engine
}

// --------------------------------------------------
// INITIALIZATION BENCHMARK
// --------------------------------------------------

func BenchmarkInitialization(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		src.NewVectorEngine()
	}
}

// --------------------------------------------------
// INSERT BENCHMARKS
// --------------------------------------------------

func benchmarkInsert(b *testing.B, numVectors int) {
	b.Helper()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		engine := src.NewVectorEngine()

		for j := 0; j < numVectors; j++ {
			err := engine.Insert(j, generateVector())

			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkInsert10K(b *testing.B) {
	benchmarkInsert(b, 10_000)
}

func BenchmarkInsert100K(b *testing.B) {
	benchmarkInsert(b, 100_000)
}

func BenchmarkInsert1M(b *testing.B) {
	benchmarkInsert(b, 1_000_000)
}

// --------------------------------------------------
// SEARCH BENCHMARKS
// --------------------------------------------------

func benchmarkSearch(b *testing.B, numVectors int) {
	b.Helper()

	// Build database before timing.
	engine := createEngine(numVectors)

	query := generateVector()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		engine.Search(query, k)
	}
}

func BenchmarkSearch10K(b *testing.B) {
	benchmarkSearch(b, 10_000)
}

func BenchmarkSearch100K(b *testing.B) {
	benchmarkSearch(b, 100_000)
}

func BenchmarkSearch1M(b *testing.B) {
	benchmarkSearch(b, 1_000_000)
}
