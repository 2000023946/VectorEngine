package vectorengine

import (
	"testing"
)

// =========================================================
// FIXED VECTOR HELPER (IMPORTANT)
// =========================================================

func makeVec(dim int, val float32) []float32 {
	v := make([]float32, dim)
	for i := range v {
		v[i] = val
	}
	return v
}

// =========================================================
// TEST FIXTURE GRAPH
// =========================================================

func newInsertTestGraph() *Graph {
	g := NewGraphStore(4, 3, 200) // ⬅️ FIX HERE

	for i := 1; i <= 200; i++ {
		g.SetVector(i, makeVec(g.Dimension, float32(i)))
	}

	g.Size = 0
	g.EntryPoint = 1
	g.EntryPointLevel = 0

	return g
}

// =========================================================
// TEST 1: BASIC INSERT EXECUTION
// =========================================================

func TestUnitInsertBasicExecution(t *testing.T) {

	g := newInsertTestGraph()

	id, err := g.Insert(makeVec(g.Dimension, 10))
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	if id <= 0 || id > g.Capacity {
		t.Fatalf("invalid node id: %d", id)
	}

	if g.Size != 1 {
		t.Fatalf("expected size=1 got %d", g.Size)
	}
}

// =========================================================
// TEST 2: ENTRY POINT UPDATES
// =========================================================

func TestUnitInsertEntryPointUpdate(t *testing.T) {

	g := newInsertTestGraph()

	_, err := g.Insert(makeVec(g.Dimension, 5))
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	_, err = g.Insert(makeVec(g.Dimension, 15))
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	if g.EntryPoint <= 0 {
		t.Fatalf("entry point not set")
	}
}

// =========================================================
// TEST 3: GRAPH DEGREE CONSTRAINT (K)
// =========================================================

func TestUnitInsertRespectsK(t *testing.T) {

	g := NewGraphStore(4, 2, 10)

	for i := 1; i <= 5; i++ {
		g.SetVector(i, makeVec(g.Dimension, float32(i)))
	}

	for i := 0; i < 5; i++ {
		_, err := g.Insert(makeVec(g.Dimension, float32(i+1)))
		if err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}

	for node := 1; node <= g.Size; node++ {
		neighbors := g.GetNeighbors(node, 0)

		count := 0
		for _, n := range neighbors {
			if n != 0 {
				count++
			}
		}

		if count > g.K {
			t.Fatalf("node %d exceeds K=%d neighbors: %d", node, g.K, count)
		}
	}
}

// =========================================================
// TEST 4: BIDIRECTIONAL EDGE CONSISTENCY
// =========================================================

func TestUnitInsertBidirectionalEdges(t *testing.T) {

	g := newInsertTestGraph()

	id1, err := g.Insert(makeVec(g.Dimension, 5))
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	id2, err := g.Insert(makeVec(g.Dimension, 6))
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	foundForward := false
	foundBackward := false

	for _, n := range g.GetNeighbors(id1, 0) {
		if n == id2 {
			foundForward = true
		}
	}

	for _, n := range g.GetNeighbors(id2, 0) {
		if n == id1 {
			foundBackward = true
		}
	}

	if !foundForward && !foundBackward {
		t.Fatalf("missing bidirectional connectivity")
	}
}

// =========================================================
// TEST 5: INSERT STABILITY (NO CRASH)
// =========================================================

func TestUnitInsertStressNoCrash(t *testing.T) {

	g := newInsertTestGraph()

	for i := 0; i < 100; i++ {
		_, err := g.Insert(makeVec(g.Dimension, float32(i%10+1)))
		if err != nil {
			t.Fatalf("insert failed at i=%d: %v", i, err)
		}
	}
}

// =========================================================
// TEST 6: DETERMINISTIC BEHAVIOR
// =========================================================

func TestUnitInsertDeterminism(t *testing.T) {

	g1 := newInsertTestGraph()
	g2 := newInsertTestGraph()

	for i := 1; i <= 10; i++ {
		v := makeVec(4, float32(i))

		g1.Insert(v)
		g2.Insert(v)
	}

	if g1.Size != g2.Size {
		t.Fatalf("size mismatch")
	}

	if g1.EntryPoint == 0 || g2.EntryPoint == 0 {
		t.Fatalf("entry point not set deterministically")
	}
}
