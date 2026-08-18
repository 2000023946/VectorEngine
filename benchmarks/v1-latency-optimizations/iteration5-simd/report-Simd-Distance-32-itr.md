## VectorEngine — SIMD Optimization Summary

Implemented an **ARM64 NEON SIMD distance kernel** in Go assembly and added `float64 → int32` quantization to enable supported integer SIMD operations.

### Performance Results

| Benchmark   |  Baseline |      SIMD |                 Change |
| ----------- | --------: | --------: | ---------------------: |
| Search 10K  |   2.27 ms |  0.704 ms |         **69% faster** |
| Search 100K |  22.53 ms |   7.79 ms |         **65% faster** |
| Search 1M   | 236.13 ms |  71.66 ms | **70% faster / 3.29×** |
| Insert 10K  |  0.473 ms |  1.465 ms |        **210% slower** |
| Insert 100K |   4.89 ms |  14.65 ms |        **200% slower** |
| Insert 1M   | 169.06 ms | 298.94 ms |         **77% slower** |

The insert slowdown comes primarily from **quantization being performed during insertion**. This is a deliberate tradeoff: pay the conversion cost once at ingestion to make the much hotter search path substantially faster.

### Accuracy

* **Baseline:** Accuracy was not measured.
* **SIMD + quantization:** **99.96%**
* Correct results: **9,996 / 10,000**

So the optimization achieved roughly **3.3× faster 1M-vector search** while maintaining **99.96% measured top-k accuracy**.

### Quantization

```text
float64 vector
      ↓
quantize during Insert()
      ↓
int32 vector
      ↓
ARM64 NEON SIMD search
```

Quantization also reduces storage from **8 bytes → 4 bytes per dimension**, cutting the raw vector-data footprint in half.
