package tests

import (
	"testing"

	"vectorengine/src"
)

// makeVector creates a 128-dimensional vector with the
// first few dimensions populated and the rest set to 0.
func makeVector(values ...float64) []float64 {
	vector := make([]float64, 128)

	copy(vector, values)

	return vector
}

func TestInsertAndSearch(t *testing.T) {
	engine := src.NewVectorEngine()

	err := engine.Insert(
		1,
		makeVector(1, 2, 3),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = engine.Insert(
		2,
		makeVector(1.1, 2.1, 3.1),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = engine.Insert(
		3,
		makeVector(10, 10, 10),
	)
	if err != nil {
		t.Fatal(err)
	}

	query := makeVector(1, 2, 3)

	results := engine.Search(query, 2)

	if len(results) != 2 {
		t.Fatalf(
			"expected 2 results, got %d",
			len(results),
		)
	}

	// Query is identical to vector 1,
	// so vector 1 must be the closest result.
	if results[0].ID != 1 {
		t.Fatalf(
			"expected closest vector to be ID 1, got ID %d",
			results[0].ID,
		)
	}

	if results[0].Distance != 0 {
		t.Fatalf(
			"expected squared distance 0, got %f",
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
		{
			1,
			makeVector(1, 1),
		},
		{
			2,
			makeVector(2, 2),
		},
		{
			3,
			makeVector(10, 10),
		},
	}

	for _, vector := range vectors {
		if err := engine.Insert(
			vector.id,
			vector.values,
		); err != nil {
			t.Fatal(err)
		}
	}

	query := makeVector(1, 1)

	results := engine.Search(query, 2)

	if len(results) != 2 {
		t.Fatalf(
			"expected 2 results, got %d",
			len(results),
		)
	}

	if results[0].ID != 1 {
		t.Errorf(
			"expected ID 1 first, got ID %d",
			results[0].ID,
		)
	}

	if results[1].ID != 2 {
		t.Errorf(
			"expected ID 2 second, got ID %d",
			results[1].ID,
		)
	}
}

func TestDimensionMismatch(t *testing.T) {
	engine := src.NewVectorEngine()

	// Correct dimension: 128.
	err := engine.Insert(
		1,
		makeVector(1, 2, 3),
	)

	if err != nil {
		t.Fatal(err)
	}

	// Incorrect dimension: 127.
	invalidVector := make([]float64, 127)

	err = engine.Insert(
		2,
		invalidVector,
	)

	if err == nil {
		t.Fatal("expected dimension mismatch error")
	}
}

func TestQueryDimensionMismatch(t *testing.T) {
	engine := src.NewVectorEngine()

	err := engine.Insert(
		1,
		makeVector(1, 2, 3),
	)

	if err != nil {
		t.Fatal(err)
	}

	// Query has 127 dimensions instead of 128.
	invalidQuery := make([]float64, 127)

	results := engine.Search(
		invalidQuery,
		1,
	)

	if results != nil {
		t.Fatal(
			"expected nil result for invalid query dimension",
		)
	}
}

func TestKGreaterThanDatasetSize(t *testing.T) {
	engine := src.NewVectorEngine()

	err := engine.Insert(
		1,
		makeVector(1, 2, 3),
	)

	if err != nil {
		t.Fatal(err)
	}

	results := engine.Search(
		makeVector(1, 2, 3),
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

	err := engine.Insert(
		1,
		makeVector(1, 2, 3),
	)

	if err != nil {
		t.Fatal(err)
	}

	results := engine.Search(
		makeVector(1, 2, 3),
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

	err := engine.Insert(
		1,
		makeVector(0, 0),
	)

	if err != nil {
		t.Fatal(err)
	}

	err = engine.Insert(
		2,
		makeVector(3, 4),
	)

	if err != nil {
		t.Fatal(err)
	}

	results := engine.Search(
		makeVector(0, 0),
		2,
	)

	if len(results) != 2 {
		t.Fatalf(
			"expected 2 results, got %d",
			len(results),
		)
	}

	if results[0].ID != 1 {
		t.Fatalf(
			"expected ID 1 first, got ID %d",
			results[0].ID,
		)
	}

	// We now use squared Euclidean distance:
	//
	// 3² + 4² = 25
	//
	// No sqrt is performed.
	expectedSquaredDistance := 25.0

	if results[1].Distance != expectedSquaredDistance {
		t.Fatalf(
			"expected squared distance %f, got %f",
			expectedSquaredDistance,
			results[1].Distance,
		)
	}
}
