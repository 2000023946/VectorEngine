package vectorengine

import "testing"

// -------------------- LARGE FIXTURE --------------------

func newBenchGraph(b *testing.B) *Graph {
	// large capacity so no resizing / pressure
	g := NewGraphStore(64, 16, 100000)
	return g
}

// -------------------- INIT BENCHMARK --------------------

func BenchmarkGraphInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewGraphStore(64, 16, 100000)
	}
}

// -------------------- INSERT BENCHMARK --------------------

func BenchmarkInsert(b *testing.B) {
	g := newBenchGraph(b)

	vec := make([]float32, 64)
	for i := range vec {
		vec[i] = float32(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := g.Insert(vec)
		if err != nil {
			b.Fatalf("insert failed: %v", err)
		}
	}
}

// -------------------- BATCH INSERT BENCHMARK --------------------

func BenchmarkBulkInsert(b *testing.B) {
	g := newBenchGraph(b)

	vec := make([]float32, 64)
	for i := range vec {
		vec[i] = float32(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id := i % 100000
		_ = id // ensure distribution idea

		_, err := g.Insert(vec)
		if err != nil {
			b.Fatalf("bulk insert failed: %v", err)
		}
	}
}

// -------------------- SEARCH BENCHMARK --------------------

func BenchmarkSearch(b *testing.B) {
	g := newBenchGraph(b)

	vec := make([]float32, 64)
	for i := range vec {
		vec[i] = float32(i)
	}

	// preload graph so search is meaningful
	for i := 1; i <= 50000; i++ {
		_, err := g.Insert(vec)
		if err != nil {
			b.Fatalf("preload insert failed: %v", err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := g.Search(vec)
		if err != nil {
			b.Fatalf("search failed: %v", err)
		}
	}
}

// -------------------- SEARCH AFTER HEAVY INSERT LOAD --------------------

func BenchmarkSearchAfterHeavyInsert(b *testing.B) {
	g := newBenchGraph(b)

	vec := make([]float32, 64)
	for i := range vec {
		vec[i] = float32(i)
	}

	// heavy insert phase
	for i := 1; i <= 80000; i++ {
		_, err := g.Insert(vec)
		if err != nil {
			b.Fatalf("preload insert failed: %v", err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := g.Search(vec)
		if err != nil {
			b.Fatalf("search failed: %v", err)
		}
	}
}
