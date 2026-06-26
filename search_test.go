package vectorengine

import (
	"testing"
)

// -------------------- HELPERS --------------------

func newSearchGraph() *Graph {
	g := NewGraphStore(3, 2, 5)

	// Ensure graph is "non-empty"
	g.Size = 5

	// valid entry point setup
	g.EntryPoint = 1
	g.EntryPointLevel = 1

	// set vectors so distance logic (if used later) is stable
	g.SetVector(1, []float32{1, 1, 1})
	g.SetVector(2, []float32{2, 2, 2})
	g.SetVector(3, []float32{3, 3, 3})

	// connect simple structure
	g.AddNeighbor(1, 2, 0)
	g.AddNeighbor(2, 3, 0)

	return g
}

// -------------------- EDGE CASE TESTS --------------------

func TestSearchEmptyGraph(t *testing.T) {
	g := NewGraphStore(3, 2, 5)

	_, err := g.Search([]float32{1, 1, 1})
	if err == nil {
		t.Fatalf("expected error for empty graph")
	}
}

func TestSearchQueryMismatch(t *testing.T) {
	g := newSearchGraph()

	_, err := g.Search([]float32{1, 1}) // wrong dim
	if err == nil {
		t.Fatalf("expected query dimension mismatch error")
	}
}

func TestSearchInvalidEntryPointLow(t *testing.T) {
	g := newSearchGraph()
	g.EntryPoint = 0

	_, err := g.Search([]float32{1, 1, 1})
	if err == nil {
		t.Fatalf("expected invalid entry point error")
	}
}

func TestSearchInvalidEntryPointHigh(t *testing.T) {
	g := newSearchGraph()
	g.EntryPoint = 999

	_, err := g.Search([]float32{1, 1, 1})
	if err == nil {
		t.Fatalf("expected invalid entry point error")
	}
}

// -------------------- CORE SEARCH PATH --------------------

func TestSearchBasicExecution(t *testing.T) {
	g := newSearchGraph()

	_, err := g.Search([]float32{1, 1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// -------------------- ENTRY POINT LEVEL CLAMP --------------------

func TestSearchEntryPointLevelClamp(t *testing.T) {
	g := newSearchGraph()

	g.EntryPointLevel = 999 // force invalid high level

	_, err := g.Search([]float32{1, 1, 1})
	if err != nil {
		t.Fatalf("should clamp entry point level instead of failing")
	}
}

// -------------------- FULL FLOW SIMULATION --------------------

// This test ensures the loop executes all layers safely
func TestSearchLayerLoopExecution(t *testing.T) {
	g := newSearchGraph()

	g.EntryPointLevel = 2 // ensure multiple iterations

	_, err := g.Search([]float32{1, 1, 1})
	if err != nil {
		t.Fatalf("unexpected error in layered traversal: %v", err)
	}
}

// -------------------- STRUCTURAL INTEGRITY --------------------

func TestSearchReturnsValidNode(t *testing.T) {
	g := newSearchGraph()

	result, err := g.Search([]float32{1, 1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result <= 0 {
		t.Fatalf("invalid result node: %d", result)
	}
}
