package vectorengine

// import "fmt"

// func (g *Graph) Insert(vec []float32) error {

// 	// STEP 0: validate
// 	if len(vec) != g.Dimension {
// 		return fmt.Errorf("vector dimension mismatch")
// 	}

// 	newID := g.Size

// 	if newID >= g.Capacity {
// 		return fmt.Errorf("graph capacity exceeded")
// 	}

// 	// STEP 1: empty graph case
// 	if g.Size == 0 {
// 		g.SetVector(0, vec)
// 		g.Size = 1
// 		return nil
// 	}

// 	// STEP 2: traverse graph (find good region)
// 	current, visited, err := g.Traverse(vec)
// 	if err != nil {
// 		return err
// 	}

// 	// STEP 3: candidate selection (top-K)
// 	type cand struct {
// 		id   int
// 		dist float32
// 	}

// 	best := make([]cand, 0, g.K)

// 	findWorst := func() int {
// 		worst := 0
// 		for i := 1; i < len(best); i++ {
// 			if best[i].dist > best[worst].dist {
// 				worst = i
// 			}
// 		}
// 		return worst
// 	}

// 	process := func(id int) error {
// 		var d float32

// 		if dist, ok := visited[id]; ok {
// 			d = dist
// 		} else {
// 			var err error
// 			d, err = EuclideanDistance(vec, g.GetVector(id))
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		c := cand{id: id, dist: d}

// 		if len(best) < g.K {
// 			best = append(best, c)
// 			return nil
// 		}

// 		worst := findWorst()
// 		if c.dist < best[worst].dist {
// 			best[worst] = c
// 		}

// 		return nil
// 	}

// 	// STEP 4: evaluate current + neighbors
// 	if err := process(current); err != nil {
// 		return err
// 	}

// 	for _, nid := range g.GetNeighbors(current) {
// 		if err := process(nid); err != nil {
// 			return err
// 		}
// 	}

// 	// STEP 5: write vector into flat storage
// 	g.SetVector(newID, vec)
// 	g.Size++

// 	// STEP 6: connect graph (flat model)
// 	for i := 0; i < len(best); i++ {
// 		nid := best[i].id

// 		// add edge new -> neighbor
// 		g.AddNeighbor(newID, nid)

// 		// add reverse edge (bounded)
// 		g.AddNeighbor(nid, newID)
// 	}

// 	return nil
// }
