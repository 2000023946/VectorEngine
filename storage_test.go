package vectorengine

import (
	"testing"
)

func TestVectorSetGet(t *testing.T) {
	g := NewGraphStore(3, 2, 10)

	vec := []float32{1, 2, 3}
	g.SetVector(0, vec)

	got := g.GetVector(0)

	for i := 0; i < len(vec); i++ {
		if got[i] != vec[i] {
			t.Fatalf("vector mismatch at %d: got %f want %f", i, got[i], vec[i])
		}
	}
}

func TestGetNeighborsBounds(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	neighbors := g.GetNeighbors(0, 0)

	if len(neighbors) != g.K {
		t.Fatalf("expected %d neighbors, got %d", g.K, len(neighbors))
	}
}

func TestAddNeighborSingleLayer(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	g.AddNeighbor(0, 42, 0)

	n := g.GetNeighbors(0, 0)

	if n[0] != 42 {
		t.Fatalf("expected 42, got %v", n)
	}
}

func TestAddNeighborFillsOnlyK(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	g.AddNeighbor(0, 1, 0)
	g.AddNeighbor(0, 2, 0)
	g.AddNeighbor(0, 3, 0) // should be ignored

	n := g.GetNeighbors(0, 0)

	count := 0
	for _, v := range n {
		if v != 0 {
			count++
		}
	}

	if count > g.K {
		t.Fatalf("overflow: got %d, max %d", count, g.K)
	}
}

func TestLayerIsolation(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	// insert into layer 0 and layer 1
	g.AddNeighbor(0, 10, 0)
	g.AddNeighbor(0, 99, 1)

	l0 := g.GetNeighbors(0, 0)
	l1 := g.GetNeighbors(0, 1)

	if l0[0] != 10 {
		t.Fatalf("layer 0 broken: %v", l0)
	}

	if l1[0] != 99 {
		t.Fatalf("layer 1 broken: %v", l1)
	}
}

func TestNoCrossLayerCorruption(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	g.AddNeighbor(0, 1, 0)
	g.AddNeighbor(0, 2, 1)

	l0 := g.GetNeighbors(0, 0)
	l1 := g.GetNeighbors(0, 1)

	if l0[0] != 1 {
		t.Fatalf("layer 0 corrupted: %v", l0)
	}

	if l1[0] != 2 {
		t.Fatalf("layer 1 corrupted: %v", l1)
	}
}

func TestMultipleNodesIsolation(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	g.AddNeighbor(0, 11, 0)
	g.AddNeighbor(1, 22, 0)

	n0 := g.GetNeighbors(0, 0)
	n1 := g.GetNeighbors(1, 0)

	if n0[0] != 11 || n1[0] != 22 {
		t.Fatalf("node isolation failed: %v %v", n0, n1)
	}
}

func TestGenerateRandomLayer(t *testing.T) {
	g := NewGraphStore(2, 2, 100)

	for i := 0; i < 1000; i++ {
		l := g.GenerateRandomLayer()

		if l < 0 {
			t.Fatalf("negative layer: %d", l)
		}

		if l > 50 { // sanity bound, not strict HNSW distribution test
			t.Fatalf("layer too large: %d", l)
		}
	}
}

func TestStressInsert(t *testing.T) {
	g := NewGraphStore(2, 2, 50)

	for i := 0; i < 20; i++ {
		for j := 0; j < 2; j++ {
			g.AddNeighbor(i, j+1, j%2)
		}
	}

	for i := 0; i < 20; i++ {
		for layer := 0; layer < 2; layer++ {
			n := g.GetNeighbors(i, layer)

			if len(n) != g.K {
				t.Fatalf("invalid K size at node %d layer %d", i, layer)
			}
		}
	}
}
