
---

```md
# VectorEngine

A high-performance in-memory vector storage and search engine built in Go.

VectorEngine is a baseline vector database implementation that supports high-dimensional vector insertion and brute-force nearest neighbor search. It is designed for benchmarking, experimentation, and as a foundation for future approximate nearest neighbor (ANN) indexing systems such as HNSW or IVF.

---

## 🧠 Overview

VectorEngine provides:

- High-dimensional vector insertion (float32)
- Brute-force similarity search (Euclidean distance or custom metric)
- Dynamic in-memory storage using Go slices
- Benchmark-driven performance analysis
- Full unit test coverage for core functionality

This project prioritizes clarity, correctness, and measurable performance over scalability in its current version (v1).

---

## ⚙️ Features

### ✔ Core Functionality
- Insert vectors into an in-memory store
- Search for nearest vectors using full scan comparison
- Configurable vector dimension size

### ✔ Performance Benchmarking
- Insert and search benchmarks using Go testing framework
- Memory allocation tracking (`B/op`, `allocs/op`)
- Nanosecond-level performance measurement

### ✔ Testing
- Unit tests for insert, search, and distance functions
- Edge case validation (dimension mismatch, empty store, etc.)
- Designed for full test coverage of core logic

---

## 📁 Project Structure

```

.
├── bench_test.go
├── distance.go
├── distance_test.go
├── insert_test.go
├── search.go
├── search_test.go
├── storage.go
├── store_test.go
├── main.go
├── coverage.out
├── docs/
│   ├── timeline.md
│   └── versions/
│       └── version1.md
├── run_benchmark_tests.sh
├── run_test_reports.sh

````

---

## 🚀 Core API

### Insert

```go
Insert(vec []float32)
````

Adds a vector to the in-memory store.

* Stores vectors in a dynamically growing slice
* May trigger memory reallocations during growth
* No deduplication or indexing in v1

---

### Search

```go
Search(query []float32) ([]float32, error)
```

Finds the nearest vector using brute-force search.

* Iterates through all stored vectors
* Computes distance per vector
* Returns the closest match
* Time complexity: **O(N × D)**

---

## 🧪 Benchmarks (Apple M2)

### Insert Benchmark

```
BenchmarkInsert
~65 µs/op
~500 KB/op
~16 allocs/op
```

### Search Benchmark

```
BenchmarkSearch
~16 ms/op
0 B/op
0 allocs/op
```

---

## 📊 Performance Summary

| Operation | Time   | Memory     | Allocations | Complexity                            |
| --------- | ------ | ---------- | ----------- | ------------------------------------- |
| Insert    | ~65 µs | ~500 KB/op | 16 allocs   | Amortized O(1) with resizing overhead |
| Search    | ~16 ms | 0 B/op     | 0 allocs    | O(N × D)                              |

---

## 🧠 Key Observations

### Insert

* Fast under normal conditions
* Memory overhead increases due to slice resizing
* Allocation cost driven by dynamic growth strategy

### Search

* Fully allocation-free
* Performance limited by linear scan
* Bottleneck at scale due to O(N) complexity

---

## ⚠️ Limitations

* No approximate nearest neighbor (ANN) indexing
* Linear search does not scale beyond large datasets
* Insert incurs memory overhead due to resizing
* No SIMD or vectorized optimizations

---

## 🛠️ Future Improvements (Roadmap)

### Priority 1: Search Optimization

* Implement HNSW (Hierarchical Navigable Small World graph)
* OR IVF clustering (partition-based search)
* Target improvement: **~16ms → <1ms**

---

### Priority 2: Insert Optimization

* Preallocate storage capacity
* Reduce slice resizing frequency
* Lower allocations from 16 → ~1–3 allocs/op

---

### Priority 3: Memory Optimization

* Reduce vector copying overhead
* Introduce pooling or reuse strategies
* Optimize GC pressure

---

### Priority 4: Profiling & Analysis

* Integrate `pprof` for CPU and memory profiling
* Identify allocation hotspots
* Optimize search hot path

---

## 🧪 Running Tests

Run all tests:

```bash
go test ./...
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

Run coverage report:

```bash
go test ./... -cover
go tool cover -html=coverage.out
```

---

## 📄 Documentation

* Version history: `docs/versions/version1.md`
* Development timeline: `docs/timeline.md`

---

## 🧠 Conclusion

VectorEngine v1 demonstrates a working baseline vector database with:

* Correct insertion and retrieval logic
* Benchmark-driven performance visibility
* Full unit testing of core components

It serves as a foundation for building a scalable vector search system with future support for ANN indexing and production-grade optimizations.

---

```

---
