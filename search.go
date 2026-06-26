package vectorengine

import (
	"errors"
)

func (g *Graph) Search(query []float32) (int, error) {

	// -------------------- VALIDATE GRAPH --------------------
	if g.Size == 0 {
		return -1, errors.New("empty graph")
	}

	// -------------------- VALIDATE QUERY --------------------
	if len(query) != g.Dimension {
		return -1, errors.New("query dimension mismatch")
	}

	// -------------------- VALIDATE ENTRY POINT --------------------
	if g.EntryPoint <= 0 || g.EntryPoint > g.Capacity {
		return -1, errors.New("invalid entry point")
	}

	// -------------------- START FROM ENTRY POINT --------------------
	current := g.EntryPoint

	// -------------------- START FROM ENTRY POINT LEVEL --------------------
	level := g.EntryPointLevel

	if level >= g.MaxLevels {
		level = g.MaxLevels - 1
	}

	// -------------------- TOP → BOTTOM DESCENT --------------------
	for l := level; l >= 0; l-- {

		next, err := g.GreedyTraverseLayer(current, query, l)
		if err != nil {
			return -1, err
		}

		// update current best node
		current = next
	}

	return current, nil
}
