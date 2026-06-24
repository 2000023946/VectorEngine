package vectorengine

import "testing"

func buildSearchGraph() *Graph {
	g := NewGraphStore(2, 2, 10)

	// IMPORTANT: build graph using Insert so Size is valid
	_ = g.Insert([]float32{10, 10}) // node 0
	_ = g.Insert([]float32{5, 5})   // node 1
	_ = g.Insert([]float32{1, 1})   // node 2

	// optional structure (kept for traversal behavior)
	g.AddNeighbor(0, 1)
	g.AddNeighbor(1, 2)

	return g
}

func TestSearchEmptyStore(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	_, err := g.Search([]float32{1, 1})

	if err == nil {
		t.Fatal("expected error for empty store")
	}
}

func TestSearchInvalidDimension(t *testing.T) {
	g := buildSearchGraph()

	_, err := g.Search([]float32{1, 2, 3})

	if err == nil {
		t.Fatal("expected dimension mismatch error")
	}
}

func TestSearchFindsNearestNode(t *testing.T) {
	g := buildSearchGraph()

	res, err := g.Search([]float32{0, 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ID != 2 {
		t.Fatalf("expected node 2, got %d", res.ID)
	}

	if res.Distance <= 0 {
		t.Fatalf("expected positive distance, got %f", res.Distance)
	}
}

func TestSearchReturnsDistance(t *testing.T) {
	g := buildSearchGraph()

	res, err := g.Search([]float32{0, 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Distance <= 0 {
		t.Fatalf("expected positive distance, got %f", res.Distance)
	}
}

func TestSearchExactMatch(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	_ = g.Insert([]float32{0, 0})

	// optional edge (not required for correctness)
	g.AddNeighbor(0, 0)

	res, err := g.Search([]float32{0, 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ID != 0 {
		t.Fatalf("expected node 0, got %d", res.ID)
	}

	if res.Distance != 0 {
		t.Fatalf("expected distance 0, got %f", res.Distance)
	}
}

func TestSearchMultiHopTraversal(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	_ = g.Insert([]float32{20, 20}) // 0
	_ = g.Insert([]float32{10, 10}) // 1
	_ = g.Insert([]float32{1, 1})   // 2

	g.AddNeighbor(0, 1)
	g.AddNeighbor(1, 2)

	res, err := g.Search([]float32{0, 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ID != 2 {
		t.Fatalf("expected node 2, got %d", res.ID)
	}
}

func TestSearchUsesTraverseCorrectly(t *testing.T) {
	g := buildSearchGraph()

	_, err := g.Search([]float32{0, 0})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
