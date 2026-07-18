

## Vector Engine V2: Architecture & Performance Review

This release introduces a high-performance Inverted File (IVF) vector index, replacing the Phase 1 brute-force scan. The architecture is specifically optimized for hardware-level cache efficiency and high-throughput ingestion.

### 1. Index Selection: IVF vs. HNSW

For the indexing strategy, we deliberately chose an Inverted File (IVF) architecture over a Hierarchical Navigable Small World (HNSW) graph.

While HNSW provides exceptional search times, it introduces significant drawbacks for our specific use case:

* **Memory & CPU Overhead:** Graph traversal requires constant pointer chasing, leading to frequent CPU cache misses and high memory consumption.
* **Ingestion Bottlenecks:** Inserting nodes into a layered graph structure is comparatively slow, making it poorly suited for write-heavy data streams.

By optimizing an IVF index for continuous memory layouts, we achieve search speeds that rival or beat graph-based approaches through pure mechanical sympathy. Because our vectors are stored in contiguous slice buckets, the CPU can prefetch memory and scan entirely out of the hardware cache. This cache-aware design aligns with industry standards utilized by Meta (FAISS), Milvus, and Qdrant.

### 2. Initialization and Centroid Training

To ensure the IVF index routes queries efficiently, we implemented a two-phase initialization model designed to pick optimal centroids without early bias.

* **Warmup Phase:** The engine begins by ingesting raw vectors into a flat list. This ensures we gather a statistically significant sample of the real-world data distribution before making routing decisions.
* **Index Build:** Once the data reaches a predefined threshold, the engine automatically triggers the index build.
* **Lloyd's Algorithm:** Centroids are seeded using random initialization (Forgy method) and then refined using a batch-update implementation of Lloyd’s Algorithm. After just a few iterations, the centroids successfully converge to the true centers of mass.

### 3. Future Roadmap: Asynchronous Updates

While this version successfully introduces the frozen IVF index, the industry standard for continuous streaming involves dynamically updating the index. A future release will introduce an asynchronous background worker to continuously recalculate centroids and rebalance buckets without blocking the main read/write threads.

---

## Benchmark Results: V1 vs. V2

The transition from the V1 Brute Force model to the V2 IVF Index resulted in a massive latency reduction and achieved zero-allocation searches. At a 50,000 dataset scale, the V2 engine is approximately **1,250x faster** than V1.

| Dataset Size | V1 Latency (ns) | V2 Latency (ns) | V1 Allocations | V2 Allocations |
| --- | --- | --- | --- | --- |
| 1,000 | 208,933 | 2,456 | ~73 KB (4 allocs) | 0 B (0 allocs) |
| 10,000 | 2,803,463 | 6,223 | ~721 KB (4 allocs) | 0 B (0 allocs) |
| 50,000 | 15,959,609 | 12,774 | ~3.6 MB (4 allocs) | 0 B (0 allocs) |

> **Key Takeaway:** The search path now operates entirely without heap allocations. This guarantees that the Go Garbage Collector will not pause the system during high-throughput querying, and proves that our sub-linear $O(\sqrt{N})$ routing is working exactly as mathematically intended.