package src

import (
	"fmt"
	"math"
	"sort"
)

type Vector struct {
	ID     int
	Values []float64
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
	vectors   []Vector
	nodes     map[int]*VectorNode
	dimension int
	entryID   int
}

func NewVectorEngine() *VectorEngine {
	return &VectorEngine{
		vectors: make([]Vector, 0),
		nodes:   make(map[int]*VectorNode),
	}
}

func (ve *VectorEngine) Insert(id int, values []float64) error {
	// The first vector establishes the dimension
	// of the vector space.
	if ve.dimension == 0 {
		ve.dimension = len(values)
	} else if len(values) != ve.dimension {
		return fmt.Errorf(
			"dimension mismatch: expected %d, got %d",
			ve.dimension,
			len(values),
		)
	}

	vector := Vector{
		ID:     id,
		Values: values,
	}

	node := &VectorNode{
		Vector:    vector,
		Neighbors: make([]int, 0),
	}

	// First vector becomes the graph entry point.
	if len(ve.vectors) == 0 {
		ve.entryID = id
	}

	ve.vectors = append(ve.vectors, vector)
	ve.nodes[id] = node

	// Nothing to connect for the first vector.
	if len(ve.vectors) == 1 {
		return nil
	}

	// --------------------------------------------------
	// Graph construction
	// --------------------------------------------------
	//
	// Start from the entry point and walk through the
	// graph while a neighbor is closer to the new vector.
	//
	// This is a simple greedy graph baseline inspired
	// by the search behavior of HNSW.
	// --------------------------------------------------

	currentID := ve.entryID
	current := ve.nodes[currentID]

	currentDistance := distance(values, current.Vector.Values)

	for {
		nextID := currentID
		nextDistance := currentDistance

		// Check every neighbor of the current node.
		for _, neighborID := range current.Neighbors {
			neighbor := ve.nodes[neighborID]

			d := distance(values, neighbor.Vector.Values)

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
	current.Neighbors = append(current.Neighbors, id)

	// Add the reverse connection so the graph can
	// navigate back toward the new node.
	node.Neighbors = append(node.Neighbors, currentID)

	return nil
}

func distance(a []float64, b []float64) float64 {
	var sum float64

	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

func (ve *VectorEngine) Search(query []float64, k int) []Result {
	if len(query) != ve.dimension {
		return nil
	}

	if len(ve.vectors) == 0 || k <= 0 {
		return []Result{}
	}

	// --------------------------------------------------
	// Greedy graph search
	// --------------------------------------------------

	currentID := ve.entryID
	current := ve.nodes[currentID]

	currentDistance := distance(query, current.Vector.Values)

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
			d := distance(query, neighbor.Vector.Values)

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
			ID:       id,
			Distance: distance(query, node.Vector.Values),
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
