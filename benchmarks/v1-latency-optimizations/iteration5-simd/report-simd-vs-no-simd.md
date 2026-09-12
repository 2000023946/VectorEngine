### Int32 Distance Kernel: NEON vs. Scalar Go

This experiment isolates the effect of **ARM64 NEON SIMD acceleration**. Both versions use the same `int32` quantized vector representation, the same `int64` squared-distance accumulation, the same brute-force `O(N·D)` search algorithm, and the same top-k selection. The only variable is the implementation of the distance kernel:

* **Scalar:** Go `for` loop performing one dimension at a time.
* **NEON:** ARM64 assembly using SIMD registers to process multiple dimensions in parallel.

| Benchmark   | Int32 + Scalar Go |  Int32 + NEON |   Speedup |
| ----------- | ----------------: | ------------: | --------: |
| Search 10K  |          1.188 ms |  **0.704 ms** | **1.69×** |
| Search 100K |         11.891 ms |  **7.791 ms** | **1.53×** |
| Search 1M   |        122.775 ms | **71.665 ms** | **1.71×** |
| Accuracy    |            99.95% |    **99.96%** |     ~same |

At **1M vectors**, NEON reduces search latency from **122.8 ms → 71.7 ms**, a **41.1% reduction in latency** and approximately **1.71× higher throughput** for the distance computation/search workload.

The accuracy remains effectively identical, confirming that the assembly implementation is producing equivalent distance results.

**Conclusion:** the benchmark demonstrates that the performance improvement is attributable to **NEON SIMD execution**, not to the int32 representation or quantization itself. This establishes **int32 + NEON** as the optimized SIMD baseline against which we can now evaluate **int8 + NEON**.
