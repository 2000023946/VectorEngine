package vectorengine

import "testing"

func TestUnitEuclideanDistance(t *testing.T) {
	a := []float32{0, 0}
	b := []float32{3, 4}

	dist, err := EuclideanDistance(a, b)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dist != 25 {
		t.Fatalf("expected 25, got %f", dist)
	}
}

func TestUnitEuclideanDistanceDimensionMismatch(t *testing.T) {
	a := []float32{1, 2}
	b := []float32{1, 2, 3}

	_, err := EuclideanDistance(a, b)

	if err == nil {
		t.Fatalf("expected dimension mismatch error")
	}
}
