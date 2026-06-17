package vectorengine

import "testing"

func TestSearch(t *testing.T) {

	store := NewVectorStore(2, 10)

	store.insert([]float32{1, 1})
	store.insert([]float32{5, 5})
	store.insert([]float32{10, 10})

	result, err := store.Search([]float32{2, 2})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Index != 0 {
		t.Fatalf("expected index 0, got %d", result.Index)
	}
}

func TestSearchEmptyStore(t *testing.T) {

	store := NewVectorStore(2, 10)

	_, err := store.Search([]float32{1, 1})

	if err == nil {
		t.Fatalf("expected error for empty store")
	}
}

func TestSearchWrongDimension(t *testing.T) {

	store := NewVectorStore(2, 10)

	store.insert([]float32{1, 1})

	_, err := store.Search([]float32{1, 1, 1})

	if err == nil {
		t.Fatalf("expected dimension error")
	}
}

func TestSearchExactMatch(t *testing.T) {

	store := NewVectorStore(2, 10)

	store.insert([]float32{3, 4})
	store.insert([]float32{7, 8})

	result, err := store.Search([]float32{3, 4})

	if err != nil {
		t.Fatalf("unexpected error")
	}

	if result.Index != 0 {
		t.Fatalf("expected index 0")
	}

	if result.Distance != 0 {
		t.Fatalf("expected distance 0")
	}
}

func TestSearchEuclideanErrorPropagation(t *testing.T) {

	store := NewVectorStore(2, 10)

	// valid vector first
	store.insert([]float32{1, 1})

	// ❗ inject invalid vector directly (bypasses Insert validation)
	store.vectors = append(store.vectors, []float32{1, 1, 1})

	_, err := store.Search([]float32{2, 2})

	if err == nil {
		t.Fatalf("expected Euclidean error propagation, got nil")
	}

	expected := "vectors must have same dimension"
	if err.Error() != expected {
		t.Fatalf("expected '%s', got '%s'", expected, err.Error())
	}
}
