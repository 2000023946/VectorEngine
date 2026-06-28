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
	AccSize    = 50_000
	AccQueries = 1_000
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

	// ONLY track top-K (fast approach)
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

// -------------------- RECALL@K (SINGLE PRED VERSION) --------------------

// since your Search() returns ONLY 1 result,
// we measure: "is it in the true top-K?"
func recallAtKSingle(pred int, trueTopK []int) float64 {
	for _, v := range trueTopK {
		if v == pred {
			return 1.0
		}
	}
	return 0.0
}

// -------------------- TEST --------------------

func TestGraphRecallAtK(t *testing.T) {

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

		// ANN single result (YOUR CURRENT ENGINE)
		pred, err := g.Search(query)
		if err != nil {
			t.Fatalf("search failed: %v", err)
		}

		recall := recallAtKSingle(pred, trueTopK)

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

	t.Logf("📊 VECTORENGINE ACCURACY RESULTS (Search -> Top1 Eval)")
	t.Logf("Dataset size: %d", AccSize)
	t.Logf("K (ground truth): %d", AccK)
	t.Logf("Queries: %d", AccQueries)
	t.Logf("Avg Recall@%d (Top1 hit in TopK): %.4f", AccK, avgRecall)
	t.Logf("Min Recall: %.4f", minRecall)
	t.Logf("Max Recall: %.4f", maxRecall)

	// -------------------- ASSERTION --------------------

	if avgRecall < 0.2 {
		t.Fatalf("recall too low for single-shot search: %.4f", avgRecall)
	}
}
