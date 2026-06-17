package vectorengine

import "testing"

func TestNewVectorStore(t *testing.T) {
	const dim = 128
	const capacity = 100000

	store := NewVectorStore(dim, capacity)

	if store.dimension != dim {
		t.Errorf("expected dimension %d, got %d", dim, store.dimension)
	}

	if len(store.vectors) != 0 {
		t.Errorf("expected length 0, got %d", len(store.vectors))
	}

	if cap(store.vectors) != capacity {
		t.Errorf("expected capacity %d, got %d", capacity, cap(store.vectors))
	}
}
