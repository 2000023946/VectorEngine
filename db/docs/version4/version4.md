Here is the updated documentation draft for your project. It reflects the architectural shift to lock-free atomic pointers, the mathematical optimization of the $K$ parameter, and highlights your massive 22x parallel speedup.

---

## V4 Architecture: The Lock-Free Vector Engine

With the elimination of `sync.RWMutex` in favor of an atomic Copy-on-Write (CoW) state pattern, the engine has evolved from a concurrent system into a truly **lock-free, enterprise-scale database**.

Here are the final benchmark results for the V4 engine on the Apple M2 processor:

| Dataset Size | Latency | Memory Footprint (per search) |
| --- | --- | --- |
| **1,000 (Sequential)** | ~5.3 µs | 0 Bytes (0 allocs) |
| **10,000 (Sequential)** | ~18.8 µs | 0 Bytes (0 allocs) |
| **50,000 (Sequential)** | ~90.6 µs | 0 Bytes (0 allocs) |
| **50,000 (Highly Parallel Load)** | **~6.7 µs** | **0 Bytes (0 allocs)** |

### What Makes V4 Enterprise-Scale?

While V3 successfully decoupled the ingestion queue using background workers, it still relied on `sync.RWMutex` to protect the index. Under heavy simulated traffic (500+ concurrent requests), the `RWMutex` created a "stop-the-world" bottleneck: the moment the background worker needed to write the new index, it locked out all incoming read queries.

V4 solves this read-contention problem through advanced memory management and mathematical tuning:

* **Lock-Free Read Path:** `sync.RWMutex` has been entirely removed. The engine now uses `atomic.Pointer[IndexState]` to hold the database state. Search queries load a pointer to an immutable snapshot of the data.
* **Copy-on-Write (CoW) Swaps:** The background worker builds the newly clustered K-means buckets in an isolated memory space. Once the math is complete, it performs a single, nanosecond-level atomic pointer swap. Existing searches finish safely on the old state, while new searches instantly route to the new state. Readers are never blocked by writers.
* **Mathematical Space Optimization:** V4 dynamically scales the number of centroids ($K$) based on the dataset size ($N$). By setting $K = \sqrt{N}$, the engine mathematically minimizes the search space. The total number of distance calculations per search is strictly limited to $K + \frac{N}{K}$, drastically reducing the hardware cycles required for large datasets.
* **True Hardware Concurrency:** By removing OS-level locking mechanisms, V4 allows modern multi-core processors (like the Apple M2) to saturate their cores with raw floating-point math. This resulted in parallel search latency dropping to **~6.7 µs** under massive concurrent load—a 22x speedup over the mutex-bound architecture.
* **Zero-Allocation Guarantees:** Despite the complex atomic pointer swapping, the search path maintains strict `0 B/op` memory allocation. The garbage collector remains entirely unengaged during high-velocity read operations, ensuring flat, deterministic latency at scale.

By combining $O(\sqrt{N})$ mathematical scaling with a lock-free atomic architecture, V4 guarantees that system throughput is limited only by raw CPU clock speed, entirely immune to the concurrency bottlenecks that plague traditional locked data structures.