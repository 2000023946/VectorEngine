package vectorengine

import (
	"fmt"
	"sort"
)

func (g *Graph) Insert(vec []float32) error {
	if g.Dimension != len(vec) {
		return fmt.Errorf("Vector to insert dimension does not match Graph dimension")
	}
	newID := g.LastID + 1
	// STEP 1: Create Node
	newNode := &Node{
		ID:        newID,
		Vector:    vec,
		Neighbors: []int{},
	}

	if len(g.Nodes) == 0 {
		g.Nodes[newID] = newNode
		g.LastID = newID
	}

	current, visited, err := g.Traverse(vec)
	if err != nil {
		return err
	}

	candidateSet := map[int]bool{}
	candidateSet[current] = true

	for _, nID := range g.Nodes[current].Neighbors {
		candidateSet[nID] = true
	}

	type cand struct {
		id   int
		dist float32
	}

	var list []cand

	for id := range candidateSet {
		if d, seen := visited[id]; seen {
			list = append(list, cand{id, d})
		} else {
			d, err := EuclideanDistance(vec, g.Nodes[id].Vector)
			if err != nil {
				return err
			}
			list = append(list, cand{id, d})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].dist < list[j].dist
	})

	for i := 0; i < g.K && i < len(list); i++ {
		nid := list[i].id
		newNode.Neighbors = append(newNode.Neighbors, nid)

		// connect nid -> new node
		if len(g.Nodes[nid].Neighbors) < g.K {
			g.Nodes[nid].Neighbors = append(g.Nodes[nid].Neighbors, newNode.ID)
		}
	}

	g.Nodes[newID] = newNode
	g.LastID = newID
	return nil
}
