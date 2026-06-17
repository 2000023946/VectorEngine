package vectorengine

import (
	"fmt"
	"math"
)

type Result struct {
	Index    int
	Distance float32
}

func (vs *VectorStore) Search(query []float32) (Result, error) {

	if len(query) != vs.dimension {
		return Result{}, fmt.Errorf("invalid query dimension")
	}

	if len(vs.vectors) == 0 {
		return Result{}, fmt.Errorf("store is empty")
	}

	best := Result{
		Index:    -1,
		Distance: math.MaxFloat32,
	}

	for i, vec := range vs.vectors {
		dist, err := EuclideanDistance(query, vec)
		if err != nil {
			return Result{}, err
		}

		if dist < best.Distance {
			best.Index = i
			best.Distance = dist
		}
	}

	return best, nil
}
