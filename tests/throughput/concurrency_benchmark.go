package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"

	"vectorengine/src"
)

const (
	dimension = 128

	// Initial dataset used by search and mixed workloads.
	initialVectors = 1_000_000

	// Number of operations measured at each concurrency level.
	totalOperations = 1000

	// Number of nearest neighbors requested.
	k = 10

	// Mixed workload:
	// 80% search
	// 20% insert
	searchPercentage = 80
	insertPercentage = 20
)

type OperationType int

const (
	SearchOperation OperationType = iota
	InsertOperation
)

type Operation struct {
	Type  OperationType
	Query []float64
	ID    int
}

type BenchmarkResult struct {
	Workload    string
	Concurrency int
	Operations  int

	TotalTime time.Duration

	Throughput float64

	Average time.Duration
	P50     time.Duration
	P95     time.Duration
	P99     time.Duration
}

// --------------------------------------------------
// Generate one vector.
// --------------------------------------------------

func generateVector() []float64 {
	vector := make([]float64, dimension)

	for i := range vector {
		vector[i] = rand.Float64()
	}

	return vector
}

// --------------------------------------------------
// Percentile calculation.
//
// Input must already be sorted.
// --------------------------------------------------

func percentile(
	sorted []time.Duration,
	p float64,
) time.Duration {

	if len(sorted) == 0 {
		return 0
	}

	index := int(
		float64(len(sorted)-1) * p,
	)

	return sorted[index]
}

// --------------------------------------------------
// Build a VectorEngine containing initialVectors.
//
// This setup happens OUTSIDE the timed benchmark.
// --------------------------------------------------

func buildEngine() *src.VectorEngine {
	fmt.Printf(
		"Building engine with %d vectors...\n",
		initialVectors,
	)

	engine := src.NewVectorEngine()

	for i := 0; i < initialVectors; i++ {
		vector := generateVector()

		if err := engine.Insert(i, vector); err != nil {
			panic(err)
		}
	}

	fmt.Println("Engine ready.")

	return engine
}

// --------------------------------------------------
// Generate the fixed workload before timing.
//
// Search-only:
//
// 100% Search
//
// Insert-only:
//
// 100% Insert
//
// Mixed:
//
// 80% Search
// 20% Insert
//
// The operations are shuffled so the mixed workload
// doesn't simply do all searches followed by inserts.
// --------------------------------------------------

func generateOperations(
	workload string,
	startID int,
) []Operation {

	operations := make(
		[]Operation,
		totalOperations,
	)

	switch workload {

	case "search":

		for i := 0; i < totalOperations; i++ {
			operations[i] = Operation{
				Type:  SearchOperation,
				Query: generateVector(),
			}
		}

	case "insert":

		for i := 0; i < totalOperations; i++ {
			operations[i] = Operation{
				Type:  InsertOperation,
				ID:    startID + i,
				Query: generateVector(),
			}
		}

	case "mixed":

		searchCount :=
			totalOperations *
				searchPercentage / 100

		for i := 0; i < totalOperations; i++ {

			if i < searchCount {

				operations[i] = Operation{
					Type:  SearchOperation,
					Query: generateVector(),
				}

			} else {

				operations[i] = Operation{
					Type:  InsertOperation,
					ID:    startID + i,
					Query: generateVector(),
				}
			}
		}

		// Shuffle the workload while preserving
		// the exact 80/20 ratio.
		rand.Shuffle(
			len(operations),
			func(i, j int) {
				operations[i], operations[j] =
					operations[j], operations[i]
			},
		)

	default:
		panic("unknown workload")
	}

	return operations
}

// --------------------------------------------------
// Execute one operation.
// --------------------------------------------------

func executeOperation(
	engine *src.VectorEngine,
	operation Operation,
) {

	switch operation.Type {

	case SearchOperation:

		engine.Search(
			operation.Query,
			k,
		)

	case InsertOperation:

		if err := engine.Insert(
			operation.ID,
			operation.Query,
		); err != nil {
			panic(err)
		}
	}
}

// --------------------------------------------------
// Run benchmark at one concurrency level.
// --------------------------------------------------

func runBenchmark(
	engine *src.VectorEngine,
	workload string,
	concurrency int,
	startID int,
) BenchmarkResult {

	// --------------------------------------------------
	// Generate workload BEFORE measurement.
	// --------------------------------------------------

	operations :=
		generateOperations(
			workload,
			startID,
		)

	latencies :=
		make([]time.Duration, len(operations))

	var wg sync.WaitGroup

	// --------------------------------------------------
	// Split operations across goroutines.
	// --------------------------------------------------

	operationsPerWorker :=
		(len(operations) + concurrency - 1) /
			concurrency

	startTime := time.Now()

	for worker := 0; worker < concurrency; worker++ {

		start :=
			worker * operationsPerWorker

		end :=
			start + operationsPerWorker

		if end > len(operations) {
			end = len(operations)
		}

		if start >= len(operations) {
			continue
		}

		wg.Add(1)

		go func(start, end int) {

			defer wg.Done()

			for i := start; i < end; i++ {

				requestStart :=
					time.Now()

				executeOperation(
					engine,
					operations[i],
				)

				latencies[i] =
					time.Since(requestStart)
			}

		}(start, end)
	}

	wg.Wait()

	totalTime :=
		time.Since(startTime)

	// --------------------------------------------------
	// Sort latency samples.
	// --------------------------------------------------

	sortedLatencies :=
		make([]time.Duration, len(latencies))

	copy(
		sortedLatencies,
		latencies,
	)

	sort.Slice(
		sortedLatencies,
		func(i, j int) bool {
			return sortedLatencies[i] <
				sortedLatencies[j]
		},
	)

	// --------------------------------------------------
	// Average latency.
	// --------------------------------------------------

	var totalLatency time.Duration

	for _, latency := range latencies {
		totalLatency += latency
	}

	average :=
		totalLatency /
			time.Duration(len(latencies))

	// --------------------------------------------------
	// Throughput.
	//
	// Total completed operations divided by elapsed
	// seconds.
	// --------------------------------------------------

	throughput :=
		float64(len(operations)) /
			totalTime.Seconds()

	return BenchmarkResult{
		Workload:    workload,
		Concurrency: concurrency,
		Operations:  len(operations),

		TotalTime: totalTime,

		Throughput: throughput,

		Average: average,

		P50: percentile(
			sortedLatencies,
			0.50,
		),

		P95: percentile(
			sortedLatencies,
			0.95,
		),

		P99: percentile(
			sortedLatencies,
			0.99,
		),
	}
}

