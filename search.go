package vectorengine

import (
	"fmt"
)

type Result struct {
	ID       int
	Distance float32
}

func (g *Graph) Search(query []float32) (Result, error) {

	if len(query) != g.Dimension {
		return Result{}, fmt.Errorf("invalid query dimension")
	}

	if len(g.Nodes) == 0 {
		return Result{}, fmt.Errorf("store is empty")
	}

	nID, visited, err := g.Traverse(query)

	if err != nil {
		return Result{}, err
	}

	return Result{nID, visited[nID]}, nil
}
