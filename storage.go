package vectorengine

import "fmt"

type Node struct {
	ID        int
	Vector    []float32
	Neighbors []int
}

type Graph struct {
	Nodes     map[int]*Node
	K         int
	Dimension int
	LastID    int
}

func NewGraphStore(dim int, k int) *Graph {
	return &Graph{
		Dimension: dim,
		K:         k,
		Nodes:     make(map[int]*Node),
		LastID:    -1,
	}
}

func (g *Graph) Traverse(vec []float32) (int, map[int]float32, error) {
	if len(g.Nodes) == 0 {
		return -1, nil, fmt.Errorf("Cannot Traverse Empty Graph!")
	}
	// STEP 1: Greedy Navigation
	current := g.LastID
	visited := map[int]float32{}

	for {
		// get the node from the graph
		curr := g.Nodes[current]
		// initialize the best
		best := current
		bestDistance, err := EuclideanDistance(vec, curr.Vector)
		if err != nil {
			return 0, nil, err
		}
		visited[current] = bestDistance
		improved := false

		for _, nID := range curr.Neighbors {
			if d, seen := visited[nID]; seen {
				if d < bestDistance {
					bestDistance = d
					best = nID
					improved = true
				}
				continue
			}
			// compute the distance
			currDistance, err := EuclideanDistance(vec, g.Nodes[nID].Vector)
			if err != nil {
				return 0, nil, err
			}
			visited[nID] = currDistance

			if currDistance < bestDistance {
				bestDistance = currDistance
				best = nID
				improved = true
			}
		}

		if !improved {
			break
		}
		current = best
	}

	return current, visited, nil
}
