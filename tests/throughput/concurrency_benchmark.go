package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"vectorengine/src"
)

const (
	dimension = 128

	// Number of searches performed for each concurrency level.
	totalRequests = 1000

	// Search top-k.
	k = 10
)

type BenchmarkResult struct {
	Concurrency int

	Requests int

	TotalTime time.Duration

	Throughput float64

	Average time.Duration
	P50     time.Duration
	P95     time.Duration
	P99     time.Duration
}

func generateVector() []float64 {
	vector := make([]float64, dimension)

	for i := range vector {
		vector[i] = rand.Float64()
	}

	return vector
}

func percentile(sorted []time.Duration, percentile float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}

	index := int(float64(len(sorted)-1) * percentile)

	return sorted[index]
}

func runBenchmark(
	engine *src.VectorEngine,
	concurrency int,
) BenchmarkResult {

	// --------------------------------------------------
	// Generate queries before timing.
	//
	// We don't want random vector generation to become
	// part of the search benchmark.
	// --------------------------------------------------

	queries := make([][]float64, totalRequests)

	for i := range queries {
		queries[i] = generateVector()
	}

	latencies := make([]time.Duration, totalRequests)

	// Divide requests among workers.
	requestsPerWorker :=
		(totalRequests + concurrency - 1) / concurrency

	var wg sync.WaitGroup

	startTime := time.Now()

	for worker := 0; worker < concurrency; worker++ {

		start := worker * requestsPerWorker
		end := start + requestsPerWorker

		if end > totalRequests {
			end = totalRequests
		}

		if start >= totalRequests {
			continue
		}

		wg.Add(1)

		go func(start, end int) {
			defer wg.Done()

			for i := start; i < end; i++ {

				requestStart := time.Now()

				engine.Search(
					queries[i],
					k,
				)

				latencies[i] =
					time.Since(requestStart)
			}
		}(start, end)
	}

	wg.Wait()

	totalTime := time.Since(startTime)

	// --------------------------------------------------
	// Sort latency samples.
	// --------------------------------------------------

	sortedLatencies := make([]time.Duration, len(latencies))

	copy(sortedLatencies, latencies)

	sort.Slice(
		sortedLatencies,
		func(i, j int) bool {
			return sortedLatencies[i] <
				sortedLatencies[j]
		},
	)

	// --------------------------------------------------
	// Calculate average.
	// --------------------------------------------------

	var totalLatency time.Duration

	for _, latency := range latencies {
		totalLatency += latency
	}

	average :=
		totalLatency / time.Duration(len(latencies))

	// --------------------------------------------------
	// Calculate throughput.
	//
	// requests / total elapsed seconds
	// --------------------------------------------------

	throughput :=
		float64(totalRequests) /
			totalTime.Seconds()

	return BenchmarkResult{
		Concurrency: concurrency,

		Requests: totalRequests,

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

func printResult(result BenchmarkResult) {

	fmt.Printf(
		"%8d | %8d | %12.2f | %10s | %10s | %10s | %10s\n",
		result.Concurrency,
		result.Requests,
		result.Throughput,
		result.Average,
		result.P50,
		result.P95,
		result.P99,
	)
}

func main() {

	fmt.Println("==============================================")
	fmt.Println("VectorEngine Concurrency Benchmark")
	fmt.Println("==============================================")

	fmt.Printf(
		"Vectors:       %d\n",
		1_000_000,
	)

	fmt.Printf(
		"Dimension:     %d\n",
		dimension,
	)

	fmt.Printf(
		"Requests:      %d\n",
		totalRequests,
	)

	fmt.Printf(
		"Top-K:         %d\n",
		k,
	)

	fmt.Println()

	// --------------------------------------------------
	// Build engine.
	// --------------------------------------------------

	engine := src.NewVectorEngine()

	fmt.Println("Building VectorEngine...")

	// Generate enough vectors for the benchmark.
	for i := 0; i < 1_000_000; i++ {

		vector := generateVector()

		if err := engine.Insert(i, vector); err != nil {
			panic(err)
		}
	}

	fmt.Println("Engine ready.")
	fmt.Println()

	// --------------------------------------------------
	// Warmup.
	//
	// This prevents initialization/runtime startup from
	// polluting the measurements.
	// --------------------------------------------------

	fmt.Println("Warming up...")

	for i := 0; i < 50; i++ {
		engine.Search(
			generateVector(),
			k,
		)
	}

	fmt.Println("Warmup complete.")
	fmt.Println()

	// --------------------------------------------------
	// Concurrency levels.
	// --------------------------------------------------

	concurrencyLevels := []int{
		1,
		2,
		4,
		8,
		16,
		32,
	}

	fmt.Println(
		"Concurrency benchmark:",
	)

	fmt.Println()

	fmt.Printf(
		"%8s | %8s | %12s | %10s | %10s | %10s | %10s\n",
		"Workers",
		"Requests",
		"Ops/sec",
		"Average",
		"P50",
		"P95",
		"P99",
	)

	fmt.Println(
		"---------+----------+--------------+------------+------------+------------+------------",
	)

	for _, concurrency := range concurrencyLevels {

		result := runBenchmark(
			engine,
			concurrency,
		)

		printResult(result)
	}

	fmt.Println()
	fmt.Println("Benchmark complete.")
}
