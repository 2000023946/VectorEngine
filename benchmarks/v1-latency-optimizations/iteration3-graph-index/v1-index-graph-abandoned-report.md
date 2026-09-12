# v1 Index-Graph Experiment — Abandoned

## Summary

VectorEngine v1 experimented with a graph-based nearest-neighbor index on the same random 128-dimensional vector workload used for the brute-force baseline.

**Status:**  Abandoned

The graph index did not provide a useful performance/accuracy tradeoff on this workload. Index construction became substantially more expensive, while search performance improvements were not sufficient to justify the added complexity.

The experiment is therefore being discontinued. VectorEngine will return to optimizing the exact brute-force engine for the current workload.

## Results

| Operation | Dataset | Brute Force | Graph Index | Result |
| --------- | ------: | ----------: | ----------: | ------ |
| Insert | 10K | 14.59 ms | 318.53 ms | ~21.8× slower |
| Insert | 100K | 145.40 ms | 6.85 s | ~47.1× slower |
| Insert | 1M | 1.57 s | ~357.0 s | ~227× slower |
| Search | 10K | 3.72 ms | 0.249 ms | Faster |
| Search | 100K | 41.53 ms | 0.585 ms | Faster |
| Search | 1M | 483.73 ms | — | Benchmark did not complete |

## Accuracy

The graph index achieved only:

**3.66% accuracy at 10K vectors**

with 1,000 queries and `k=10`.

The required accuracy target was 90%.

## Finding

The graph index reduced search latency at 10K and 100K, but the improvement came with a very large increase in insertion/index-construction cost and poor recall.

For this experiment, the dataset consists of randomly generated 128-dimensional vectors. The experiment therefore did not produce a useful index structure for the workload.

The important engineering result is not that graph indexing is universally ineffective. Rather:

> **For this specific random-vector workload, the graph index added substantial indexing overhead without achieving the required accuracy/performance tradeoff.**

Because the goal of the project is to build a measured system rather than force an indexing strategy that does not work, the graph-index experiment is being abandoned.

## Decision

VectorEngine will not continue developing this index for the current workload.

The project will instead focus on optimizing the brute-force engine through:

1. Memory optimization
2. Computation optimization
3. CPU/cache optimization

The original brute-force benchmark remains the comparison baseline for all subsequent optimizations.

## Baseline Reference

The brute-force baseline achieved:

- **100% accuracy at 10K**
- **3.72 ms/search at 10K**
- **41.53 ms/search at 100K**
- **483.73 ms/search at 1M**
- **4 allocations/search**

The indexing experiment demonstrated that adding an index is not automatically an improvement: the index itself has a construction and maintenance cost that must be justified by the workload.
