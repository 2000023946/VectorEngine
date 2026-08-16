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

// --------------------------------------------------
// TEST DATA
// --------------------------------------------------
//
// One contiguous input buffer.
//
// This avoids creating one []float64 allocation
// for every vector.
//
// Layout:
//
// [vector 0][vector 1][vector 2]...[vector N]
//
// Each vector occupies 128 float64 values.
// --------------------------------------------------

func generateData(numVectors int) []float64 {
	data := make([]float64, numVectors*dimension)

	for i := range data {
		data[i] = rand.Float64()
	}

	return data
}

func vectorAt(data []float64, index int) []float64 {
	offset := index * dimension

	return data[offset : offset+dimension]
}

// --------------------------------------------------
// INITIALIZATION
// --------------------------------------------------
//
// Measures the cost of creating the engine.
//
// This intentionally measures the ~2.064 GB allocation.
// --------------------------------------------------

func BenchmarkInitialization(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = src.NewVectorEngine()
	}
}

// --------------------------------------------------
// INSERT
// --------------------------------------------------
//
// Important:
//
// Engine allocation happens ONCE.
//
// Input data is generated ONCE.
//
// Each benchmark iteration:
//
//     Reset()
//       ↓
//     Insert N vectors
//
// No new VectorEngine is created for every iteration.
// No input vectors are allocated during the benchmark.
// No vector payload allocation occurs inside Insert.
// --------------------------------------------------

func benchmarkInsert(b *testing.B, numVectors int) {
	b.Helper()

	// Generate the input dataset once.
	data := generateData(numVectors)

	// Allocate the engine once.
	engine := src.NewVectorEngine()

	b.ReportAllocs()
	b.ResetTimer()

	for iteration := 0; iteration < b.N; iteration++ {

		// Reset only logical state.
		//
		// We intentionally DO NOT zero the 2 GB buffer.
		// Insert overwrites every location that it uses.
		engine.Reset()

		for i := 0; i < numVectors; i++ {
			err := engine.Insert(
				i,
				vectorAt(data, i),
			)

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
// SEARCH SETUP
// --------------------------------------------------

func createEngine(
	numVectors int,
	data []float64,
) *src.VectorEngine {

	engine := src.NewVectorEngine()

	for i := 0; i < numVectors; i++ {
		err := engine.Insert(
			i,
			vectorAt(data, i),
		)

		if err != nil {
			panic(err)
		}
	}

	return engine
}

// --------------------------------------------------
// SEARCH
// --------------------------------------------------
//
// Dataset creation happens BEFORE the benchmark.
//
// Only Search() is timed.
// --------------------------------------------------

func benchmarkSearch(
	b *testing.B,
	numVectors int,
) {
	b.Helper()

	// Generate dataset once.
	data := generateData(numVectors)

	// Build database once.
	engine := createEngine(
		numVectors,
		data,
	)

	// Generate query once.
	queryData := generateData(1)
	query := vectorAt(queryData, 0)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = engine.Search(
			query,
			k,
		)
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
