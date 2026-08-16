package src

import (
	"fmt"
	"math"
	"sort"
)

const (
	// IVF configuration.
	indexBuildThreshold = 10_000
	numClusters         = 100
	numProbes           = 5
	maxKMeansIterations = 20
)

type Vector struct {
	ID     int
	Values []float64
}

type Result struct {
	ID       int
	Distance float64
}

type IVFIndex struct {
	centroids [][]float64
	buckets   [][]int
}

type VectorEngine struct {
	vectors   []Vector
	dimension int

	// IVF index is built once the engine reaches
	// indexBuildThreshold vectors.
	index *IVFIndex
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

	// Build the IVF index once the dataset reaches
	// the configured threshold.
	if ve.index == nil && len(ve.vectors) >= indexBuildThreshold {
		ve.buildIVF()
		return nil
	}

	// Once the index exists, add new vectors to the
	// bucket belonging to their nearest centroid.
	if ve.index != nil {
		ve.addToIVF(len(ve.vectors) - 1)
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

func (ve *VectorEngine) buildIVF() {
	clusterCount := numClusters

	if len(ve.vectors) < clusterCount {
		clusterCount = len(ve.vectors)
	}

	centroids := make([][]float64, clusterCount)

	// Deterministic initialization.
	//
	// We spread the initial centroids throughout
	// the dataset instead of choosing them randomly.
	for i := 0; i < clusterCount; i++ {
		index := i * len(ve.vectors) / clusterCount

		centroids[i] = copyVector(
			ve.vectors[index].Values,
		)
	}

	assignments := make([]int, len(ve.vectors))

	for iteration := 0; iteration < maxKMeansIterations; iteration++ {
		changed := false

		// ------------------------------------------
		// Assignment step
		// ------------------------------------------

		for i, vector := range ve.vectors {
			bestCluster := 0
			bestDistance := distance(
				vector.Values,
				centroids[0],
			)

			for c := 1; c < clusterCount; c++ {
				d := distance(
					vector.Values,
					centroids[c],
				)

				if d < bestDistance {
					bestDistance = d
					bestCluster = c
				}
			}

			if assignments[i] != bestCluster || iteration == 0 {
				changed = true
				assignments[i] = bestCluster
			}
		}

		// ------------------------------------------
		// Update step
		// ------------------------------------------

		sums := make([][]float64, clusterCount)
		counts := make([]int, clusterCount)

		for c := 0; c < clusterCount; c++ {
			sums[c] = make([]float64, ve.dimension)
		}

		for i, vector := range ve.vectors {
			cluster := assignments[i]

			counts[cluster]++

			for d := 0; d < ve.dimension; d++ {
				sums[cluster][d] += vector.Values[d]
			}
		}

		for c := 0; c < clusterCount; c++ {
			// Keep the previous centroid if the cluster
			// received no vectors.
			if counts[c] == 0 {
				continue
			}

			for d := 0; d < ve.dimension; d++ {
				centroids[c][d] =
					sums[c][d] / float64(counts[c])
			}
		}

		if !changed {
			break
		}
	}

	// ----------------------------------------------
	// Build inverted buckets
	// ----------------------------------------------

	buckets := make([][]int, clusterCount)

	for vectorIndex, cluster := range assignments {
		buckets[cluster] = append(
			buckets[cluster],
			vectorIndex,
		)
	}

	ve.index = &IVFIndex{
		centroids: centroids,
		buckets:   buckets,
	}
}

func (ve *VectorEngine) addToIVF(vectorIndex int) {
	vector := ve.vectors[vectorIndex]

	cluster := ve.nearestCentroid(vector.Values)

	ve.index.buckets[cluster] = append(
		ve.index.buckets[cluster],
		vectorIndex,
	)
}

func (ve *VectorEngine) nearestCentroid(vector []float64) int {
	bestCluster := 0
	bestDistance := distance(
		vector,
		ve.index.centroids[0],
	)

	for c := 1; c < len(ve.index.centroids); c++ {
		d := distance(
			vector,
			ve.index.centroids[c],
		)

		if d < bestDistance {
			bestDistance = d
			bestCluster = c
		}
	}

	return bestCluster
}

func copyVector(vector []float64) []float64 {
	result := make([]float64, len(vector))
	copy(result, vector)
	return result
}

// --------------------------------------------------
// SEARCH
// --------------------------------------------------

func (ve *VectorEngine) Search(query []float64, k int) []Result {
	if len(query) != ve.dimension {
		return nil
	}

	if k < 0 {
		k = 0
	}

	if k == 0 {
		return []Result{}
	}

	// Before the IVF threshold, use the original
	// brute-force implementation.
	if ve.index == nil {
		return ve.bruteForceSearch(query, k)
	}

	// ----------------------------------------------
	// IVF SEARCH
	// ----------------------------------------------

	// Find the nearest clusters to the query.
	type clusterDistance struct {
		cluster  int
		distance float64
	}

	clusterDistances := make(
		[]clusterDistance,
		len(ve.index.centroids),
	)

	for i, centroid := range ve.index.centroids {
		clusterDistances[i] = clusterDistance{
			cluster:  i,
			distance: distance(query, centroid),
		}
	}

	sort.Slice(
		clusterDistances,
		func(i, j int) bool {
			return clusterDistances[i].distance <
				clusterDistances[j].distance
		},
	)

	probes := numProbes

	if probes > len(clusterDistances) {
		probes = len(clusterDistances)
	}

	// Search only vectors contained in the closest
	// clusters.
	results := make([]Result, 0)

	for i := 0; i < probes; i++ {
		cluster := clusterDistances[i].cluster

		for _, vectorIndex := range ve.index.buckets[cluster] {
			vector := ve.vectors[vectorIndex]

			d := distance(
				query,
				vector.Values,
			)

			results = append(results, Result{
				ID:       vector.ID,
				Distance: d,
			})
		}
	}

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
// BRUTE FORCE
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

	// Compare the query against EVERY vector.
	for _, vector := range ve.vectors {
		d := distance(
			query,
			vector.Values,
		)

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
