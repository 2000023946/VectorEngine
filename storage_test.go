package vectorengine

import (
	"testing"
)

// -------------------- HELPERS --------------------

func newTestGraph() *Graph {
	return NewGraphStore(3, 2, 5)
}

// -------------------- ENTRY POINT TESTS --------------------

func TestEntryPointDefault(t *testing.T) {
	g := newTestGraph()

	if g.EntryPoint != 1 {
		t.Fatalf("expected EntryPoint=1 got %d", g.EntryPoint)
	}
}

func TestEntryPointLevelDefault(t *testing.T) {
	g := newTestGraph()

	if g.EntryPointLevel < 0 {
		t.Fatalf("invalid EntryPointLevel: %d", g.EntryPointLevel)
	}

	if g.EntryPointLevel >= g.MaxLevels {
		t.Fatalf("EntryPointLevel exceeds MaxLevels")
	}
}

func TestEntryPointConsistency(t *testing.T) {
	g := newTestGraph()

	if g.EntryPoint <= 0 || g.EntryPoint > g.Capacity {
		t.Fatalf("EntryPoint out of range: %d", g.EntryPoint)
	}
}

// -------------------- VECTOR TESTS --------------------

func TestSetGetVector(t *testing.T) {
	g := newTestGraph()

	vec := []float32{1, 2, 3}
	g.SetVector(1, vec)

	got := g.GetVector(1)

	for i := range vec {
		if got[i] != vec[i] {
			t.Fatalf("vector mismatch: expected %v got %v", vec, got)
		}
	}
}

func TestMultipleVectors(t *testing.T) {
	g := newTestGraph()

	g.SetVector(1, []float32{1, 1, 1})
	g.SetVector(2, []float32{2, 2, 2})

	v1 := g.GetVector(1)
	v2 := g.GetVector(2)

	if v1[0] != 1 || v2[0] != 2 {
		t.Fatalf("vector overwrite issue: v1=%v v2=%v", v1, v2)
	}
}

// -------------------- INDEXING TESTS --------------------

func TestIndexMonotonicAcrossLayers(t *testing.T) {
	g := newTestGraph()

	i0 := g.getIndex(1, 0)
	i1 := g.getIndex(1, 1)

	if i1 <= i0 {
		t.Fatalf("layer indexing incorrect: i0=%d i1=%d", i0, i1)
	}
}

func TestGetIndexSafeValid(t *testing.T) {
	g := newTestGraph()

	_, err := g.getIndexSafe(1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetIndexSafeInvalidNode(t *testing.T) {
	g := newTestGraph()

	_, err := g.getIndexSafe(999, 0)
	if err == nil {
		t.Fatalf("expected error for invalid node")
	}
}

func TestGetIndexSafeZeroNode(t *testing.T) {
	g := newTestGraph()

	_, err := g.getIndexSafe(0, 0)
	if err == nil {
		t.Fatalf("expected error for node 0")
	}
}

func TestMaxLevelsBoundary(t *testing.T) {
	g := newTestGraph()

	if g.MaxLevels <= 0 {
		t.Fatalf("MaxLevels not initialized properly")
	}

	_, err := g.getIndexSafe(1, g.MaxLevels)
	if err == nil {
		t.Fatalf("expected error for layer >= MaxLevels")
	}
}

// -------------------- NEIGHBOR TESTS --------------------

func TestAddNeighborSingle(t *testing.T) {
	g := newTestGraph()

	g.AddNeighbor(1, 2, 0)

	n := g.GetNeighbors(1, 0)

	if n[0] != 2 {
		t.Fatalf("expected neighbor 2 got %v", n)
	}
}

func TestAddNeighborMultipleSlots(t *testing.T) {
	g := newTestGraph()

	g.AddNeighbor(1, 2, 0)
	g.AddNeighbor(1, 3, 0)

	n := g.GetNeighbors(1, 0)

	if n[0] != 2 || n[1] != 3 {
		t.Fatalf("neighbors not filled correctly: %v", n)
	}
}

func TestNeighborCapacityLimit(t *testing.T) {
	g := newTestGraph()

	g.AddNeighbor(1, 2, 0)
	g.AddNeighbor(1, 3, 0)
	g.AddNeighbor(1, 4, 0)

	n := g.GetNeighbors(1, 0)

	count := 0
	for _, v := range n {
		if v != 0 {
			count++
		}
	}

	if count != 2 {
		t.Fatalf("expected K=2 neighbors, got %d", count)
	}
}

// -------------------- SAFE ACCESS TESTS --------------------

func TestGetNeighborValueSuccess(t *testing.T) {
	g := newTestGraph()

	g.AddNeighbor(1, 99, 0)

	val, err := g.GetNeighborValue(1, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != 99 {
		t.Fatalf("expected 99 got %d", val)
	}
}

func TestGetNeighborValueEmptySlot(t *testing.T) {
	g := newTestGraph()

	_, err := g.GetNeighborValue(1, 0, 0)
	if err == nil {
		t.Fatalf("expected error for empty slot")
	}
}

func TestGetNeighborValueOutOfBoundsOffset(t *testing.T) {
	g := newTestGraph()

	_, err := g.GetNeighborValue(1, 0, 99)
	if err == nil {
		t.Fatalf("expected error for invalid offset")
	}
}

func TestGetNeighborValueInvalidNode(t *testing.T) {
	g := newTestGraph()

	_, err := g.GetNeighborValue(999, 0, 0)
	if err == nil {
		t.Fatalf("expected error for invalid node")
	}
}

// -------------------- LAYER TESTS --------------------

func TestLayerIsolation(t *testing.T) {
	g := newTestGraph()

	g.AddNeighbor(1, 2, 0)
	g.AddNeighbor(1, 3, 1)

	n0 := g.GetNeighbors(1, 0)
	n1 := g.GetNeighbors(1, 1)

	if n0[0] == n1[0] {
		t.Fatalf("layers not isolated")
	}
}

func TestMaxLevelsExists(t *testing.T) {
	g := newTestGraph()

	if g.MaxLevels <= 0 {
		t.Fatalf("invalid MaxLevels")
	}
}

// -------------------- RANDOM LEVEL TEST --------------------

func TestGenerateRandomLayer(t *testing.T) {
	g := newTestGraph()

	for i := 0; i < 200; i++ {
		l := g.GenerateRandomLayer()

		if l < 0 {
			t.Fatalf("invalid negative layer: %d", l)
		}

		if l >= g.MaxLevels {
			t.Fatalf("layer exceeds MaxLevels: %d", l)
		}
	}
}
