package vectorengine

import (
	"testing"
)

// =========================================================
// TEST GRAPH BUILDER
// =========================================================

func makeTestGraph() *Graph {

	g := NewGraphStore(4, 2, 20)

	// deterministic vectors (distance = simple numeric difference)
	for i := 1; i <= 20; i++ {
		g.SetVector(i, []float32{float32(i)})
	}

	// simple chain + some shortcuts
	for i := 1; i < 20; i++ {
		g.AddNeighbor(i, i+1, 0)
	}

	// add reverse links (more realistic graph)
	for i := 2; i <= 20; i++ {
		g.AddNeighbor(i, i-1, 0)
	}

	g.Size = 20
	g.EntryPoint = 1
	g.EntryPointLevel = 0

	return g
}

// =========================================================
// 1. BASIC EXECUTION TEST
// =========================================================

func TestGenerateCandidatePoolBasic(t *testing.T) {

	g := makeTestGraph()

	query := []float32{10}

	pool, err := g.GenerateCandidatePool(query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pool) == 0 {
		t.Fatalf("expected non-empty pool")
	}

	for _, c := range pool {
		if c.id <= 0 || c.id > g.Size {
			t.Fatalf("invalid candidate id: %d", c.id)
		}
	}
}

// =========================================================
// 2. DETERMINISM TEST
// =========================================================

func TestGenerateCandidatePoolDeterminism(t *testing.T) {

	g := makeTestGraph()

	query := []float32{7}

	pool1, err := g.GenerateCandidatePool(query)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	pool2, err := g.GenerateCandidatePool(query)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(pool1) != len(pool2) {
		t.Fatalf("non-deterministic pool size")
	}

	// compare ids (order doesn't matter, so use membership check)
	set := make(map[int]bool)

	for _, c := range pool1 {
		set[c.id] = true
	}

	for _, c := range pool2 {
		if !set[c.id] {
			t.Fatalf("pool mismatch between runs")
		}
	}
}

// =========================================================
// 3. VALIDATION TESTS
// =========================================================

func TestGenerateCandidatePoolValidation(t *testing.T) {

	g := makeTestGraph()

	_, err := g.GenerateCandidatePool([]float32{1, 2, 3}) // wrong dimension
	if err == nil {
		t.Fatalf("expected dimension mismatch error")
	}

	emptyGraph := NewGraphStore(4, 2, 10)

	_, err = emptyGraph.GenerateCandidatePool([]float32{1, 2, 3, 4})
	if err == nil {
		t.Fatalf("expected empty graph error")
	}
}

// =========================================================
// 4. EF SIZE BOUND TEST
// =========================================================

func TestGenerateCandidatePoolEFBound(t *testing.T) {

	g := makeTestGraph()

	query := []float32{15}

	pool, err := g.GenerateCandidatePool(query)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(pool) > EF {
		t.Fatalf("pool exceeds EF limit: %d", len(pool))
	}
}

// =========================================================
// 5. QUALITY SANITY CHECK (VERY LIGHT)
// =========================================================

func TestGenerateCandidatePoolQuality(t *testing.T) {

	g := makeTestGraph()

	query := []float32{18}

	pool, err := g.GenerateCandidatePool(query)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	// best candidate should be close to query (>= 16-20 region)
	foundGood := false

	for _, c := range pool {
		if c.id >= 15 && c.id <= 20 {
			foundGood = true
			break
		}
	}

	if !foundGood {
		t.Fatalf("pool missing good region candidates")
	}
}

// =========================================================
// 6. STABILITY TEST (NO CRASH UNDER REPEATS)
// =========================================================

func TestGenerateCandidatePoolStability(t *testing.T) {

	g := makeTestGraph()

	query := []float32{12}

	for i := 0; i < 50; i++ {
		_, err := g.GenerateCandidatePool(query)
		if err != nil {
			t.Fatalf("iteration %d failed: %v", i, err)
		}
	}
}
