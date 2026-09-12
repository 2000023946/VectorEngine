### VectorEngine Durability Benchmark

|   Dataset | Insert Time | Persistence Wait | Reboot Time |                    Recovered |
| --------: | ----------: | ---------------: | ----------: | ---------------------------: |
|    10,000 |    19.98 ms |           4.58 s |      957 ms |       10,000 / 10,000 (100%) |
|   100,000 |   156.70 ms |          32.77 s |      9.63 s |     100,000 / 100,000 (100%) |
| 1,000,000 |      1.63 s |       5m 17.26 s |  1m 36.86 s | 1,000,000 / 1,000,000 (100%) |

### Batching / Async Persistence

|   Dataset | Insert Before Batching | Insert After Batching |      Improvement |
| --------: | ---------------------: | --------------------: | ---------------: |
|    10,000 |                11.05 s |              19.98 ms | **~553× faster** |
|   100,000 |                      — |             156.70 ms |                — |
| 1,000,000 |                      — |                1.63 s |                — |

**Result:** Batching/coalescing persistence requests reduced the overhead on `Insert()` dramatically. At 10K vectors, insertion improved from **11.05 seconds to 19.98 ms (~553×)**, while durability remained **100%** across all tested dataset sizes.
