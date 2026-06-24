package vectorengine

import (
	"reflect"
	"testing"
)

func TestVectorSetGet(t *testing.T) {
	g := NewGraphStore(3, 2, 10)

	vec := []float32{1.0, 2.0, 3.0}
	g.SetVector(0, vec)

	got := g.GetVector(0)

	if !reflect.DeepEqual(vec, got) {
		t.Fatalf("expected %v, got %v", vec, got)
	}
}

func TestMultipleVectors(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	v1 := []float32{1, 2}
	v2 := []float32{3, 4}
	v3 := []float32{5, 6}

	g.SetVector(0, v1)
	g.SetVector(1, v2)
	g.SetVector(2, v3)

	if !reflect.DeepEqual(v1, g.GetVector(0)) {
		t.Fatal("vector 0 mismatch")
	}
	if !reflect.DeepEqual(v2, g.GetVector(1)) {
		t.Fatal("vector 1 mismatch")
	}
	if !reflect.DeepEqual(v3, g.GetVector(2)) {
		t.Fatal("vector 2 mismatch")
	}
}

func TestAddNeighbor(t *testing.T) {
	g := NewGraphStore(2, 3, 10)

	// node 0 neighbors
	g.AddNeighbor(0, 1)
	g.AddNeighbor(0, 2)
	g.AddNeighbor(0, 3)

	got := g.GetNeighbors(0)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestNeighborCapacityLimit(t *testing.T) {
	g := NewGraphStore(2, 2, 10)

	g.AddNeighbor(0, 1)
	g.AddNeighbor(0, 2)
	g.AddNeighbor(0, 3) // should be ignored (K=2)

	got := g.GetNeighbors(0)
	expected := []int{1, 2}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("capacity overflow failed, got %v", got)
	}
}

func TestFlatVectorMemoryLayout(t *testing.T) {
	g := NewGraphStore(3, 2, 10)

	// manually write vectors
	g.SetVector(0, []float32{1, 2, 3})
	g.SetVector(1, []float32{4, 5, 6})

	raw := g.Vectors

	expected := []float32{
		1, 2, 3,
		4, 5, 6,
	}

	if !reflect.DeepEqual(raw[:6], expected) {
		t.Fatalf("flat layout broken: %v", raw[:6])
	}
}

func TestMassInsert(t *testing.T) {
	g := NewGraphStore(128, 10, 1000)

	vec := make([]float32, 128)
	for i := range vec {
		vec[i] = float32(i)
	}

	for i := 0; i < 1000; i++ {
		g.SetVector(i, vec)
		g.AddNeighbor(i, (i+1)%1000)
	}

	if g.Size != 0 {
		t.Log("Size not used in current version (expected)")
	}
}
