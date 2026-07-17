# VectorEngine (v1.0 Baseline)

A high-performance, single-node, in-memory vector search engine written in Go. This engine is designed to ingest and query real-time motion data streamed from an Intel DE10-Lite FPGA (ADXL345 accelerometer) via a Serial-to-HTTP bridge.

This repository represents **Version 1.0**, which implements a brute-force linear search $O(N)$ baseline to establish ground-truth accuracy and performance metrics before implementing an IVF-Flat index.

## 🏗️ Architecture (Phase 1)

The system operates across a 3-tier pipeline:
1. **Hardware Layer (Verilog):** Intel DE10-Lite FPGA reads accelerometer data via SPI and transmits it over UART.
2. **Bridge Layer (Python):** Listens to serial data and POSTs structured JSON to the vector engine.
3. **Database Layer (Go):** This engine. Stores vector embeddings in memory and performs nearest-neighbor searches using Squared Euclidean distance.

## 🚀 Key Features of v1.0
* **Primitive Optimization:** Uses `float32` slices to align with typical hardware sensor precision and reduce memory footprint.
* **Zero-Allocation Search:** The core distance calculation loop performs zero heap allocations, ensuring predictable execution latency.
* **Squared Euclidean Distance:** Avoids expensive `sqrt()` operations during the sorting phase for maximum CPU efficiency.

## 📊 Benchmarks

Tests were executed on an Apple M2 architecture. The benchmark suite queries varying dataset sizes to validate the linear $O(N)$ scaling limitations of a brute-force approach.

| Dataset Size | Latency (ns/op) | Memory Allocated per Search | Allocs |
|--------------|-----------------|-----------------------------|--------|
| 1,000        | 208,933         | ~73 KB                      | 4      |
| 10,000       | 2,803,463       | ~721 KB                     | 4      |
| 50,000       | 15,959,609      | ~3.6 MB                     | 4      |

**Analysis:**
The latency scales linearly by a factor of 10 as the dataset grows, proving true $O(N)$ complexity. At 50,000 vectors, query time approaches ~16ms. To support a 60 FPS hardware data stream (16.6ms per frame) without bottlenecking, query latency must be reduced. Additionally, the memory allocated per query scales linearly because v1.0 duplicates the dataset for sorting. 

## 🛣️ Roadmap: Version 2.0 (IVF-Flat Index)

The upcoming Version 2.0 will introduce an **Inverted File (IVF-Flat)** index to solve the scaling limitations of v1.0:
* Segment vectors into clusters (Voronoi cells) via K-Means routing.
* Reduce search complexity from $O(N)$ to sub-linear time by only scanning the nearest cluster.
* Eliminate the $O(N)$ memory allocation overhead during the search sorting phase by utilizing a fixed-size pre-allocated Min-Heap.

## 🛠️ How to Run

**1. Run Accuracy Tests**
Validates Euclidean math and exact nearest-neighbor sorting logic.
```bash
go test -v
```

**2. Run Performance Benchmarks**
Measures nanosecond latency and memory allocations.

```bash
go test -bench=. -benchmem

```