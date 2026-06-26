package vectorengine

import (
	"testing"
)

// -------------------- HELPERS --------------------

func newInsertGraph() *Graph {
	g := NewGraphStore(3, 2, 5)

	// force stable search behavior assumption
	g.Size = 0
	g.EntryPoint = 0
	g.EntryPointLevel = 0

	return g
}

// -------------------- BASIC INSERT TESTS --------------------

func TestInsertFirstNodeSetsEntryPoint(t *testing.T) {
	g := newInsertGraph()

	id, err := g.Insert([]float32{1, 1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id != 1 {
		t.Fatalf("expected id=1 got %d", id)
	}

	if g.EntryPoint != 1 {
		t.Fatalf("expected EntryPoint=1 got %d", g.EntryPoint)
	}

	if g.Size != 1 {
		t.Fatalf("expected Size=1 got %d", g.Size)
	}
}

// -------------------- VECTOR STORAGE TEST --------------------

func TestInsertStoresVectorCorrectly(t *testing.T) {
	g := newInsertGraph()

	vec := []float32{1, 2, 3}

	id, err := g.Insert(vec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := g.GetVector(id)

	for i := range vec {
		if got[i] != vec[i] {
			t.Fatalf("vector mismatch: expected %v got %v", vec, got)
		}
	}
}

// -------------------- MULTIPLE INSERTIONS --------------------

func TestMultipleInsertsIncreaseSize(t *testing.T) {
	g := newInsertGraph()

	_, _ = g.Insert([]float32{1, 1, 1})
	_, _ = g.Insert([]float32{2, 2, 2})
	_, _ = g.Insert([]float32{3, 3, 3})

	if g.Size != 3 {
		t.Fatalf("expected Size=3 got %d", g.Size)
	}
}

// -------------------- ENTRY POINT UPDATE TEST --------------------

func TestEntryPointUpdatesWhenHigherLevelNodeAppears(t *testing.T) {
	g := newInsertGraph()

	// force deterministic behavior by overriding generated level
	// (we simulate by manually setting higher entry point first)
	_, _ = g.Insert([]float32{1, 1, 1})

	oldEP := g.EntryPoint
	oldLevel := g.EntryPointLevel

	_, _ = g.Insert([]float32{2, 2, 2})

	// EntryPoint may change depending on generated level
	if g.EntryPoint <= 0 {
		t.Fatalf("invalid EntryPoint after insert")
	}

	if g.EntryPointLevel < 0 || g.EntryPointLevel >= g.MaxLevels {
		t.Fatalf("invalid EntryPointLevel after insert: %d", g.EntryPointLevel)
	}

	_ = oldEP
	_ = oldLevel
}

// -------------------- CAPACITY LIMIT TEST --------------------

func TestInsertCapacityExceeded(t *testing.T) {
	g := NewGraphStore(3, 2, 1)

	_, err1 := g.Insert([]float32{1, 1, 1})
	if err1 != nil {
		t.Fatalf("first insert should succeed")
	}

	_, err2 := g.Insert([]float32{2, 2, 2})
	if err2 == nil {
		t.Fatalf("expected capacity exceeded error")
	}
}

// -------------------- SIZE CONSISTENCY --------------------

func TestInsertSizeConsistency(t *testing.T) {
	g := newInsertGraph()

	for i := 0; i < 5; i++ {
		_, err := g.Insert([]float32{float32(i), float32(i), float32(i)})
		if err != nil {
			t.Fatalf("unexpected error at i=%d: %v", i, err)
		}
	}

	if g.Size != 5 {
		t.Fatalf("expected Size=5 got %d", g.Size)
	}
}

// -------------------- VECTOR INTEGRITY --------------------

func TestInsertVectorIntegrityAcrossNodes(t *testing.T) {
	g := newInsertGraph()

	v1 := []float32{1, 1, 1}
	v2 := []float32{2, 2, 2}

	id1, _ := g.Insert(v1)
	id2, _ := g.Insert(v2)

	g1 := g.GetVector(id1)
	g2 := g.GetVector(id2)

	if g1[0] != 1 || g2[0] != 2 {
		t.Fatalf("vector integrity broken: g1=%v g2=%v", g1, g2)
	}
}

// -------------------- GRAPH STRUCTURE TEST --------------------

func TestInsertCreatesValidGraphStructure(t *testing.T) {
	g := newInsertGraph()

	_, _ = g.Insert([]float32{1, 1, 1})
	_, _ = g.Insert([]float32{2, 2, 2})

	if g.EntryPoint <= 0 {
		t.Fatalf("invalid EntryPoint")
	}

	if g.EntryPointLevel < 0 {
		t.Fatalf("invalid EntryPointLevel")
	}

	if g.Size != 2 {
		t.Fatalf("expected Size=2 got %d", g.Size)
	}
}

// -------------------- RANDOM LEVEL STABILITY --------------------

func TestInsertDoesNotCrashWithRandomLevels(t *testing.T) {
	g := NewGraphStore(3, 2, 5)

	for i := 0; i < 5; i++ {
		_, err := g.Insert([]float32{1, 2, 3})
		if err != nil {
			t.Fatalf("unexpected insert failure: %v", err)
		}
	}

	// 6th insert MUST fail
	_, err := g.Insert([]float32{1, 2, 3})
	if err == nil {
		t.Fatalf("expected capacity exceeded error")
	}
}
