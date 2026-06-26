package vectorengine

import (
	"errors"
	"math"
	"math/rand/v2"
)

type Graph struct {
	Vectors   []float32 // flat: N * Dimension
	Neighbors []int     // flat: N * K * MaxLevels

	K         int
	Dimension int
	Capacity  int
	Size      int
	MaxLevels int

	EntryPoint      int // starting node
	EntryPointLevel int // highest level of entry node
}

// -------------------- INIT --------------------

func NewGraphStore(dim int, k int, maxNodes int) *Graph {

	maxLevels := int(math.Log(float64(maxNodes))) + 1

	return &Graph{
		Vectors:   make([]float32, (maxNodes+1)*dim),
		Neighbors: make([]int, (maxNodes+1)*k*maxLevels),

		K:         k,
		Dimension: dim,
		Capacity:  maxNodes,
		Size:      0,
		MaxLevels: maxLevels,

		EntryPoint:      1,
		EntryPointLevel: 0,
	}
}

// -------------------- VECTOR OPS --------------------

func (g *Graph) GetVector(id int) []float32 {
	start := id * g.Dimension
	return g.Vectors[start : start+g.Dimension]
}

func (g *Graph) SetVector(id int, vec []float32) {
	start := id * g.Dimension
	copy(g.Vectors[start:start+g.Dimension], vec)
}

// -------------------- INDEXING CORE --------------------

func (g *Graph) getIndex(nodeIndex int, layer int) int {
	return layer*g.Capacity*g.K + nodeIndex*g.K
}

func (g *Graph) getIndexSafe(nodeIndex int, layer int) (int, error) {

	if nodeIndex <= 0 || nodeIndex > g.Capacity {
		return -1, errors.New("node index out of range")
	}

	if layer < 0 || layer >= g.MaxLevels {
		return -1, errors.New("layer out of range")
	}

	return g.getIndex(nodeIndex, layer), nil
}

// -------------------- SAFE NEIGHBOR ACCESS --------------------

func (g *Graph) GetNeighborValue(nodeIndex int, layer int, offset int) (int, error) {

	if nodeIndex <= 0 || nodeIndex > g.Capacity {
		return -1, errors.New("node index out of range")
	}

	if layer < 0 || layer >= g.MaxLevels {
		return -1, errors.New("layer out of range")
	}

	if offset < 0 || offset >= g.K {
		return -1, errors.New("offset out of range")
	}

	start := g.getIndex(nodeIndex, layer)

	if start+offset < 0 || start+offset >= len(g.Neighbors) {
		return -1, errors.New("computed index out of bounds")
	}

	val := g.Neighbors[start+offset]

	if val == 0 {
		return -1, errors.New("no neighbor exists at this position")
	}

	return val, nil
}

// -------------------- NEIGHBORS --------------------

func (g *Graph) GetNeighbors(id int, layer int) []int {
	start := g.getIndex(id, layer)
	return g.Neighbors[start : start+g.K]
}

func (g *Graph) AddNeighbor(id int, neighbor int, layer int) {
	start := g.getIndex(id, layer)

	for i := 0; i < g.K; i++ {
		if g.Neighbors[start+i] == 0 {
			g.Neighbors[start+i] = neighbor
			return
		}
	}
}

// -------------------- LAYER GENERATION --------------------

func (g *Graph) GenerateRandomLayer() int {
	level := 0

	for rand.Float64() < 0.5 {
		level++
	}

	if level >= g.MaxLevels {
		level = g.MaxLevels - 1
	}

	return level
}
