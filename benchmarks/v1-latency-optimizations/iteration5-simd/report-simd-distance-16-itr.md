## SIMD Optimization Comparison

The NEON distance kernel was tested with **32 iterations vs. 16 iterations** per 128-dimensional vector. The 16-iteration version processed twice as many elements per loop, but the benchmark showed that the improvement was modest and inconsistent at larger scale.

| Benchmark     |      32 Iterations |    16 Iterations |      Improvement |
| ------------- | -----------------: | ---------------: | ---------------: |
| Accuracy      |             99.96% |       **99.98%** |         +0.02 pp |
| Accuracy Test |            102.48s |       **79.63s** | **22.3% faster** |
| Search 10K    |         704,147 ns |   **684,756 ns** |  **2.8% faster** |
| Search 100K   |       7,791,018 ns | **7,241,718 ns** |  **7.1% faster** |
| Search 1M     |  **71,664,956 ns** |    72,507,161 ns |      1.2% slower |
| Insert 10K    |   **1,465,248 ns** |     1,555,398 ns |      5.8% slower |
| Insert 100K   |  **14,408,554 ns** |    15,356,556 ns |      6.6% slower |
| Insert 1M     | **298,937,364 ns** |   335,282,611 ns |     10.9% slower |

### Result

The 16-iteration version produced **modest improvements at 10K–100K**, but was slightly slower at 1M. Because the differences are relatively small and benchmark noise can contribute to results at this scale, the **32-iteration implementation was retained as the more conservative baseline** rather than claiming a significant performance improvement.

The main SIMD optimization remains the ARM64 NEON distance kernel, which provides the actual low-level vectorization while maintaining approximately **99.96% retrieval accuracy**.
