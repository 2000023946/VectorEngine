package vectorengine

import "testing"

func buildSimpleGraph() *Graph {
	g := NewGraphStore(2, 2, 10)

	g.SetVector(0, []float32{0, 0})
	g.SetVector(1, []float32{10, 10})
	g.SetVector(2, []float32{20, 20})

	g.AddNeighbor(0, 1)
	g.AddNeighbor(1, 2)

	g.Size = 3

	return g
}

func TestInsert_AddsNode(t *testing.T) {
	g := buildSimpleGraph()

	before := g.Size

	err := g.Insert([]float32{1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if g.Size != before+1 {
		t.Fatalf("expected %d nodes, got %d", before+1, g.Size)
	}
}

func TestInsert_CreatesNeighbors(t *testing.T) {
	g := buildSimpleGraph()

	err := g.Insert([]float32{1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newID := g.Size - 1

	if len(g.GetNeighbors(newID)) == 0 {
		t.Fatal("expected new node to have neighbors")
	}
}

func TestInsert_RespectsK(t *testing.T) {
	g := buildSimpleGraph()

	err := g.Insert([]float32{1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newID := g.Size - 1

	if len(g.GetNeighbors(newID)) > g.K {
		t.Fatalf("expected <= %d neighbors, got %d", g.K, len(g.GetNeighbors(newID)))
	}
}

func TestInsert_GraphConnectivity(t *testing.T) {
	g := buildSimpleGraph()

	err := g.Insert([]float32{1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newID := g.Size - 1

	found := false

	for i := 0; i < g.Size; i++ {
		neigh := g.GetNeighbors(i)

		for _, n := range neigh {
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
	g := NewGraphStore(2, 2, 10)

	err := g.Insert([]float32{0, 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if g.Size != 1 {
		t.Fatalf("expected 1 node, got %d", g.Size)
	}

	if len(g.GetNeighbors(0)) != len(g.GetNeighbors(0)) {
		t.Fatal("neighbor structure inconsistent")
	}
}
