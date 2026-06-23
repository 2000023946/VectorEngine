package vectorengine

import "testing"

// helper to build a simple graph
func buildSimpleGraph() *Graph {
	g := NewGraphStore(2, 2)

	g.Nodes[0] = &Node{
		ID:        0,
		Vector:    []float32{0, 0},
		Neighbors: []int{1},
	}

	g.Nodes[1] = &Node{
		ID:        1,
		Vector:    []float32{10, 10},
		Neighbors: []int{2},
	}

	g.Nodes[2] = &Node{
		ID:        2,
		Vector:    []float32{20, 20},
		Neighbors: []int{},
	}

	g.LastID = 2

	return g
}

func TestInsert_AddsNode(t *testing.T) {
	g := buildSimpleGraph()

	before := len(g.Nodes)

	err := g.Insert([]float32{1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(g.Nodes) != before+1 {
		t.Fatalf("expected %d nodes, got %d", before+1, len(g.Nodes))
	}

	if g.LastID != 3 {
		t.Fatalf("expected LastID to be 3, got %d", g.LastID)
	}
}

func TestInsert_CreatesNeighbors(t *testing.T) {
	g := buildSimpleGraph()

	err := g.Insert([]float32{1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newNode := g.Nodes[g.LastID]

	if len(newNode.Neighbors) == 0 {
		t.Fatal("expected new node to have neighbors")
	}
}

func TestInsert_RespectsK(t *testing.T) {
	g := buildSimpleGraph()

	err := g.Insert([]float32{1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newNode := g.Nodes[g.LastID]

	if len(newNode.Neighbors) > g.K {
		t.Fatalf("expected <= %d neighbors, got %d", g.K, len(newNode.Neighbors))
	}
}

func TestInsert_GraphConnectivity(t *testing.T) {
	g := buildSimpleGraph()

	err := g.Insert([]float32{1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newID := g.LastID

	found := false

	for _, node := range g.Nodes {
		for _, n := range node.Neighbors {
			if n == newID {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Fatal("expected reverse connectivity to new node")
	}
}

func TestInsert_NoPanic(t *testing.T) {
	g := buildSimpleGraph()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic occurred: %v", r)
		}
	}()

	_ = g.Insert([]float32{2, 2})
}

func TestInsert_FirstNodeBehavior(t *testing.T) {
	g := NewGraphStore(2, 2)

	err := g.Insert([]float32{0, 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if g.LastID != 0 {
		t.Fatalf("expected LastID = 0, got %d", g.LastID)
	}

	if len(g.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(g.Nodes))
	}
}
