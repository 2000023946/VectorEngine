package src

import (
	"fmt"
	"math"
	"sort"
)

const (
	maxVectors = 2_000_000
	dimension  = 128
)

type Vector struct {
	ID     int
	Offset int
}

type Result struct {
	ID       int
	Distance float64
}

type VectorNode struct {
	Vector    Vector
	Neighbors []int
}

type VectorEngine struct {
	// One large contiguous allocation:
	//
	// 2,000,000 vectors * 128 dimensions * 8 bytes
	// = 2.048 GB
	vectorData []float64

	// Vector metadata only.
	vectors []Vector

	// Graph nodes.
	nodes map[int]*VectorNode

	dimension int
	count     int
	entryID   int
}

func NewVectorEngine() *VectorEngine {
	return &VectorEngine{
		// Preallocate ~2 GB of contiguous vector memory.
		vectorData: make([]float64, maxVectors*dimension),

		// Preallocate metadata for all vectors.
		vectors: make([]Vector, 0, maxVectors),

		// Preallocate the hash map for the expected number of vectors.
		nodes: make(map[int]*VectorNode, maxVectors),

		dimension: dimension,
	}
}

func (ve *VectorEngine) Insert(id int, values []float64) error {
	// The first vector establishes the dimension
	// of the vector space.
	if ve.count == 0 {
		ve.dimension = len(values)
	} else if len(values) != ve.dimension {
		return fmt.Errorf(
			"dimension mismatch: expected %d, got %d",
			ve.dimension,
			len(values),
		)
	}

	// Make sure we don't exceed our preallocated buffer.
	if ve.count >= maxVectors {
		return fmt.Errorf(
			"vector capacity exceeded: maximum %d vectors",
			maxVectors,
		)
	}

	// --------------------------------------------------
	// Store vector in the flat contiguous buffer.
	// --------------------------------------------------
	//
	// Vector i occupies:
	//
	// [i*dimension : (i+1)*dimension]
	//
	// No allocation occurs here.
	// --------------------------------------------------

	offset := ve.count * ve.dimension

	copy(
		ve.vectorData[offset:offset+ve.dimension],
		values,
	)

	vector := Vector{
		ID:     id,
		Offset: offset,
	}

	node := &VectorNode{
		Vector:    vector,
		Neighbors: make([]int, 0),
	}

	// First vector becomes the graph entry point.
	if ve.count == 0 {
		ve.entryID = id
	}

	ve.vectors = append(ve.vectors, vector)
	ve.nodes[id] = node
	ve.count++

	// Nothing to connect for the first vector.
	if ve.count == 1 {
		return nil
	}

	// --------------------------------------------------
	// Graph construction
	// --------------------------------------------------
	//
	// Start from the entry point and walk through the
	// graph while a neighbor is closer to the new vector.
	// --------------------------------------------------

	currentID := ve.entryID
	current := ve.nodes[currentID]

	currentDistance := ve.distanceToVector(
		values,
		current.Vector,
	)

	for {
		nextID := currentID
		nextDistance := currentDistance

		// Check every neighbor of the current node.
		for _, neighborID := range current.Neighbors {
			neighbor := ve.nodes[neighborID]

			d := ve.distanceToVector(
				values,
				neighbor.Vector,
			)

			if d < nextDistance {
				nextID = neighborID
				nextDistance = d
			}
		}

		// No neighbor is closer.
		if nextID == currentID {
			break
		}

		// Move through the graph.
		currentID = nextID
		current = ve.nodes[currentID]
		currentDistance = nextDistance
	}

	// Connect the new vector to the node where
	// the greedy search stopped.
	current.Neighbors = append(
		current.Neighbors,
		id,
	)

	// Add the reverse connection so the graph can
	// navigate back toward the new node.
	node.Neighbors = append(
		node.Neighbors,
		currentID,
	)

	return nil
}

// vectorValues returns the contiguous slice containing
// the vector's values.
//
// No allocation occurs here.
func (ve *VectorEngine) vectorValues(vector Vector) []float64 {
	return ve.vectorData[vector.Offset : vector.Offset+ve.dimension]
}

func (ve *VectorEngine) distanceToVector(
	values []float64,
	vector Vector,
) float64 {
	var sum float64

	offset := vector.Offset

	for i := 0; i < ve.dimension; i++ {
		diff := values[i] - ve.vectorData[offset+i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

func (ve *VectorEngine) Search(query []float64, k int) []Result {
	if len(query) != ve.dimension {
		return nil
	}

	if ve.count == 0 || k <= 0 {
		return []Result{}
	}

	// --------------------------------------------------
	// Greedy graph search
	// --------------------------------------------------

	currentID := ve.entryID
	current := ve.nodes[currentID]

	currentDistance := ve.distanceToVector(
		query,
		current.Vector,
	)

	visited := make(map[int]bool)

	visited[currentID] = true

	for {
		nextID := currentID
		nextDistance := currentDistance

		for _, neighborID := range current.Neighbors {
			if visited[neighborID] {
				continue
			}

			neighbor := ve.nodes[neighborID]

			d := ve.distanceToVector(
				query,
				neighbor.Vector,
			)

			visited[neighborID] = true

			if d < nextDistance {
				nextID = neighborID
				nextDistance = d
			}
		}

		if nextID == currentID {
			break
		}

		currentID = nextID
		current = ve.nodes[currentID]
		currentDistance = nextDistance
	}

	// --------------------------------------------------
	// Return the best result discovered by the graph.
	// --------------------------------------------------

	results := make([]Result, 0, len(visited))

	for id := range visited {
		node := ve.nodes[id]

		results = append(results, Result{
			ID: id,
			Distance: ve.distanceToVector(
				query,
				node.Vector,
			),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	if k > len(results) {
		k = len(results)
	}

	return results[:k]
}
