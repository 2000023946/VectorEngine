package vectorengine

import "testing"

// helper
func makeVec(dim int, v float32) []float32 {
	vec := make([]float32, dim)
	for i := range vec {
		vec[i] = v
	}
	return vec
}

//
// =========================
// 1. PURE INSERT BENCHMARK
// =========================
// Measures ONLY Insert cost
//

func BenchmarkInsert(b *testing.B) {
	dim := 128
	k := 10
	maxNodes := 100000

	vec := makeVec(dim, 1.0)

	g := NewGraphStore(dim, k, maxNodes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = g.Insert(vec)
	}
}

//
// =========================
// 2. PURE SEARCH BENCHMARK
// =========================
// Measures ONLY Search cost (prebuilt graph)
//

func BenchmarkSearch(b *testing.B) {
	dim := 128
	k := 10
	maxNodes := 100000

	g := NewGraphStore(dim, k, maxNodes)

	vec := makeVec(dim, 1.0)

	// build graph ONCE (this is allowed because Search benchmark must have data)
	for i := 0; i < maxNodes; i++ {
		_ = g.Insert(vec)
	}

	query := makeVec(dim, 0.5)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = g.Search(query)
	}
}

//
// =========================
// 3. PURE BUILD GRAPH BENCHMARK
// =========================
// Measures ONLY struct + memory allocation cost
//

func BenchmarkBuildGraph(b *testing.B) {
	dim := 128
	k := 10
	maxNodes := 100000

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewGraphStore(dim, k, maxNodes)
	}
}
