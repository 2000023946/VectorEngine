## Durability & Persistence Benchmark

| Dataset | Insert Time | Persistence |  Reboot | Recovery |
| ------: | ----------: | ----------: | ------: | -------: |
|     10K |    19.98 ms |      4.58 s |  0.96 s |     100% |
|    100K |   156.70 ms |     32.77 s |  9.63 s |     100% |
|      1M |      1.63 s |    317.26 s | 96.86 s |     100% |

### Batching Improvement

Batching reduced the number of full database writes by **coalescing queued persistence requests into a single snapshot** rather than performing a disk save for every insert.

The clearest improvement is insert throughput:

* **10K:** ~11.05 s → **19.98 ms** → **~553× faster**
* **100K:** ~156.7 ms in the final implementation
* **1M:** **1.63 s** for all 1M inserts

### Durability Result

All tested datasets achieved **100% recovery** after the in-memory engine was discarded and reconstructed from disk.

**Key result:** asynchronous batched persistence preserves durability while keeping disk I/O off the critical path of `Insert()`.
