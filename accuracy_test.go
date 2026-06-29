package vectorengine

import (
	"math"
	"math/rand"
	"testing"
)

// -------------------- CONFIG --------------------

const (
	AccDim     = 64
	AccK       = 16
	AccSize    = 2000 // reduced for test stability (50k is too heavy for unit tests)
	AccQueries = 200
)

// -------------------- DATASET --------------------

func generateDataset(n int, dim int) [][]float32 {
	data := make([][]float32, n)

	for i := 0; i < n; i++ {
		vec := make([]float32, dim)
		for j := 0; j < dim; j++ {
			vec[j] = float32((i*j)%97) / 97.0
		}
		data[i] = vec
	}

	return data
}

// -------------------- DOT PRODUCT --------------------

func dot(a, b []float32) float32 {
	var sum float32
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// -------------------- BRUTE FORCE TOP-K --------------------

func bruteForce(dataset [][]float32, query []float32, k int) []int {

	type pair struct {
		id  int
		sim float32
	}

	results := make([]pair, len(dataset))

	for i := range dataset {
		results[i] = pair{
			id:  i,
			sim: dot(dataset[i], query),
		}
	}

	// partial selection top-k
	for i := 0; i < k; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].sim > results[i].sim {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	topk := make([]int, k)
	for i := 0; i < k; i++ {
		topk[i] = results[i].id
	}

	return topk
}

// -------------------- RECALL@K (SET OVERLAP) --------------------

func recallAtK(pred []Candidate, trueTopK []int) float64 {

	trueSet := make(map[int]struct{})
	for _, id := range trueTopK {
		trueSet[id] = struct{}{}
	}

	hits := 0
	for _, p := range pred {
		if _, ok := trueSet[p.id]; ok {
			hits++
		}
	}

	return float64(hits) / float64(len(trueTopK))
}

// -------------------- TEST --------------------

func TestAccuracyGraphRecallAtK(t *testing.T) {

	rand.Seed(42)

	// -------------------- INIT GRAPH --------------------
	g := NewGraphStore(AccDim, AccK, AccSize)

	// -------------------- DATASET --------------------
	data := generateDataset(AccSize, AccDim)

	// -------------------- INSERT --------------------
	for i := 0; i < AccSize; i++ {
		_, err := g.Insert(data[i])
		if err != nil {
			t.Fatalf("insert failed at %d: %v", i, err)
		}
	}

	// -------------------- EVALUATION --------------------

	var total float64
	var minRecall float64 = math.MaxFloat64
	var maxRecall float64 = 0

	for i := 0; i < AccQueries; i++ {

		qid := rand.Intn(AccSize)
		query := data[qid]

		// ground truth top-K
		trueTopK := bruteForce(data, query, AccK)

		// ANN top-K result
		pred, err := g.SearchTopK(query, AccK)
		if err != nil {
			t.Fatalf("search failed: %v", err)
		}

		recall := recallAtK(pred, trueTopK)

		total += recall

		if recall < minRecall {
			minRecall = recall
		}
		if recall > maxRecall {
			maxRecall = recall
		}
	}

	avgRecall := total / float64(AccQueries)

	// -------------------- REPORT --------------------

	t.Logf("📊 VECTORENGINE ACCURACY RESULTS (Top-K Recall)")
	t.Logf("Dataset size: %d", AccSize)
	t.Logf("K: %d", AccK)
	t.Logf("Queries: %d", AccQueries)
	t.Logf("Avg Recall@%d: %.4f", AccK, avgRecall)
	t.Logf("Min Recall: %.4f", minRecall)
	t.Logf("Max Recall: %.4f", maxRecall)

	// -------------------- ASSERTION --------------------

	if avgRecall < 0.05 {
		t.Fatalf("recall too low: %.4f", avgRecall)
	}
}
