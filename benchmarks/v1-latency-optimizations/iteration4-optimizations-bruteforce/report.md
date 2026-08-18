
# VectorEngine — Brute-Force Optimization

## Baseline

The original VectorEngine used per-vector allocations and scattered `[]float64` storage.

Search was exact brute force:

```text
O(N × D)
```

## Optimizations

### 1. Flat Contiguous Memory

Preallocated:

```text
2,000,000 × 128 × 8 bytes ≈ 2.048 GB
```

All vectors are stored sequentially in one `[]float64` buffer.

Benefits:

* Better cache locality
* Sequential memory access
* Less pointer chasing
* Lower allocation overhead

### 2. Zero-Allocation Insert

Vectors are copied directly into the preallocated buffer.

```text
Before: ~1 allocation/vector
After:  0 allocations/vector
```

### 3. Squared Distance

Removed `sqrt()` from distance calculations.

Since:

```text
sqrt(a) < sqrt(b) ↔ a < b
```

squared distance preserves nearest-neighbor ordering.

### 4. O(k) Search Memory

Only `k` results are stored instead of allocating results proportional to the dataset size.

### 5. Memory Reuse

`Reset()` clears only the logical vector count.

The 2 GB backing buffer is reused instead of being reallocated.

## Results

| Benchmark   |    Before |     After | Improvement |
| ----------- | --------: | --------: | ----------: |
| Insert 10K  |  14.59 ms |  0.473 ms |   **30.8×** |
| Insert 100K | 145.40 ms |  4.887 ms |   **29.8×** |
| Insert 1M   |   1.566 s |  169.1 ms |    **9.3×** |
| Search 10K  |  3.719 ms |  2.270 ms |   **1.64×** |
| Search 100K |  41.53 ms |  22.52 ms |   **1.84×** |
| Search 1M   | 483.73 ms | 236.13 ms |   **2.05×** |

### Allocation Results

Optimized insertion:

```text
10K:  0 B/op, 0 allocs/op
100K: 0 B/op, 0 allocs/op
1M:   0 B/op, 0 allocs/op
```

## Current Architecture

```text
Query
  ↓
Sequential scan
  ↓
Squared Euclidean distance
  ↓
Top-k results
```

The algorithm is still exact brute force:

```text
Search: O(N × D)
Insert: O(D)
```




