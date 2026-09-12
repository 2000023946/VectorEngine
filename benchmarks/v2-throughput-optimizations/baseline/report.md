# VectorEngine — Concurrency Baseline

**Date:** August 21, 2026
**Dataset:** 1M vectors × 128 dimensions
**CPU:** 8 logical CPUs

## Results

### Search

| Concurrency | Throughput |   P50 |   P95 |   P99 |
| ----------: | ---------: | ----: | ----: | ----: |
|           1 |    27.81/s |  31ms |  66ms | 123ms |
|           4 |    32.41/s | 117ms | 182ms | 241ms |
|           8 |    28.15/s | 260ms | 463ms | 557ms |
|          32 |    32.82/s | 913ms | 1.46s | 1.74s |

**Finding:** More concurrent searches do not significantly increase throughput. Latency rises sharply because a single search already heavily utilizes the CPU/memory subsystem.

### Insert

Single-concurrency insert throughput reached **1.67M ops/sec** with ~473 ns average latency.

**Finding:** Inserts are extremely cheap compared with searches. The mutex provides correctness but limits concurrent write scaling.

### Mixed

At 32 concurrency:

* **30.87 ops/sec**
* **P95: 4.38s**
* **P99: 7.19s**

High concurrency causes severe tail-latency growth.

## Conclusion

The bottleneck is **brute-force search**, not insufficient goroutines.

### Next step

**Micro-batching search**

Test batch sizes:

```text
1 → 2 → 4 → 8 → 16 → 32
```

with a small maximum wait time.

Goal:

> Increase throughput through better cache/memory/SIMD utilization while keeping the added latency small.

After batching, move to **indexing**, which reduces the amount of data scanned per query.
