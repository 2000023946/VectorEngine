package vectorengine

func (g *Graph) GreedyTraverseLayer(nodeIndex int, query []float32, layer int) (int, error) {

	// start node
	current := nodeIndex

	currentVec := g.GetVector(current)
	bestDist, err := EuclideanDistance(currentVec, query)
	if err != nil {
		return -1, err
	}

	for {
		improved := false
		bestCandidate := current
		bestCandidateDist := bestDist

		// get neighbors of CURRENT NODE (NOT flat array index)
		neighbors := g.GetNeighbors(current, layer)

		for _, nb := range neighbors {
			if nb == 0 {
				continue
			}

			nbVec := g.GetVector(nb)

			dist, err := EuclideanDistance(nbVec, query)
			if err != nil {
				return -1, err
			}

			if dist < bestCandidateDist {
				bestCandidateDist = dist
				bestCandidate = nb
				improved = true
			}
		}

		// if no improvement → stop
		if !improved {
			return current, nil
		}

		// move to better node
		current = bestCandidate
		bestDist = bestCandidateDist
	}
}
