package db

import (
	"math"
	"math/rand"
	"sync"
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

	// Background Concurrency Fields
	insertQueue chan []float32 // Hardware-like FIFO queue
	mu          sync.RWMutex   // Protects reads/writes during phase transitions
}

func NewVectorEngine(threshold int, queueSize int) *VectorEngine {
	engine := &VectorEngine{
		Phase:       PhaseWarmup,
		Threshold:   threshold,
		insertQueue: make(chan []float32, queueSize),
	}

	// Spin up the background worker goroutine
	go engine.backgroundWorker()

	return engine
}

// Insert drops the vector into the queue and returns immediately (Non-blocking)
func (e *VectorEngine) Insert(vec []float32) {
	e.insertQueue <- vec
}

// backgroundWorker continuously processes vectors from the queue
func (e *VectorEngine) backgroundWorker() {
	for vec := range e.insertQueue {
		e.processInsert(vec)
	}
}

// processInsert executes routing logic inside the background thread
func (e *VectorEngine) processInsert(vec []float32) {
	e.mu.Lock()
	if e.Phase == PhaseIndexed {
		// FAST PATH: Route to nearest centroid in O(K)
		cIdx := e.findNearestCentroid(vec)
		e.Buckets[cIdx] = append(e.Buckets[cIdx], vec)
		e.mu.Unlock()
		return
	}

	// WARMUP PATH: Flat append
	e.RawVectors = append(e.RawVectors, vec)
	shouldBuild := len(e.RawVectors) >= e.Threshold
	e.mu.Unlock()

	// Transition trigger
	if shouldBuild {
		e.buildIndex()
	}
}

// Search routes the query safely under a Read Lock
func (e *VectorEngine) Search(query []float32) ([]float32, float32) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.Phase == PhaseIndexed {
		// 1. Find nearest frozen centroid
		cIdx := e.findNearestCentroid(query)

		// 2. Scan only that specific contiguous bucket
		return findNearestInList(query, e.Buckets[cIdx])
	}

	// Fallback: Brute force the raw list if still warming up
	return findNearestInList(query, e.RawVectors)
}

// buildIndex executes Lloyd's Algorithm without locking readers during computation
func (e *VectorEngine) buildIndex() {
	// 1. Create a snapshot copy of RawVectors under Read Lock
	e.mu.RLock()
	N := len(e.RawVectors)
	if N == 0 {
		e.mu.RUnlock()
		return
	}
	rawCopy := make([][]float32, N)
	copy(rawCopy, e.RawVectors)
	e.mu.RUnlock()

	K := int(math.Sqrt(float64(N)))
	dim := len(rawCopy[0])

	// STEP 1: Random Initialization (Forgy Method)
	centroids := make([][]float32, K)
	perm := rand.Perm(N)
	for i := 0; i < K; i++ {
		centroids[i] = append([]float32(nil), rawCopy[perm[i]]...)
	}

	// STEP 2: Lloyd's Algorithm (Runs completely unlocked)
	maxIterations := 20
	tolerance := float32(1e-5)

	for iter := 0; iter < maxIterations; iter++ {
		tempBuckets := make([][][]float32, K)
		for _, vec := range rawCopy {
			cIdx := findNearestInCentroids(vec, centroids)
			tempBuckets[cIdx] = append(tempBuckets[cIdx], vec)
		}

		maxShift := float32(0.0)
		for i := 0; i < K; i++ {
			if len(tempBuckets[i]) == 0 {
				continue
			}

			newCentroid := make([]float32, dim)
			for _, vec := range tempBuckets[i] {
				for d := 0; d < dim; d++ {
					newCentroid[d] += vec[d]
				}
			}

			for d := 0; d < dim; d++ {
				newCentroid[d] /= float32(len(tempBuckets[i]))
			}

			shift := euclideanSq(centroids[i], newCentroid)
			if shift > maxShift {
				maxShift = shift
			}

			centroids[i] = newCentroid
		}

		if maxShift < tolerance {
			break
		}
	}

	// STEP 3: Create final buckets
	buckets := make([][][]float32, K)
	for _, vec := range rawCopy {
		cIdx := findNearestInCentroids(vec, centroids)
		buckets[cIdx] = append(buckets[cIdx], vec)
	}

	// STEP 4: Lock briefly to swap state and enter PhaseIndexed
	e.mu.Lock()
	e.Centroids = centroids
	e.Buckets = buckets
	e.RawVectors = nil // Free memory used during warmup
	e.Phase = PhaseIndexed
	e.mu.Unlock()
}

// findNearestCentroid uses internal Centroids state
func (e *VectorEngine) findNearestCentroid(vec []float32) int {
	return findNearestInCentroids(vec, e.Centroids)
}

// Helper function to search arbitrary centroid slices
func findNearestInCentroids(vec []float32, centroids [][]float32) int {
	bestIdx := -1
	minDist := float32(math.MaxFloat32)
	for i, centroid := range centroids {
		dist := euclideanSq(vec, centroid)
		if dist < minDist {
			minDist = dist
			bestIdx = i
		}
	}
	return bestIdx
}

// findNearestInList powers both the Phase 1 scan and Phase 2 bucket scan
func findNearestInList(query []float32, list [][]float32) ([]float32, float32) {
	if len(list) == 0 {
		return nil, -1
	}

	var bestVec []float32
	minDist := float32(math.MaxFloat32)
	for _, v := range list {
		dist := euclideanSq(query, v)
		if dist < minDist {
			minDist = dist
			bestVec = v
		}
	}
	return bestVec, minDist
}

// euclideanSq calculates squared Euclidean distance
func euclideanSq(a, b []float32) float32 {
	var sum float32
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return sum
}
