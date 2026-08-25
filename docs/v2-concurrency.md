# VectorEngine Concurrency Optimization

## Overview

VectorEngine was tested for concurrent search and insert performance on a 1M-vector dataset.

### Concurrency Results

Search throughput remained around **25–38 ops/sec** as concurrency increased from 1 to 32 workers, while latency increased significantly.

Insert throughput reached approximately **1–3M ops/sec**, showing that concurrent inserts perform well.

## Search Batching

We tested search batching with batch sizes of **2, 8, and 16 vectors at a time** using Go channels and goroutines.

|   Batch Size | Result         |
| -----------: | -------------- |
| 1 (baseline) | ~27–38 ops/sec |
|            2 | ~25–36 ops/sec |
|            8 | ~25–33 ops/sec |
|           16 | ~28–35 ops/sec |

Batching did **not provide a consistent performance improvement**, so it was not worth adding the additional channel/goroutine complexity.

## ARM64 Optimization

We also investigated ARM64/NEON optimization for the search distance calculation.

The goal was to accelerate the vector distance computation using ARM64 SIMD rather than relying on the existing Go implementation.

The ARM64 approach did **not produce a meaningful improvement**, so it was not adopted.

## Conclusion

The concurrency optimization phase is complete. There was no meaningful change from the optimized baseline solution.

* Concurrent search: **implemented and benchmarked**
* Concurrent insert: **implemented and benchmarked**
* Search batching: **tested and rejected**
* ARM64/NEON optimization: **tested and rejected**

The measured results indicate that further optimization should focus on the **search algorithm and indexing**, rather than additional concurrency or batching mechanisms.
