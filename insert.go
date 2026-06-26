package vectorengine

import "errors"

// Insert adds a vector into the graph using HNSW-style construction.
func (g *Graph) Insert(vec []float32) (int, error) {

	// -------------------- VALIDATE INPUT --------------------
	if len(vec) != g.Dimension {
		return -1, errors.New("dimension mismatch")
	}

	// -------------------- ASSIGN NODE ID --------------------
	g.Size++
	newID := g.Size

	if newID <= 0 || newID > g.Capacity {
		return -1, errors.New("capacity exceeded")
	}

	// -------------------- STORE VECTOR --------------------
	g.SetVector(newID, vec)

	// -------------------- ASSIGN RANDOM LEVEL --------------------
	newLevel := g.GenerateRandomLayer()

	// -------------------- FIRST NODE CASE --------------------
	if g.Size == 1 {
		g.EntryPoint = newID
		g.EntryPointLevel = newLevel
		return newID, nil
	}

	// -------------------- FIND NEAREST NEIGHBOR (BOTTOM SEARCH) --------------------
	// NOTE: This assumes you have a vector-based search function.
	// If Search() is query-based, you should replace it with SearchVector().
	current, err := g.Search(vec)
	if err != nil {
		return -1, err
	}

	// -------------------- CONNECT NODES ACROSS LEVELS --------------------
	maxL := newLevel
	if maxL >= g.MaxLevels {
		maxL = g.MaxLevels - 1
	}

	for l := 0; l <= maxL; l++ {

		// connect new node -> existing node
		g.AddNeighbor(newID, current, l)

		// connect existing node -> new node (bidirectional graph)
		g.AddNeighbor(current, newID, l)
	}

	// -------------------- UPDATE ENTRY POINT IF NEEDED --------------------
	if newLevel > g.EntryPointLevel {
		g.EntryPoint = newID
		g.EntryPointLevel = newLevel
	}

	return newID, nil
}