// --------------------------------------------------
// Print table header.
// --------------------------------------------------

func printHeader() {

	fmt.Printf(
		"%10s | %10s | %12s | %10s | %10s | %10s | %10s\n",
		"Concurrency",
		"Operations",
		"Ops/sec",
		"Average",
		"P50",
		"P95",
		"P99",
	)

	fmt.Println(
		"------------+------------+--------------+------------+------------+------------+------------",
	)
}

// --------------------------------------------------
// Print benchmark result.
// --------------------------------------------------

func printResult(
	result BenchmarkResult,
) {

	fmt.Printf(
		"%10d | %10d | %12.2f | %10s | %10s | %10s | %10s\n",
		result.Concurrency,
		result.Operations,
		result.Throughput,
		result.Average,
		result.P50,
		result.P95,
		result.P99,
	)
}

// --------------------------------------------------
// Warmup searches.
//
// Prevent runtime startup effects from contaminating
// the benchmark.
// --------------------------------------------------

func warmup(
	engine *src.VectorEngine,
) {

	fmt.Println("Warming up...")

	for i := 0; i < 50; i++ {
		engine.Search(
			generateVector(),
			k,
		)
	}

	fmt.Println("Warmup complete.")
}

// --------------------------------------------------
// Run a workload across all concurrency levels.
// --------------------------------------------------

func runWorkload(
	workload string,
	concurrencyLevels []int,
) {

	fmt.Println()
	fmt.Println(
		"============================================================",
	)

	switch workload {

	case "search":

		fmt.Println("SEARCH-ONLY")
		fmt.Println("100% Search")

	case "insert":

		fmt.Println("INSERT-ONLY")
		fmt.Println("100% Insert")

	case "mixed":

		fmt.Println("MIXED WORKLOAD")

		fmt.Printf(
			"%d%% Search / %d%% Insert\n",
			searchPercentage,
			insertPercentage,
		)
	}

	fmt.Println(
		"============================================================",
	)

	fmt.Println()

	// --------------------------------------------------
	// Fresh engine for every workload.
	// --------------------------------------------------

	engine := buildEngine()

	// Search warmup is only useful for workloads
	// containing search.
	if workload == "search" ||
		workload == "mixed" {

		warmup(engine)
	}

	fmt.Println()

	printHeader()

	for _, concurrency := range concurrencyLevels {

		result :=
			runBenchmark(
				engine,
				workload,
				concurrency,
				initialVectors,
			)

		printResult(result)

		// Allow the system to settle between tests.
		time.Sleep(
			500 * time.Millisecond,
		)
	}
}

// --------------------------------------------------
// Main.
// --------------------------------------------------

func main() {

	fmt.Println(
		"============================================================",
	)

	fmt.Println(
		"VectorEngine Concurrency / Throughput Benchmark",
	)

	fmt.Println(
		"============================================================",
	)

	fmt.Println()

	fmt.Printf(
		"Logical CPUs:  %d\n",
		runtime.GOMAXPROCS(0),
	)

	fmt.Printf(
		"Initial vectors: %d\n",
		initialVectors,
	)

	fmt.Printf(
		"Dimension:       %d\n",
		dimension,
	)

	fmt.Printf(
		"Operations:      %d per level\n",
		totalOperations,
	)

	fmt.Printf(
		"Search K:        %d\n",
		k,
	)

	fmt.Println()

	// --------------------------------------------------
	// Concurrency levels.
	// --------------------------------------------------

	concurrencyLevels :=
		[]int{
			1,
			2,
			4,
			8,
			16,
			32,
		}

	// --------------------------------------------------
	// TEST 1
	// Search only
	// --------------------------------------------------

	runWorkload(
		"search",
		concurrencyLevels,
	)

	// --------------------------------------------------
	// TEST 2
	// Insert only
	// --------------------------------------------------

	runWorkload(
		"insert",
		concurrencyLevels,
	)

	// --------------------------------------------------
	// TEST 3
	// Mixed workload
	// --------------------------------------------------

	runWorkload(
		"mixed",
		concurrencyLevels,
	)

	fmt.Println()

	fmt.Println(
		"============================================================",
	)

	fmt.Println(
		"Benchmark complete.",
	)

	fmt.Println(
		"============================================================",
	)
}
