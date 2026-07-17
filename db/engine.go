package db

import (
	"errors"
	"sync"
	"time"
)

// Vector represents a single motion window from the DE10-Lite
type Vector struct {
	ID        string
	Timestamp time.Time
	Data      []float32 // Explicit float32 to cut memory footprint in half
}

// VectorEngine is the in-memory database
type VectorEngine struct {
	mu      sync.RWMutex // Protects concurrent API reads/writes
	vectors []Vector
	dim     int // Dimensionality of the vector (e.g., 30 for 10x(X,Y,Z))
}

// NewEngine initializes the DB with a pre-allocated capacity
func NewEngine(capacity int, dimensions int) *VectorEngine {
	return &VectorEngine{
		// Pre-allocate the slice to prevent expensive dynamic resizing
		vectors: make([]Vector, 0, capacity),
		dim:     dimensions,
	}
}

// Insert safely adds a new vector to the engine
func (e *VectorEngine) Insert(v Vector) error {
	if len(v.Data) != e.dim {
		return errors.New("vector dimension mismatch")
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.vectors = append(e.vectors, v)
	return nil
}

// Count returns the current number of stored vectors
func (e *VectorEngine) Count() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.vectors)
}
