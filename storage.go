package vectorengine

type Graph struct {
	Vectors        []float32 // flat: N * Dimension
	Neighbors      []int     // flat: N * K (fixed blocks)
	NeighborCounts []int     // actual used neighbors per node

	K         int
	Dimension int
	Capacity  int
	Size      int
}

func NewGraphStore(dim int, k int, maxNodes int) *Graph {
	return &Graph{
		Vectors:        make([]float32, maxNodes*dim),
		Neighbors:      make([]int, maxNodes*k),
		NeighborCounts: make([]int, maxNodes),

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

func (g *Graph) GetNeighbors(id int) []int {
	start := id * g.K
	count := g.NeighborCounts[id]

	return g.Neighbors[start : start+count]
}

func (g *Graph) AddNeighbor(id int, neighbor int) {
	start := id * g.K
	count := g.NeighborCounts[id]

	if count >= g.K {
		return // or replace policy (important later)
	}

	g.Neighbors[start+count] = neighbor
	g.NeighborCounts[id]++
}
