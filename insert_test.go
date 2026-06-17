package vectorengine

import "testing"

func TestInsert(t *testing.T) {
	store := NewVectorStore(3, 3)
	vec := []float32{1, 2, 3}

	err := store.insert(vec)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

}

func TestInvalidInsert(t *testing.T) {
	store := NewVectorStore(4, 3)

	vec := []float32{1, 2, 3}

	err := store.insert(vec)

	if err == nil {
		t.Fatalf("expected dimension error, got nil")
	}
}
