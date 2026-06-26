package vectorengine

import (
	"math"
	"testing"
)

// -------------------- TEST HELPERS --------------------

func makeSimpleGraph() *Graph {
	g := NewGraphStore(2, 2, 5) // dim=2, K=2, maxNodes=5

	// vectors designed for predictable distance behavior
	// node 1 = origin
	g.SetVector(1, []float32{0, 0})

	// node 2 = slightly closer to query
	g.SetVector(2, []float32{1, 1})

	// node 3 = farthest
	g.SetVector(3, []float32{10, 10})

	return g
}

// helper distance for validation
func dist(a, b []float32) float64 {
	sum := 0.0
	for i := range a {
		d := float64(a[i] - b[i])
		sum += d * d
	}
	return math.Sqrt(sum)
}

// -------------------- BASIC SUCCESS PATH --------------------

func TestGreedyTraverseLayerImproves(t *testing.T) {
	g := makeSimpleGraph()

	// build graph connections:
	// 1 → 2 → 3 (chain)
	g.AddNeighbor(1, 2, 0)
	g.AddNeighbor(2, 3, 0)

	query := []float32{2, 2}

	result, err := g.GreedyTraverseLayer(1, query, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// should end closer to node 2 than node 1
	if result != 2 {
		t.Fatalf("expected to converge to node 2, got %d", result)
	}
}

// -------------------- NO IMPROVEMENT STOPS --------------------

func TestGreedyTraverseStopsWhenNoImprovement(t *testing.T) {
	g := makeSimpleGraph()

	// node 1 isolated (no neighbors)
	query := []float32{100, 100}

	result, err := g.GreedyTraverseLayer(1, query, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 1 {
		t.Fatalf("expected to stay at node 1, got %d", result)
	}
}

// -------------------- BEST NEIGHBOR SELECTION --------------------

func TestGreedyTraverseSelectsBestNeighbor(t *testing.T) {
	g := makeSimpleGraph()

	// node 1 connects to both node 2 and 3
	g.AddNeighbor(1, 2, 0)
	g.AddNeighbor(1, 3, 0)

	query := []float32{1, 1}

	result, err := g.GreedyTraverseLayer(1, query, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// node 2 is closer than node 3
	if result != 2 {
		t.Fatalf("expected node 2 (closer), got %d", result)
	}
}

// -------------------- MULTI-STEP DESCENT --------------------

func TestGreedyTraverseMultiStep(t *testing.T) {
	g := makeSimpleGraph()

	// chain: 1 → 2 → 3
	g.AddNeighbor(1, 2, 0)
	g.AddNeighbor(2, 3, 0)

	query := []float32{9, 9}

	result, err := g.GreedyTraverseLayer(1, query, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// should reach deepest node in chain
	if result != 3 {
		t.Fatalf("expected node 3, got %d", result)
	}
}

// -------------------- EMPTY NEIGHBOR SAFETY --------------------

func TestGreedyTraverseEmptyNeighbors(t *testing.T) {
	g := makeSimpleGraph()

	// no neighbors at all
	query := []float32{1, 1}

	result, err := g.GreedyTraverseLayer(1, query, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 1 {
		t.Fatalf("expected to remain at node 1, got %d", result)
	}
}

// -------------------- DISTANCE INTEGRITY CHECK --------------------

func TestGreedyTraverseActuallyImprovesDistance(t *testing.T) {
	g := makeSimpleGraph()

	g.AddNeighbor(1, 2, 0)

	query := []float32{2, 2}

	startDist := dist(g.GetVector(1), query)

	result, err := g.GreedyTraverseLayer(1, query, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	endDist := dist(g.GetVector(result), query)

	if endDist > startDist {
		t.Fatalf("greedy traversal did not improve distance")
	}
}

// -------------------- INVALID QUERY SAFETY --------------------

func TestGreedyTraverseInvalidVectorsHandled(t *testing.T) {
	g := makeSimpleGraph()

	// invalid query size (should trigger EuclideanDistance error depending on your implementation)
	query := []float32{1}

	_, err := g.GreedyTraverseLayer(1, query, 0)
	if err == nil {
		t.Fatalf("expected error for invalid query dimension")
	}
}
