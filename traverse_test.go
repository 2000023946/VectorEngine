package vectorengine

// import "testing"

// func buildTraverseGraph() *Graph {
// 	g := NewGraphStore(2, 2, 10)

// 	// IMPORTANT: use Insert so Size is correct
// 	_ = g.Insert([]float32{10, 10}) // node 0
// 	_ = g.Insert([]float32{1, 1})   // node 1

// 	g.AddNeighbor(0, 1)

// 	return g
// }

// func TestTraverseEmptyGraph(t *testing.T) {
// 	g := NewGraphStore(2, 2, 10)

// 	query := []float32{0, 0}

// 	id, visited, err := g.Traverse(query)

// 	if err == nil {
// 		t.Fatal("expected error for empty graph")
// 	}

// 	if id != -1 {
// 		t.Fatalf("expected id -1, got %d", id)
// 	}

// 	if visited != nil {
// 		t.Fatal("expected nil visited map")
// 	}
// }
// func TestTraverseSingleHop(t *testing.T) {
// 	g := buildTraverseGraph()

// 	query := []float32{0, 0}

// 	bestID, visited, err := g.Traverse(query)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if bestID != 1 {
// 		t.Fatalf("expected traversal to end at node 1, got %d", bestID)
// 	}

// 	if len(visited) != 2 {
// 		t.Fatalf("expected 2 visited nodes, got %d", len(visited))
// 	}
// }
// func TestTraverseNoImprovement(t *testing.T) {
// 	g := NewGraphStore(2, 2, 10)

// 	_ = g.Insert([]float32{0, 0})   // node 0
// 	_ = g.Insert([]float32{10, 10}) // node 1

// 	g.AddNeighbor(0, 1)

// 	query := []float32{0, 0}

// 	bestID, visited, err := g.Traverse(query)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if bestID != 0 {
// 		t.Fatalf("expected traversal to remain at node 0, got %d", bestID)
// 	}

// 	if len(visited) != 2 {
// 		t.Fatalf("expected 2 visited nodes, got %d", len(visited))
// 	}
// }

// func TestTraverseMultiHop(t *testing.T) {
// 	g := NewGraphStore(2, 2, 10)

// 	_ = g.Insert([]float32{10, 10}) // 0
// 	_ = g.Insert([]float32{5, 5})   // 1
// 	_ = g.Insert([]float32{1, 1})   // 2

// 	g.AddNeighbor(0, 1)
// 	g.AddNeighbor(1, 2)

// 	query := []float32{0, 0}

// 	bestID, visited, err := g.Traverse(query)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if bestID != 2 {
// 		t.Fatalf("expected traversal to end at node 2, got %d", bestID)
// 	}

// 	if len(visited) != 3 {
// 		t.Fatalf("expected 3 visited nodes, got %d", len(visited))
// 	}
// }

// func TestTraverseRecordsDistances(t *testing.T) {
// 	g := NewGraphStore(2, 2, 10)

// 	_ = g.Insert([]float32{3, 4}) // node 0

// 	query := []float32{0, 0}

// 	_, visited, err := g.Traverse(query)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	distance, ok := visited[0]
// 	if !ok {
// 		t.Fatal("expected node 0 distance to be recorded")
// 	}

// 	if distance <= 0 {
// 		t.Fatalf("expected positive distance, got %f", distance)
// 	}
// }
