## VectorEngine — SIMD Optimization & Benchmark Validation

The optimization had **two parts**:

1. **Quantization:** convert the original `float64` vectors to `int32` during insertion so the distance kernel can operate on integer data.
2. **ARM64 NEON:** implement the squared-distance kernel in assembly using SIMD instructions to process multiple `int32` dimensions in parallel.

The architecture became:

```text
float64 input
     ↓
float64 → int32 quantization
     ↓
contiguous int32 storage
     ↓
NEON SIMD distance calculation
     ↓
int64 squared-distance accumulation
```

Quantization reduces the raw vector-data footprint from **8 bytes → 4 bytes per dimension**, while the NEON kernel accelerates the hot search path.

### Isolating the NEON Improvement

Initially, comparing the original `float64` implementation against the optimized version could not tell us how much of the performance improvement came from **quantization** versus **NEON**.

Therefore, we ran a controlled experiment using the **same int32 representation and quantization**, changing only the distance implementation:

* **Scalar:** normal Go loop, one dimension at a time.
* **NEON:** ARM64 assembly using SIMD instructions.

Both use the same `int32` data, `int64` accumulation, brute-force `O(N·D)` search, and top-k selection.

| Benchmark   | Int32 + Scalar Go |  Int32 + NEON |   Speedup |
| ----------- | ----------------: | ------------: | --------: |
| Search 10K  |          1.188 ms |  **0.704 ms** | **1.69×** |
| Search 100K |         11.891 ms |  **7.791 ms** | **1.53×** |
| Search 1M   |        122.775 ms | **71.665 ms** | **1.71×** |
| Accuracy    |            99.95% |    **99.96%** |     ~same |

At 1M vectors, NEON reduces latency from **122.8 ms → 71.7 ms**, a **41.6% reduction in latency** and approximately **1.71× speedup**.

This controlled comparison demonstrates that the performance improvement is specifically attributable to the **NEON SIMD distance kernel**, rather than simply switching from `float64` to `int32`.

### SIMD Implementation

The original scalar distance calculation processes:

```text
128 dimensions
      ↓
128 scalar iterations
```

The NEON implementation processes **4 `int32` values per SIMD iteration**:

```text
128 dimensions
      ↓
32 SIMD iterations
```

An additional 8-values-per-iteration implementation was tested, but it did not provide a consistent improvement at 1M vectors, so the 4-value SIMD implementation was retained.

### Result

The optimization therefore established two distinct benefits:

```text
float64 → int32
    ↓
2× smaller vector representation
+ integer-compatible SIMD computation

int32 → NEON
    ↓
~1.7× faster search than scalar int32
```

The **int32 + NEON implementation is now the validated SIMD baseline** for the next optimization: evaluating whether **int8 quantization + an optimized int8 NEON kernel** can provide further memory and performance improvements.
