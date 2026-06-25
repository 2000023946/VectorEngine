package vectorengine

import (
	"math"
	"math/rand/v2"
)

type Graph struct {
	Vectors   []float32 // flat: N * Dimension
	Neighbors []int     // flat: N * K * Log(Capacity) (levels) (fixed blocks)
	Levels    []int     // actual used neighbors per node

	K         int
	Dimension int
	Capacity  int
	Size      int
}

func NewGraphStore(dim int, k int, maxNodes int) *Graph {
	levels := int(math.Log(float64(maxNodes))) + 1
	return &Graph{
		Vectors:   make([]float32, maxNodes*dim),
		Neighbors: make([]int, maxNodes*k*levels),
		Levels:    make([]int, levels),

		K:         k,
		Dimension: dim,
		Capacity:  maxNodes,
		Size:      0,
	}
}

func (g *Graph) GetVector(id int) []float32 {
	start := id * g.Dimension
	return g.Vectors[start : start+g.Dimension]
}

func (g *Graph) SetVector(id int, vec []float32) {
	start := id * g.Dimension
	copy(g.Vectors[start:start+g.Dimension], vec)
}

func (g *Graph) GetNeighbors(id int, layer int) []int {
	start := layer*g.Capacity*g.K + id*g.K

	return g.Neighbors[start : start+g.K]
}

func (g *Graph) AddNeighbor(id int, neighbor int, layer int) {
	start := layer*g.Capacity*g.K + id*g.K
	for i := 0; i < g.K; i++ {
		if g.Neighbors[start+i] == 0 {
			g.Neighbors[start+i] = neighbor
			break
		}
	}
}

func (g *Graph) GenerateRandomLayer() int {
	level := 0

	for rand.Float64() < 0.5 {
		level++
	}

	return level
}
