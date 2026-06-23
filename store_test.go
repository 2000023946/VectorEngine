package vectorengine

import "testing"

func TestNewGraphStore(t *testing.T) {
	dim := 128
	k := 10

	g := NewGraphStore(dim, k)

	if g == nil {
		t.Fatal("expected graph, got nil")
	}

	if g.Dimension != dim {
		t.Errorf("expected dimension %d, got %d", dim, g.Dimension)
	}

	if g.K != k {
		t.Errorf("expected K %d, got %d", k, g.K)
	}

	if g.LastID != -1 {
		t.Errorf("expected LastID -1, got %d", g.LastID)
	}

	if g.Nodes == nil {
		t.Fatal("expected nodes map to be initialized")
	}

	if len(g.Nodes) != 0 {
		t.Errorf("expected empty nodes map, got %d", len(g.Nodes))
	}
}

func TestTraverseEmptyGraph(t *testing.T) {
	g := NewGraphStore(2, 2)

	query := []float32{0, 0}

	id, visited, err := g.Traverse(query)

	if err == nil {
		t.Fatal("expected error for empty graph")
	}

	if id != -1 {
		t.Fatalf("expected id -1, got %d", id)
	}

	if visited != nil {
		t.Fatal("expected nil visited map")
	}
}

func TestTraverseSingleHop(t *testing.T) {
	g := NewGraphStore(2, 2)

	g.Nodes[1] = &Node{
		ID:        1,
		Vector:    []float32{10, 10},
		Neighbors: []int{2},
	}

	g.Nodes[2] = &Node{
		ID:        2,
		Vector:    []float32{1, 1},
		Neighbors: []int{},
	}

	g.LastID = 1

	query := []float32{0, 0}

	bestID, visited, err := g.Traverse(query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bestID != 2 {
		t.Fatalf("expected traversal to end at node 2, got %d", bestID)
	}

	if len(visited) != 2 {
		t.Fatalf("expected 2 visited nodes, got %d", len(visited))
	}
}

func TestTraverseNoImprovement(t *testing.T) {
	g := NewGraphStore(2, 2)

	g.Nodes[1] = &Node{
		ID:        1,
		Vector:    []float32{0, 0},
		Neighbors: []int{2},
	}

	g.Nodes[2] = &Node{
		ID:        2,
		Vector:    []float32{10, 10},
		Neighbors: []int{},
	}

	g.LastID = 1

	query := []float32{0, 0}

	bestID, visited, err := g.Traverse(query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bestID != 1 {
		t.Fatalf("expected traversal to remain at node 1, got %d", bestID)
	}

	if len(visited) != 2 {
		t.Fatalf("expected 2 visited nodes, got %d", len(visited))
	}
}

func TestTraverseMultiHop(t *testing.T) {
	g := NewGraphStore(2, 2)

	g.Nodes[1] = &Node{
		ID:        1,
		Vector:    []float32{10, 10},
		Neighbors: []int{2},
	}

	g.Nodes[2] = &Node{
		ID:        2,
		Vector:    []float32{5, 5},
		Neighbors: []int{3},
	}

	g.Nodes[3] = &Node{
		ID:        3,
		Vector:    []float32{1, 1},
		Neighbors: []int{},
	}

	g.LastID = 1

	query := []float32{0, 0}

	bestID, visited, err := g.Traverse(query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bestID != 3 {
		t.Fatalf("expected traversal to end at node 3, got %d", bestID)
	}

	if len(visited) != 3 {
		t.Fatalf("expected 3 visited nodes, got %d", len(visited))
	}
}

func TestTraverseRecordsDistances(t *testing.T) {
	g := NewGraphStore(2, 2)

	g.Nodes[1] = &Node{
		ID:        1,
		Vector:    []float32{3, 4},
		Neighbors: []int{},
	}

	g.LastID = 1

	query := []float32{0, 0}

	_, visited, err := g.Traverse(query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	distance, ok := visited[1]
	if !ok {
		t.Fatal("expected node 1 distance to be recorded")
	}

	if distance <= 0 {
		t.Fatalf("expected positive distance, got %f", distance)
	}
}
