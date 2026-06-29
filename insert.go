package vectorengine

import (
	"errors"
	"sort"
)

// =========================================================
// INSERT USING EF CANDIDATE POOL
// GRAPH DEGREE = g.K (no separate M)
// =========================================================

func (g *Graph) Insert(vec []float32) (int, error) {

	// -------------------- VALIDATE --------------------
	if len(vec) != g.Dimension {
		return -1, errors.New("dimension mismatch")
	}

	if g.K <= 0 {
		return -1, errors.New("invalid graph degree K")
	}

	// -------------------- ASSIGN NODE ID --------------------
	g.Size++
	newID := g.Size

	if newID > g.Capacity {
		return -1, errors.New("capacity exceeded")
	}

	// -------------------- STORE VECTOR --------------------
	g.SetVector(newID, vec)

	// -------------------- FIRST NODE --------------------
	if g.Size == 1 {
		g.EntryPoint = newID
		g.EntryPointLevel = 0
		return newID, nil
	}

	// =====================================================
	// STEP 1: GET EF CANDIDATE POOL
	// =====================================================

	pool, err := g.GenerateCandidatePool(vec)
	if err != nil {
		return -1, err
	}

	if len(pool) == 0 {
		return -1, errors.New("empty candidate pool")
	}

	// =====================================================
	// STEP 2: SORT BY DISTANCE (BEST FIRST)
	// =====================================================

	sort.Slice(pool, func(i, j int) bool {
		return pool[i].dist < pool[j].dist
	})

	// =====================================================
	// STEP 3: SELECT TOP K NEIGHBORS (GRAPH DEGREE)
	// =====================================================

	k := g.K
	if k > len(pool) {
		k = len(pool)
	}

	// =====================================================
	// STEP 4: CONNECT GRAPH (BIDIRECTIONAL)
	// =====================================================

	for i := 0; i < k; i++ {
		nb := pool[i]

		// new → neighbor
		g.AddNeighbor(newID, nb.id, 0)

		// neighbor → new
		g.AddNeighbor(nb.id, newID, 0)
	}

	// =====================================================
	// STEP 5: UPDATE ENTRY POINT (simple heuristic)
	// =====================================================

	if pool[0].id != 0 {
		g.EntryPoint = pool[0].id
	}

	return newID, nil
}
