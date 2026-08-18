## Int32 vs. Int8 Quantization

We evaluated moving the VectorEngine from **int32 to int8 storage** to reduce memory usage.

With int32:

```text
128 dimensions × 4 bytes = 512 bytes/vector
```

With int8:

```text
128 dimensions × 1 byte = 128 bytes/vector
```

So int8 provides a **4× reduction in raw vector storage**. However, this comes with a computational tradeoff.

### Why int8 increases computation

The stored values are `int8`, but squared Euclidean distance requires a wider representation:

```text
int8
  ↓
widen to int16/int32
  ↓
subtract
  ↓
square
  ↓
accumulate into int64
```

For example:

```text
127 - (-127) = 254
254² = 64,516
```

Neither the difference nor the squared result fits in `int8`.

Our initial int8 implementation therefore required additional **sign extension, widening, and scalar operations** during the distance calculation. This largely offset the cache and memory-bandwidth advantage of the smaller representation.

### Benchmark

| Metric      |  Int32 + NEON |  Int8 + NEON |
| ----------- | ------------: | -----------: |
| Accuracy    |    **99.96%** |       98.03% |
| Search 10K  |  **0.704 ms** |     1.198 ms |
| Search 100K |  **7.791 ms** |    11.968 ms |
| Search 1M   | **71.665 ms** |   119.374 ms |
| Init memory |       1.04 GB |   **272 MB** |
| Insert 1M   |      298.9 ms | **231.6 ms** |

Int8 successfully reduced memory usage by roughly **74%** and improved large-scale insertion performance. However, search became approximately **1.67× slower at 1M vectors**, while accuracy dropped from **99.96% to 98.03%**.

### Decision

For the current VectorEngine, **int8 is not worth the tradeoff**.

The primary workload is vector search, and the additional computation required by int8 distance evaluation outweighs its memory advantage. Therefore, we will retain:

```text
float64 input
      ↓
int32 quantization
      ↓
contiguous int32 storage
      ↓
NEON SIMD L2 distance
```

This gives us a strong balance of **accuracy, memory efficiency, and search performance**.

We can revisit int8 later if we implement a genuinely SIMD-native int8 distance kernel, but **for the current version, int32 + NEON is the better engineering choice**.
