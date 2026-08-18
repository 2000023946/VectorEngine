## Parallel Search Optimization

We parallelized the existing **int32 + NEON brute-force search** by splitting the vector dataset into independent partitions. Each goroutine searches its partition and computes a local Top-K; the local results are then merged into the final global Top-K.

### Single-threaded vs. Parallel

| Search | Int32 + NEON | 4 Workers |     6 Workers |
| ------ | -----------: | --------: | ------------: |
| 10K    |     0.704 ms |  0.235 ms |      0.241 ms |
| 100K   |     7.791 ms |  1.810 ms |      1.850 ms |
| 1M     |    71.665 ms | 20.164 ms | **17.275 ms** |

At **1M vectors**, parallel search with 6 workers reduced latency from **71.7 ms → 17.3 ms**, approximately **4.15× faster**.

### Finding the Optimal Worker Count

| Workers |     Search 1M |  Relative |
| ------: | ------------: | --------: |
|       4 |     20.164 ms |     1.00× |
|   **6** | **17.275 ms** | **1.17×** |
|       8 |     47.032 ms |     0.37× |

**6 workers is currently optimal.** Increasing to 8 workers actually hurt performance because the search is heavily memory-bound; additional workers increase memory/cache contention and scheduling overhead.

### Result

The current baseline is therefore:

> **Int32 quantization + ARM64 NEON SIMD + 6-way parallel brute-force search**

with **99.99% measured accuracy** at 10K vectors.
