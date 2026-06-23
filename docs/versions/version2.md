
---

# 📄 VectorEngine — Version 2 Report

## 🧠 Overview

VectorEngine v2 introduces a major architectural shift from a brute-force vector database to a **graph-based Approximate Nearest Neighbor (ANN) system**.

Instead of scanning all vectors during search (v1), this version builds a **navigation graph structure** that enables fast greedy traversal toward nearest neighbors.

This is the first step toward an **HNSW-inspired indexing system**, focusing on improving search performance at the cost of more expensive insertion.

---

# ⚙️ Core Architecture

## 🧩 Storage Model

VectorEngine v2 uses:

* In-memory graph structure
* Nodes represent vectors
* Edges represent similarity-based neighbor links
* Fixed dimensional vectors (e.g., 128D)

```go
type Node struct {
	ID        int
	Vector    []float32
	Neighbors []int
}
```

---

## 🧠 Graph Structure

```go
type Graph struct {
	Nodes     map[int]*Node
	K         int
	Dimension int
	LastID    int
}
```

### Key properties:

* `Nodes`: adjacency map of vector graph
* `K`: max neighbors per node (graph sparsity control)
* `LastID`: starting entry point for traversal
* `Dimension`: vector size constraint

---

# 🚀 Core API

---

## 📥 Insert (Graph-based indexing)

```go
Insert(vec []float32) error
```

### Description:

Inserts a vector into the graph and builds local neighborhood connections.

### Behavior:

1. Finds approximate insertion point using **greedy traversal**
2. Collects candidate neighbor set (current node + neighbors)
3. Computes Euclidean distances
4. Sorts candidates by similarity
5. Connects new node to top-K closest nodes
6. Updates reverse connections (bidirectional graph)

---

### Key characteristics:

* Uses **local graph exploration instead of global scan**
* Builds **sparse neighbor graph**
* Maintains **K-bounded connectivity**
* Enables efficient future navigation

---

### Complexity:

| Step                 | Cost                   |
| -------------------- | ---------------------- |
| Traverse             | O(log N-ish effective) |
| Candidate collection | O(K)                   |
| Distance computation | O(K × D)               |
| Sorting              | O(K log K)             |

---

### Tradeoff:

> Insert became more expensive to enable faster search

---

# 🔍 Search (Graph traversal ANN)

```go
Search(query []float32) (Result, error)
```

---

## Description:

Performs nearest neighbor search using **greedy graph navigation**.

Instead of scanning all nodes, it:

* starts from `LastID`
* moves to better neighbors iteratively
* stops when no improvement is found

---

## Result:

```go
type Result struct {
	ID       int
	Distance float32
}
```

---

## Behavior:

1. Start from entry node (`LastID`)
2. Compute distance to query
3. Explore neighbors
4. Move to closer node if found
5. Stop when no improvement exists

---

## Complexity:

* Theoretical: O(N × D)
* Practical: ~O(log N) to small local region traversal

---

## Key improvement over v1:

> Search no longer scans the full dataset

Instead, it explores:

* only a small connected region of the graph
* highly localized node clusters

---

# 📊 Benchmark Results (Apple M2)

---

## 🟢 Insert Benchmark

```
BenchmarkInsert-8
38 iterations
26944160 ns/op
9006894 B/op
109969 allocs/op
```

### Interpretation:

* ⏱ ~26.9 ms per insert
* 💾 ~9.0 MB memory allocation per insert batch
* 🧠 ~110,000 allocations per run

---

### Root cause:

Insert is expensive due to:

* greedy traversal before insertion
* candidate set construction (map usage)
* repeated Euclidean distance computations
* full sorting of candidate list
* bidirectional edge updates

---

## 🟢 Search Benchmark

```
BenchmarkSearch-8
317,827 iterations
3813 ns/op
520 B/op
5 allocs/op
```

### Interpretation:

* ⏱ ~3.8 microseconds per search
* 💾 ~520 bytes memory per query
* 🧠 5 allocations per query

---

# 📊 Performance Summary

| Operation | Time     | Memory | Allocations | Complexity              |
| --------- | -------- | ------ | ----------- | ----------------------- |
| Insert    | ~26.9 ms | ~9 MB  | ~110k       | Graph construction cost |
| Search    | ~3.8 µs  | ~520 B | 5           | Local graph traversal   |

---

# 🧠 Key Improvements Over Version 1

---

## 🚀 1. Search Performance Revolution

| Version | Search Time |
| ------- | ----------- |
| v1      | ~16 ms      |
| v2      | ~3.8 µs     |

### Why?

* eliminated full dataset scan
* replaced with graph traversal
* reduced search space dramatically

---

## 🧩 2. Introduced ANN-style indexing

VectorEngine now behaves like:

> a simplified HNSW-like graph index

Core ideas introduced:

* neighbor links
* greedy navigation
* locality-based search
* bounded degree (K)

---

## 🧠 3. Structured memory layout

Instead of flat list:

* graph-based representation
* adjacency relationships
* navigable structure

---

# ⚠️ Current Limitations

---

## 1. Insert is extremely expensive

* high allocation count (~110k ops)
* full candidate sorting every insert
* repeated distance computations

---

## 2. No multi-layer hierarchy yet

* single-layer graph only
* no HNSW-level skipping structure

---

## 3. Greedy search is not guaranteed optimal

* can get stuck in local minima
* depends heavily on entry node (`LastID`)

---

## 4. No beam search / exploration width

* only follows single best path

---

# 🔬 Design Classification

VectorEngine v2 is best described as:

> 🟡 Single-layer greedy ANN graph (HNSW-inspired prototype)

---

# 🚀 Recommended Next Version (v3 direction)

---

## Priority 1: Optimize Insert (critical)

* reduce allocations
* reuse buffers
* replace map with slices where possible
* avoid full sorting (use partial selection)

---

## Priority 2: Add Beam Search

* explore multiple candidate paths
* improve recall stability
* reduce local minima issues

---

## Priority 3: Introduce HNSW structure

* multi-layer graph
* probabilistic level assignment
* efConstruction / efSearch tuning

---

## Priority 4: Profiling

* use `pprof`
* identify:

  * allocation hotspots
  * Euclidean distance cost
  * sorting overhead

---

# 🧾 Conclusion

VectorEngine v2 successfully transforms the system from a:

> ❌ brute-force vector database (v1)

into a:

> 🟢 graph-based ANN index (v2)

This version achieves:

* **massive search speed improvement**
* **structured vector navigation**
* **foundation for HNSW implementation**

However, it introduces:

* high insert cost
* memory overhead
* need for optimization before scaling

---

# 🚀 Final Summary

> v1 = simple + slow search
> v2 = structured graph + fast search + expensive insert
> v3 (next) = optimized graph + HNSW-ready system

---

