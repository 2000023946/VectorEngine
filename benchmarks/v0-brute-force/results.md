# VectorEngine v0 — Brute-Force Benchmark Results

**Version:** v0.1  
**Date:** August 16, 2026  
**Platform:** Apple M2, macOS, arm64  
**Vector Dimension:** 128  
**Search K:** 10  

---

## 1. Overview

VectorEngine v0 uses a brute-force search strategy.

For every query, the engine:

1. Iterates through every stored vector.
2. Calculates the Euclidean distance.
3. Sorts all results by distance.
4. Returns the closest `k` vectors.

This version establishes the performance and correctness baseline for future VectorEngine implementations.

The purpose of this benchmark is not to demonstrate scalability. It is to establish **measured evidence** for the limitations that the indexing layer must solve.

---

## 2. Correctness

### Accuracy Test

The accuracy test uses:

- **10,000 vectors**
- **1,000 queries**
- **128 dimensions**
- **K = 10**
- Exact brute-force search as the ground truth

### Result

| Metric | Result |
|---|---:|
| Queries | 1,000 |
| Expected results | 10,000 |
| Correct results | 10,000 |
| Accuracy | **100.00%** |

### Conclusion

The brute-force implementation produces the exact nearest-neighbor results for the tested workload.

**Accuracy Grade: A+**

This gives us a reliable correctness baseline for future approximate indexes such as IVF.

---

# 3. Performance Results

## Insert Performance

| Dataset | Time / op | Grade | Memory / op | Grade | Allocs / op | Grade |
|---:|---:|:---:|---:|:---:|---:|:---:|
| 10K | 14.22 ms | A | 11.71 MB | A | 10,019 | D |
| 100K | 144.86 ms | A | 120.17 MB | B | 100,030 | F |
| 1M | 1.58 s | B | 1.20 GB | F | 1,000,039 | F |

### Insert observations

Insertion time scales approximately linearly with the number of vectors.

However, memory allocation and allocation count grow substantially with dataset size.

At 1M vectors:

- ~1.58 seconds to build the dataset
- ~1.20 GB of allocation activity
- ~1 million heap allocations

The allocation behavior is therefore an important optimization target for future versions.

---

# 4. Search Performance

| Dataset | Time / op | Grade | Memory / op | Grade | Allocs / op | Grade |
|---:|---:|:---:|---:|:---:|---:|:---:|
| 10K | 3.65 ms | A | 163.93 KB | A | 4 | A |
| 100K | 41.27 ms | B | 1.53 MB | B | 4 | A |
| 1M | 453.68 ms | C | 15.27 MB | C | 4 | A |

### Search observations

Search latency increases substantially as the dataset grows:

```text
10K       →   3.65 ms
100K      →  41.27 ms
1M        → 453.68 ms