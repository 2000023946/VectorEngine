package vectorengine

import "fmt"

func (g *Graph) Traverse(vec []float32) (int, map[int]float32, error) {
	if g.Size == 0 {
		return -1, nil, fmt.Errorf("cannot traverse empty graph")
	}

	current := 0
	visited := make(map[int]float32)

	for {
		currVec := g.GetVector(current)

		best := current
		bestDistance, err := EuclideanDistance(vec, currVec)
		if err != nil {
			return -1, nil, err
		}

		visited[current] = bestDistance
		improved := false

		// ✅ FIXED: use flat neighbors
		neighbors := g.GetNeighbors(current)

		for _, nID := range neighbors {

			// skip if already visited
			if d, seen := visited[nID]; seen {
				if d < bestDistance {
					bestDistance = d
					best = nID
					improved = true
				}
				continue
			}

			nVec := g.GetVector(nID)

			currDistance, err := EuclideanDistance(vec, nVec)
			if err != nil {
				return -1, nil, err
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
