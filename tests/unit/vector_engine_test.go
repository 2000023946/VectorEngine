package tests

import (
	"math"
	"testing"

	"vectorengine/src"
)

func TestInsertAndSearch(t *testing.T) {
	engine := src.NewVectorEngine()

	err := engine.Insert(1, []float64{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}

	err = engine.Insert(2, []float64{1.1, 2.1, 3.1})
	if err != nil {
		t.Fatal(err)
	}

	err = engine.Insert(3, []float64{10, 10, 10})
	if err != nil {
		t.Fatal(err)
	}

	query := []float64{1, 2, 3}

	results := engine.Search(query, 2)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// The query is identical to vector 1,
	// so vector 1 must be the closest result.
	if results[0].ID != 1 {
		t.Fatalf(
			"expected closest vector to be ID 1, got ID %d",
			results[0].ID,
		)
	}

	if results[0].Distance != 0 {
		t.Fatalf(
			"expected distance 0, got %f",
			results[0].Distance,
		)
	}
}

func TestSearchReturnsClosestVectors(t *testing.T) {
	engine := src.NewVectorEngine()

	vectors := []struct {
		id     int
		values []float64
	}{
		{1, []float64{1, 1}},
		{2, []float64{2, 2}},
		{3, []float64{10, 10}},
	}

	for _, vector := range vectors {
		if err := engine.Insert(vector.id, vector.values); err != nil {
			t.Fatal(err)
		}
	}

	query := []float64{1, 1}

	results := engine.Search(query, 2)

	if results[0].ID != 1 {
		t.Errorf("expected ID 1 first, got ID %d", results[0].ID)
	}

	if results[1].ID != 2 {
		t.Errorf("expected ID 2 second, got ID %d", results[1].ID)
	}
}

func TestDimensionMismatch(t *testing.T) {
	engine := src.NewVectorEngine()

	err := engine.Insert(1, []float64{1, 2, 3})

	if err != nil {
		t.Fatal(err)
	}

	// The engine has dimension 3.
	// This vector has dimension 4 and should be rejected.
	err = engine.Insert(2, []float64{1, 2, 3, 4})

	if err == nil {
		t.Fatal("expected dimension mismatch error")
	}
}

func TestQueryDimensionMismatch(t *testing.T) {
	engine := src.NewVectorEngine()

	err := engine.Insert(1, []float64{1, 2, 3})

	if err != nil {
		t.Fatal(err)
	}

	results := engine.Search(
		[]float64{1, 2},
		1,
	)

	if results != nil {
		t.Fatal("expected nil result for invalid query dimension")
	}
}

func TestKGreaterThanDatasetSize(t *testing.T) {
	engine := src.NewVectorEngine()

	err := engine.Insert(1, []float64{1, 2, 3})

	if err != nil {
		t.Fatal(err)
	}

	results := engine.Search(
		[]float64{1, 2, 3},
		10,
	)

	if len(results) != 1 {
		t.Fatalf(
			"expected 1 result, got %d",
			len(results),
		)
	}
}

func TestNegativeK(t *testing.T) {
	engine := src.NewVectorEngine()

	err := engine.Insert(1, []float64{1, 2, 3})

	if err != nil {
		t.Fatal(err)
	}

	results := engine.Search(
		[]float64{1, 2, 3},
		-1,
	)

	if len(results) != 0 {
		t.Fatalf(
			"expected 0 results, got %d",
			len(results),
		)
	}
}

func TestDistanceOrdering(t *testing.T) {
	engine := src.NewVectorEngine()

	err := engine.Insert(1, []float64{0, 0})
	if err != nil {
		t.Fatal(err)
	}

	err = engine.Insert(2, []float64{3, 4})
	if err != nil {
		t.Fatal(err)
	}

	results := engine.Search(
		[]float64{0, 0},
		2,
	)

	expectedDistance := 5.0

	if math.Abs(results[1].Distance-expectedDistance) > 0.000001 {
		t.Fatalf(
			"expected distance %f, got %f",
			expectedDistance,
			results[1].Distance,
		)
	}
}
