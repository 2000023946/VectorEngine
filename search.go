package vectorengine

import (
	"fmt"
)

type Result struct {
	ID       int
	Distance float32
}

func (g *Graph) Search(query []float32) (Result, error) {

	// STEP 1: validate dimension
	if len(query) != g.Dimension {
		return Result{}, fmt.Errorf("invalid query dimension")
	}

	// STEP 2: empty graph check (use Size, not len(Vectors))
	if g.Size == 0 {
		return Result{}, fmt.Errorf("store is empty")
	}

	// STEP 3: traverse graph (find best candidate)
	bestID, _, err := g.Traverse(query)
	if err != nil {
		return Result{}, err
	}

	// STEP 4: compute final distance ONLY once (clean design)
	bestVec := g.GetVector(bestID)

	dist, err := EuclideanDistance(query, bestVec)
	if err != nil {
		return Result{}, err
	}

	return Result{
		ID:       bestID,
		Distance: dist,
	}, nil
}
