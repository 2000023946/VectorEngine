package vectorengine

import (
	"errors"
	"sort"
)

// =========================================================
// SEARCH TOP-K USING EF CANDIDATE POOL
// =========================================================

func (g *Graph) SearchTopK(query []float32, k int) ([]Candidate, error) {

	// -------------------- VALIDATE --------------------
	if g.Size == 0 {
		return nil, errors.New("empty graph")
	}

	if len(query) != g.Dimension {
		return nil, errors.New("dimension mismatch")
	}

	if k <= 0 {
		return nil, errors.New("invalid k")
	}

	// =====================================================
	// STEP 1: GET EF CANDIDATE POOL
	// =====================================================

	pool, err := g.GenerateCandidatePool(query)
	if err != nil {
		return nil, err
	}

	if len(pool) == 0 {
		return nil, errors.New("empty candidate pool")
	}

	// =====================================================
	// STEP 2: SORT BY DISTANCE (best first)
	// =====================================================

	sort.Slice(pool, func(i, j int) bool {
		return pool[i].dist < pool[j].dist
	})

	// =====================================================
	// STEP 3: RETURN TOP-K
	// =====================================================

	if k > len(pool) {
		k = len(pool)
	}

	return pool[:k], nil
}
