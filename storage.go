package vectorengine

import "fmt"

type VectorStore struct {
	dimension int
	vectors   [][]float32
}

func NewVectorStore(dim int, capacity int) *VectorStore {
	return &VectorStore{
		dimension: dim,
		vectors:   make([][]float32, 0, capacity),
	}
}

func (vs *VectorStore) insert(vec []float32) error {
	if len(vec) != vs.dimension {
		return fmt.Errorf("Expected dimension %d, got %d dimension", vs.dimension, len(vec))
	}
	vs.vectors = append(vs.vectors, vec)
	return nil
}
