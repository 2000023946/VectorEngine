package src

import (
	"fmt"
	"math"
	"sort"
)

const (
	// IVF is built once the engine reaches this many vectors.
	IVFThreshold = 10_000

	// Number of IVF clusters.
	NumCentroids = 100

	// Number of centroids searched for each query.
	NumProbes = 60

	// Number of Lloyd's iterations used to train the centroids.
	LloydIterations = 5
)

type Vector struct {
	ID     int
	Values []float64
}

type Result struct {
	ID       int
	Distance float64
}

type VectorEngine struct {
	vectors   []Vector
	dimension int

	// IVF index.
	centroids [][]float64
	buckets   [][]int
	indexed   bool
}

func NewVectorEngine() *VectorEngine {
	return &VectorEngine{
		vectors: make([]Vector, 0),
	}
}

func (ve *VectorEngine) Insert(id int, values []float64) error {
	// The first vector establishes the dimension
	// of the vector space.
	if ve.dimension == 0 {
		ve.dimension = len(values)
	} else if len(values) != ve.dimension {
		return fmt.Errorf(
			"dimension mismatch: expected %d, got %d",
			ve.dimension,
			len(values),
		)
	}

	vector := Vector{
		ID:     id,
		Values: values,
	}

	ve.vectors = append(ve.vectors, vector)

	// Build the IVF index once the threshold is reached.
	if !ve.indexed && len(ve.vectors) >= IVFThreshold {
		ve.buildIndex()
		return nil
	}

	// If the index already exists, place the new vector
	// into its nearest IVF bucket.
	if ve.indexed {
		centroid := ve.nearestCentroid(values)
		ve.buckets[centroid] = append(
			ve.buckets[centroid],
			len(ve.vectors)-1,
		)
	}

	return nil
}

func distance(a []float64, b []float64) float64 {
	var sum float64

	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

// --------------------------------------------------
// IVF INDEX
// --------------------------------------------------

func (ve *VectorEngine) buildIndex() {
	// Do not attempt to create more centroids than vectors.
	numCentroids := NumCentroids

	if len(ve.vectors) < numCentroids {
		numCentroids = len(ve.vectors)
	}

	ve.centroids = make([][]float64, numCentroids)

	// Simple deterministic initialization.
	// We choose evenly spaced vectors as the initial centroids.
	for i := 0; i < numCentroids; i++ {
		index := i * len(ve.vectors) / numCentroids

		centroid := make([]float64, ve.dimension)
		copy(centroid, ve.vectors[index].Values)

		ve.centroids[i] = centroid
	}

	// Lloyd's algorithm.
	for iteration := 0; iteration < LloydIterations; iteration++ {
		sums := make([][]float64, numCentroids)
		counts := make([]int, numCentroids)

		for i := 0; i < numCentroids; i++ {
			sums[i] = make([]float64, ve.dimension)
		}

		// Assign every vector to its nearest centroid.
		for _, vector := range ve.vectors {
			centroid := ve.nearestCentroid(vector.Values)

			for d := 0; d < ve.dimension; d++ {
				sums[centroid][d] += vector.Values[d]
			}

			counts[centroid]++
		}

		// Recalculate centroid positions.
		for i := 0; i < numCentroids; i++ {
			if counts[i] == 0 {
				continue
			}

			for d := 0; d < ve.dimension; d++ {
				ve.centroids[i][d] =
					sums[i][d] / float64(counts[i])
			}
		}
	}

	// Create inverted buckets.
	ve.buckets = make([][]int, numCentroids)

	for index, vector := range ve.vectors {
		centroid := ve.nearestCentroid(vector.Values)

		ve.buckets[centroid] = append(
			ve.buckets[centroid],
			index,
		)
	}

	ve.indexed = true
}

func (ve *VectorEngine) nearestCentroid(vector []float64) int {
	best := 0
	bestDistance := math.Inf(1)

	for i, centroid := range ve.centroids {
		d := distance(vector, centroid)

		if d < bestDistance {
			bestDistance = d
			best = i
		}
	}

	return best
}

// --------------------------------------------------
// SEARCH
// --------------------------------------------------

func (ve *VectorEngine) Search(query []float64, k int) []Result {
	if len(query) != ve.dimension {
		return nil
	}

	if k <= 0 {
		return []Result{}
	}

	// Before the IVF threshold, use exact brute force.
	if !ve.indexed {
		return ve.bruteForceSearch(query, k)
	}

	// Find distances to every centroid.
	type centroidDistance struct {
		index    int
		distance float64
	}

	centroidDistances := make(
		[]centroidDistance,
		0,
		len(ve.centroids),
	)

	for i, centroid := range ve.centroids {
		centroidDistances = append(
			centroidDistances,
			centroidDistance{
				index:    i,
				distance: distance(query, centroid),
			},
		)
	}

	// Closest centroids first.
	sort.Slice(
		centroidDistances,
		func(i, j int) bool {
			return centroidDistances[i].distance <
				centroidDistances[j].distance
		},
	)

	numProbes := NumProbes

	if numProbes > len(centroidDistances) {
		numProbes = len(centroidDistances)
	}

	// Search only vectors inside the closest buckets.
	results := make([]Result, 0)

	for i := 0; i < numProbes; i++ {
		bucket := ve.buckets[centroidDistances[i].index]

		for _, vectorIndex := range bucket {
			vector := ve.vectors[vectorIndex]

			d := distance(query, vector.Values)

			results = append(results, Result{
				ID:       vector.ID,
				Distance: d,
			})
		}
	}

	// Sort candidates by actual vector distance.
	sort.Slice(
		results,
		func(i, j int) bool {
			return results[i].Distance < results[j].Distance
		},
	)

	if k > len(results) {
		k = len(results)
	}

	return results[:k]
}

// --------------------------------------------------
// BRUTE FORCE FALLBACK
// --------------------------------------------------

func (ve *VectorEngine) bruteForceSearch(
	query []float64,
	k int,
) []Result {

	results := make(
		[]Result,
		0,
		len(ve.vectors),
	)

	for _, vector := range ve.vectors {
		d := distance(query, vector.Values)

		results = append(results, Result{
			ID:       vector.ID,
			Distance: d,
		})
	}

	sort.Slice(
		results,
		func(i, j int) bool {
			return results[i].Distance <
				results[j].Distance
		},
	)

	if k > len(results) {
		k = len(results)
	}

	return results[:k]
}
