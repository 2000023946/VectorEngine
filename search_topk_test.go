package vectorengine

import (
	"testing"
)

// =========================================================
// TEST FIXTURE GRAPH
// =========================================================

func newTopKTestGraph() *Graph {

	g := NewGraphStore(4, 3, 30)

	// deterministic 4D vectors
	for i := 1; i <= 30; i++ {
		g.SetVector(i, makeVec(g.Dimension, float32(i)))
	}

	g.Size = 30
	g.EntryPoint = 1
	g.EntryPointLevel = 0

	// chain connectivity
	for i := 1; i < 30; i++ {
		g.AddNeighbor(i, i+1, 0)
	}

	return g
}

// =========================================================
// TEST 1: BASIC TOP-K EXECUTION
// =========================================================

func TestUnitSearchTopKBasic(t *testing.T) {

	g := newTopKTestGraph()

	res, err := g.SearchTopK(makeVec(g.Dimension, 10), 5)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(res) != 5 {
		t.Fatalf("expected 5 results got %d", len(res))
	}

	for i := 1; i < len(res); i++ {
		if res[i].dist < res[i-1].dist {
			t.Fatalf("results not sorted by distance")
		}
	}
}

// =========================================================
// TEST 2: K CLAMPING
// =========================================================

func TestUnitSearchTopKClamp(t *testing.T) {

	g := newTopKTestGraph()

	res, err := g.SearchTopK(makeVec(g.Dimension, 10), 1000)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	// clamp should never exceed dataset size
	if len(res) > g.Size {
		t.Fatalf("expected clamped size <= %d got %d", g.Size, len(res))
	}

	// if k is huge, we should return all available results
	if len(res) == 0 {
		t.Fatalf("expected non-empty result")
	}
}

// =========================================================
// TEST 3: INVALID K
// =========================================================

func TestUnitSearchTopKInvalidK(t *testing.T) {

	g := newTopKTestGraph()

	_, err := g.SearchTopK(makeVec(g.Dimension, 10), 0)
	if err == nil {
		t.Fatalf("expected error for k=0")
	}
}

// =========================================================
// TEST 4: EMPTY GRAPH
// =========================================================

func TestUnitSearchTopKEmptyGraph(t *testing.T) {

	g := NewGraphStore(4, 3, 10)

	_, err := g.SearchTopK(makeVec(g.Dimension, 1), 5)
	if err == nil {
		t.Fatalf("expected error for empty graph")
	}
}

// =========================================================
// TEST 5: DIMENSION MISMATCH
// =========================================================

func TestUnitSearchTopKDimensionMismatch(t *testing.T) {

	g := newTopKTestGraph()

	_, err := g.SearchTopK([]float32{1, 2, 3}, 5)
	if err == nil {
		t.Fatalf("expected dimension mismatch error")
	}
}

// =========================================================
// TEST 6: DETERMINISM
// =========================================================

func TestUnitSearchTopKDeterminism(t *testing.T) {

	g := newTopKTestGraph()

	r1, err := g.SearchTopK(makeVec(g.Dimension, 15), 10)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	r2, err := g.SearchTopK(makeVec(g.Dimension, 15), 10)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(r1) != len(r2) {
		t.Fatalf("non deterministic result size")
	}

	for i := range r1 {
		if r1[i].id != r2[i].id {
			t.Fatalf("non deterministic results at index %d", i)
		}
	}
}
