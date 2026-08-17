## VectorEngine — SIMD Optimization Summary

Implemented an **ARM64 NEON SIMD distance kernel** in Go assembly and added `float64 → int32` quantization to enable supported integer SIMD operations.

### SIMD Optimization

The original scalar distance calculation processed **128 dimensions individually**. The NEON implementation processes **4 `int32` values at a time**, reducing the distance calculation from:

```text
128 scalar iterations → 32 SIMD iterations
```

This produced a significant search performance improvement:

| Benchmark   |  Baseline | SIMD (32 iterations) |            Improvement |
| ----------- | --------: | -------------------: | ---------------------: |
| Search 10K  |   2.27 ms |         **0.704 ms** |         **69% faster** |
| Search 100K |  22.53 ms |          **7.79 ms** |         **65% faster** |
| Search 1M   | 236.13 ms |         **71.66 ms** | **70% faster / 3.29×** |

### 32 vs. 16 SIMD Iterations

A second optimization attempted to process **8 `int32` values per iteration**, reducing the loop from **32 → 16 iterations**.

The results were modest and inconsistent:

| Benchmark   |     32 Iterations |    16 Iterations |
| ----------- | ----------------: | ---------------: |
| Search 10K  |        704,147 ns |   **684,756 ns** |
| Search 100K |      7,791,018 ns | **7,241,718 ns** |
| Search 1M   | **71,664,956 ns** |    72,507,161 ns |
| Accuracy    |            99.96% |       **99.98%** |

Because the 16-iteration version did not provide a consistent performance improvement, the **32-iteration NEON implementation was retained**.

### Quantization

Go's ARM64 assembly support did not provide the floating-point SIMD operations needed for the original `float64` representation. The vectors were therefore quantized during insertion:

```text
float64 vector
      ↓
float64 → int32 quantization
      ↓
int32 vector storage
      ↓
ARM64 NEON SIMD distance calculation
```

Quantization reduces vector storage from **8 bytes → 4 bytes per dimension**, cutting the raw vector-data footprint in half. The conversion cost is paid once during insertion, while the search path benefits from integer SIMD operations.

### Final Result

The final implementation reduced the distance calculation from **128 scalar iterations to 32 NEON SIMD iterations**, producing approximately **3.3× faster 1M-vector search** while maintaining **99.96% measured top-k accuracy**.
