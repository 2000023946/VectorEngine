package db

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
)

const (
	PhaseWarmup  = 0
	PhaseIndexed = 1
)

// IndexState holds the active database snapshot.
// Once created, it is IMMUTABLE. Never modify its slices directly.
type IndexState struct {
	Phase       int
	Centroids   [][]float32
	Buckets     [][][]float32
	FlatList    [][]float32
	SnapshotLen int // Tracks how many vectors are in this specific snapshot
}

type VectorEngine struct {
	// Replaces RWMutex for all read operations
	state atomic.Pointer[IndexState]

	// writeMu only protects RawVectors during inserts
	writeMu    sync.Mutex
	RawVectors [][]float32

	Threshold  int
	K          int
	isBuilding atomic.Bool // Prevents multiple background builds from stacking up
}

// NewVectorEngine initializes the database in PhaseWarmup
func NewVectorEngine(threshold, k int) *VectorEngine {
	e := &VectorEngine{
		Threshold: threshold,
		K:         k,
	}

	initialState := &IndexState{
		Phase:       PhaseWarmup,
		FlatList:    make([][]float32, 0),
		SnapshotLen: 0,
	}
	e.state.Store(initialState)

	return e
}

// Insert appends a new vector and determines if an index build is needed
func (e *VectorEngine) Insert(vec []float32) {
	e.writeMu.Lock()
	e.RawVectors = append(e.RawVectors, vec)
	currentLen := len(e.RawVectors)

	// If we are still in Warmup, we must update the state so Searches see the new vector
	currentState := e.state.Load()
	if currentState.Phase == PhaseWarmup {
		// Copy-on-Write for the FlatList to ensure thread safety
		newFlatList := make([][]float32, currentLen)
		copy(newFlatList, e.RawVectors)

		newState := &IndexState{
			Phase:       PhaseWarmup,
			FlatList:    newFlatList,
			SnapshotLen: currentLen,
		}
		e.state.Store(newState)
	}
	e.writeMu.Unlock()

	// If we crossed the threshold, trigger the background build.
	// isBuilding ensures we don't spawn 500 routines if Locust hits us hard.
	if currentLen >= e.Threshold && e.isBuilding.CompareAndSwap(false, true) {
		go e.buildIndex()
	}
}

// Search executes completely lock-free using the atomic pointer
func (e *VectorEngine) Search(query []float32) []float32 {
	// Atomic load takes ~1ns. Zero queuing. Zero blocking.
	currentState := e.state.Load()

	// 1. Phase Warmup: Global Search
	if currentState.Phase == PhaseWarmup {
		return e.fallbackGlobalSearch(query, currentState.FlatList)
	}

	// 2. Phase Indexed: K-Means Search
	if len(currentState.Centroids) == 0 {
		return nil // Safety check
	}

	// Find the closest centroid
	bestCentroidIdx := 0
	minDist := float32(math.MaxFloat32)

	for i, centroid := range currentState.Centroids {
		dist := l2Distance(query, centroid)
		if dist < minDist {
			minDist = dist
			bestCentroidIdx = i
		}
	}

	// Retrieve the bucket
	targetBucket := currentState.Buckets[bestCentroidIdx]

	// Poison Pill check: If bucket is empty, fallback to FlatList
	if len(targetBucket) == 0 {
		return e.fallbackGlobalSearch(query, currentState.FlatList)
	}

	// Search inside the targeted bucket
	return e.fallbackGlobalSearch(query, targetBucket)
}

// buildIndex runs in the background. It takes its time without blocking users.
func (e *VectorEngine) buildIndex() {
	defer e.isBuilding.Store(false)

	// 1. Grab a quick snapshot of the data safely
	e.writeMu.Lock()
	snapshotLen := len(e.RawVectors)
	snapshot := make([][]float32, snapshotLen)
	copy(snapshot, e.RawVectors)
	e.writeMu.Unlock()

	// 2. Background Math: Initialize Centroids randomly from snapshot
	centroids := make([][]float32, e.K)
	for i := 0; i < e.K; i++ {
		centroids[i] = snapshot[rand.Intn(snapshotLen)]
	}

	// 3. Lloyd's Algorithm (1 iteration for simplicity, increase for accuracy)
	buckets := make([][][]float32, e.K)
	for _, vec := range snapshot {
		bestIdx := 0
		minDist := float32(math.MaxFloat32)

		for i, c := range centroids {
			dist := l2Distance(vec, c)
			if dist < minDist {
				minDist = dist
				bestIdx = i
			}
		}
		buckets[bestIdx] = append(buckets[bestIdx], vec)
	}

	// 4. The Atomic Swap: Instantly upgrade the database state
	newState := &IndexState{
		Phase:       PhaseIndexed,
		Centroids:   centroids,
		Buckets:     buckets,
		FlatList:    snapshot, // Keep snapshot as fallback for empty buckets
		SnapshotLen: snapshotLen,
	}

	e.state.Store(newState)
}

// Size returns the true number of inserted vectors
func (e *VectorEngine) Size() int {
	e.writeMu.Lock()
	defer e.writeMu.Unlock()
	return len(e.RawVectors)
}

// fallbackGlobalSearch is a helper for exhaustive searching
func (e *VectorEngine) fallbackGlobalSearch(query []float32, dataset [][]float32) []float32 {
	if len(dataset) == 0 {
		return nil
	}

	var bestMatch []float32
	minDist := float32(math.MaxFloat32)

	for _, vec := range dataset {
		dist := l2Distance(query, vec)
		if dist < minDist {
			minDist = dist
			bestMatch = vec
		}
	}
	return bestMatch
}

// l2Distance calculates squared Euclidean distance
func l2Distance(a, b []float32) float32 {
	var sum float32
	for i := 0; i < len(a) && i < len(b); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return sum
}
