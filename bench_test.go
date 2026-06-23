package vectorengine

import "testing"

// helper to build a vector
func makeVec(dim int, v float32) []float32 {
	vec := make([]float32, dim)
	for i := range vec {
		vec[i] = v
	}
	return vec
}

func BenchmarkInsert(b *testing.B) {
	dim := 128
	k := 10

	vec := makeVec(dim, 1.0)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g := NewGraphStore(dim, k)

		// build graph
		for j := 0; j < 5000; j++ {
			_ = g.Insert(vec)
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	dim := 128
	k := 10

	g := NewGraphStore(dim, k)

	// preload graph
	vec := makeVec(dim, 1.0)

	for i := 0; i < 100000; i++ {
		_ = g.Insert(vec)
	}

	query := makeVec(dim, 0.5)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = g.Search(query)
	}
}
