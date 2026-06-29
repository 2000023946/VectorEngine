package vectorengine

import (
	"errors"
)

const EF = 64

// =========================================================
// CANDIDATE
// =========================================================

type Candidate struct {
	id   int
	dist float32
}

// =========================================================
// MAX HEAP (frontier only)
// worst (largest distance) at root
// =========================================================

type MaxHeap struct {
	data []Candidate
}

func (h *MaxHeap) less(i, j int) bool {
	return h.data[i].dist > h.data[j].dist
}

func (h *MaxHeap) Len() int {
	return len(h.data)
}

func (h *MaxHeap) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// -------------------- heapify --------------------

func (h *MaxHeap) heapifyUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !h.less(i, parent) {
			break
		}
		h.swap(i, parent)
		i = parent
	}
}

func (h *MaxHeap) heapifyDown(i int) {
	n := len(h.data)

	for {
		left := 2*i + 1
		right := 2*i + 2
		largest := i

		if left < n && h.less(left, largest) {
			largest = left
		}
		if right < n && h.less(right, largest) {
			largest = right
		}
		if largest == i {
			break
		}

		h.swap(i, largest)
		i = largest
	}
}

// -------------------- heap ops --------------------

func (h *MaxHeap) push(c Candidate) {
	h.data = append(h.data, c)
	h.heapifyUp(len(h.data) - 1)
}

func (h *MaxHeap) pop() Candidate {
	top := h.data[0]

	last := len(h.data) - 1
	h.data[0] = h.data[last]
	h.data = h.data[:last]

	if len(h.data) > 0 {
		h.heapifyDown(0)
	}

	return top
}

func (h *MaxHeap) top() Candidate {
	return h.data[0]
}

// =========================================================
// SEARCH (EF CANDIDATE POOL VERSION)
// =========================================================

func (g *Graph) GenerateCandidatePool(query []float32) ([]Candidate, error) {

	// -------------------- VALIDATE --------------------
	if g.Size == 0 {
		return nil, errors.New("empty graph")
	}

	if len(query) != g.Dimension {
		return nil, errors.New("dimension mismatch")
	}

	// =====================================================
	// PHASE 1: GREEDY ROUTING (find entry region)
	// =====================================================

	current := g.EntryPoint

	level := g.EntryPointLevel
	if level >= g.MaxLevels {
		level = g.MaxLevels - 1
	}

	for l := level; l >= 0; l-- {
		next, err := g.GreedyTraverseLayer(current, query, l)
		if err != nil {
			return nil, err
		}
		current = next
	}

	start := current

	// =====================================================
	// PHASE 2: EF SEARCH
	// =====================================================

	visited := make([]bool, g.Capacity+1)

	dist := func(a int) float32 {
		d, _ := EuclideanDistance(g.GetVector(a), query)
		return d
	}

	// frontier (heap)
	heap := &MaxHeap{
		data: make([]Candidate, 0, EF),
	}

	// result set (true EF output)
	results := make([]Candidate, 0, EF)

	// -------------------- init --------------------
	startCand := Candidate{start, dist(start)}

	heap.push(startCand)
	results = append(results, startCand)
	visited[start] = true

	// =====================================================
	// EXPLORE
	// =====================================================

	for heap.Len() > 0 {

		curr := heap.pop()

		// expand neighbors
		neighbors := g.GetNeighbors(curr.id, 0)

		for _, nb := range neighbors {

			if nb == 0 || visited[nb] {
				continue
			}

			visited[nb] = true

			d := dist(nb)
			c := Candidate{id: nb, dist: d}

			// -------------------- push to frontier --------------------
			if heap.Len() < EF {
				heap.push(c)
			} else if d < heap.top().dist {
				heap.pop()
				heap.push(c)
			}

			// -------------------- maintain EF result set --------------------
			results = append(results, c)

			if len(results) > EF {

				// remove worst from results (linear scan, OK for EF=64)
				worst := 0
				for i := 1; i < len(results); i++ {
					if results[i].dist > results[worst].dist {
						worst = i
					}
				}

				results[worst] = results[len(results)-1]
				results = results[:len(results)-1]
			}
		}
	}

	return results, nil
}
