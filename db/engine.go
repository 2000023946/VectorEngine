package db

import (
	"fmt"
	"math"
	"math/rand"
)

type Phase int

const (
	PhaseWarmup Phase = iota
	PhaseIndexed
)

type VectorEngine struct {
	Phase      Phase
	Threshold  int
	RawVectors [][]float32   // Phase 1: Flat list for O(N) append
	Centroids  [][]float32   // Phase 2: Frozen routing points
	Buckets    [][][]float32 // Phase 2: Contiguous memory blocks
}

func NewVectorEngine(threshold int) *VectorEngine {
	return &VectorEngine{
		Phase:     PhaseWarmup,
		Threshold: threshold,
	}
}

// Insert handles routing data based on the current phase
func (e *VectorEngine) Insert(vec []float32) {
	if e.Phase == PhaseIndexed {
		// FAST PATH: Route to nearest centroid in O(K)
		cIdx := e.findNearestCentroid(vec)
		// Append to contiguous slice in memory
		e.Buckets[cIdx] = append(e.Buckets[cIdx], vec)
		return
	}

	// WARMUP PATH: Flat append
	e.RawVectors = append(e.RawVectors, vec)

	// Transition trigger
	if len(e.RawVectors) >= e.Threshold {
		fmt.Printf("Threshold of %d reached. Building IVF Index...\n", e.Threshold)
		e.buildIndex()
	}
}

// buildIndex executes Random Initialization and Lloyd's Algorithm
func (e *VectorEngine) buildIndex() {
	N := len(e.RawVectors)
	K := int(math.Sqrt(float64(N)))
	dim := len(e.RawVectors[0])

	// STEP 1: Random Initialization (Forgy Method)
	e.Centroids = make([][]float32, K)
	perm := rand.Perm(N)
	for i := 0; i < K; i++ {
		e.Centroids[i] = append([]float32(nil), e.RawVectors[perm[i]]...) // Deep copy
	}

	// STEP 2: Lloyd's Algorithm (Batch Update)
	maxIterations := 20
	tolerance := float32(1e-5) // Convergence threshold

	for iter := 0; iter < maxIterations; iter++ {
		// Global Assignment: Create temporary buckets
		tempBuckets := make([][][]float32, K)
		for _, vec := range e.RawVectors {
			cIdx := e.findNearestCentroid(vec)
			tempBuckets[cIdx] = append(tempBuckets[cIdx], vec)
		}

		// Global Update: Calculate new means
		maxShift := float32(0.0)
		for i := 0; i < K; i++ {
			if len(tempBuckets[i]) == 0 {
				continue // Skip empty buckets
			}

			newCentroid := make([]float32, dim)
			for _, vec := range tempBuckets[i] {
				for d := 0; d < dim; d++ {
					newCentroid[d] += vec[d]
				}
			}

			// Calculate the mathematical average
			for d := 0; d < dim; d++ {
				newCentroid[d] /= float32(len(tempBuckets[i]))
			}

			// Track how far this centroid moved
			shift := euclideanSq(e.Centroids[i], newCentroid)
			if shift > maxShift {
				maxShift = shift
			}

			// Update centroid position
			e.Centroids[i] = newCentroid
		}

		// Convergence Check
		if maxShift < tolerance {
			fmt.Printf("Lloyd's Algorithm converged after %d iterations.\n", iter+1)
			break
		}
	}

	// STEP 3: Freeze the Index
	e.Buckets = make([][][]float32, K)
	for _, vec := range e.RawVectors {
		cIdx := e.findNearestCentroid(vec)
		e.Buckets[cIdx] = append(e.Buckets[cIdx], vec)
	}

	// STEP 4: Memory Cleanup
	e.RawVectors = nil // Free the heap memory used during warmup
	e.Phase = PhaseIndexed
	fmt.Println("Index frozen. Engine is now in PhaseIndexed.")
}
