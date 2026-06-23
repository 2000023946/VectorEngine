package vectorengine

import "testing"

func buildSearchGraph() *Graph {
	g := NewGraphStore(2, 2)

	// Build a simple chain:
	// 0 -> 1 -> 2 (2 is closest to origin)
	g.Nodes[0] = &Node{
		ID:        0,
		Vector:    []float32{10, 10},
		Neighbors: []int{1},
	}

	g.Nodes[1] = &Node{
		ID:        1,
		Vector:    []float32{5, 5},
		Neighbors: []int{2},
	}

	g.Nodes[2] = &Node{
		ID:        2,
		Vector:    []float32{1, 1},
		Neighbors: []int{},
	}

	g.LastID = 0

	return g
}

func TestSearchEmptyStore(t *testing.T) {
	g := NewGraphStore(2, 2)

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

	result, err := g.Search([]float32{0, 0})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 2 {
		t.Fatalf("expected node 2, got %d", result.ID)
	}
}

func TestSearchReturnsDistance(t *testing.T) {
	g := buildSearchGraph()

	result, err := g.Search([]float32{0, 0})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Distance <= 0 {
		t.Fatalf("expected positive distance, got %f", result.Distance)
	}
}

func TestSearchExactMatch(t *testing.T) {
	g := NewGraphStore(2, 2)

	g.Nodes[0] = &Node{
		ID:        0,
		Vector:    []float32{0, 0},
		Neighbors: []int{},
	}

	g.LastID = 0

	result, err := g.Search([]float32{0, 0})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 0 {
		t.Fatalf("expected node 0, got %d", result.ID)
	}

	if result.Distance != 0 {
		t.Fatalf("expected distance 0, got %f", result.Distance)
	}
}

func TestSearchMultiHopTraversal(t *testing.T) {
	g := NewGraphStore(2, 2)

	g.Nodes[0] = &Node{ID: 0, Vector: []float32{20, 20}, Neighbors: []int{1}}
	g.Nodes[1] = &Node{ID: 1, Vector: []float32{10, 10}, Neighbors: []int{2}}
	g.Nodes[2] = &Node{ID: 2, Vector: []float32{1, 1}, Neighbors: []int{}}

	g.LastID = 0

	result, err := g.Search([]float32{0, 0})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 2 {
		t.Fatalf("expected node 2, got %d", result.ID)
	}
}

func TestSearchUsesTraverseCorrectly(t *testing.T) {
	g := buildSearchGraph()

	result, err := g.Search([]float32{0, 0})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := g.Nodes[result.ID]
	if !ok {
		t.Fatal("returned node does not exist in graph")
	}
}
