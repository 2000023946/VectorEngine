package vectorengine

import "testing"

// -------------------- CONFIG --------------------

const (
	BenchDim      = 64
	BenchK        = 16
	BenchCapacity = 1_000_000
	BenchPreload  = 500_000
)

// -------------------- FIXTURE --------------------

func newBenchGraph1M() *Graph {
	return NewGraphStore(BenchDim, BenchK, BenchCapacity)
}

// -------------------- VECTOR --------------------

func makeBenchVec() []float32 {
	vec := make([]float32, BenchDim)
	for i := range vec {
		vec[i] = float32(i)
	}
	return vec
}

// -------------------- INIT --------------------

func BenchmarkGraphInit1M(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewGraphStore(BenchDim, BenchK, BenchCapacity)
	}
}

// -------------------- INSERT (1M SCALE) --------------------

func BenchmarkInsert1M(b *testing.B) {
	g := newBenchGraph1M()

	vec := makeBenchVec()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := g.Insert(vec)
		if err != nil {
			b.Fatalf("insert failed: %v", err)
		}
	}
}

// -------------------- BULK INSERT (1M SCALE) --------------------

func BenchmarkBulkInsert1M(b *testing.B) {
	g := newBenchGraph1M()

	vec := makeBenchVec()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := g.Insert(vec)
		if err != nil {
			b.Fatalf("bulk insert failed: %v", err)
		}
	}
}

// -------------------- SEARCH (500K PRELOAD) --------------------

func BenchmarkSearch1M(b *testing.B) {
	g := newBenchGraph1M()

	vec := makeBenchVec()

	// -------------------- PRELOAD PHASE --------------------
	for i := 0; i < BenchPreload; i++ {
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

// -------------------- SEARCH AFTER FULL SCALE LOAD --------------------

func BenchmarkSearchFull1M(b *testing.B) {
	g := newBenchGraph1M()

	vec := makeBenchVec()

	// -------------------- FULL LOAD PHASE --------------------
	for i := 0; i < BenchCapacity; i++ {
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
