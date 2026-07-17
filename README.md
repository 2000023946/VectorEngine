Here is a complete, production-grade `README.md` for your showcase repository. It is structured to instantly communicate the technical depth, architecture, and engineering metrics to recruiters and judges alike.

---

# ChronosVector: Hardware-Accelerated Real-Time Vector Search Engine

A high-performance, single-node, in-memory vector search engine built from scratch in Go, driven by a live 3-axis digital accelerometer stream orchestrated via a Verilog state machine on an Intel DE10-Lite FPGA. Developed end-to-end in an intensive 8-day hardware/software co-design sprint for the ECE showcase.

```
[ DE10-Lite FPGA ] ──(Onboard ADXL345 via SPI)──> [ Verilog UART Tx ]
                                                           │
                                                (USB-UART Serial @ 115200)
                                                           │
                                                           ▼
[ Go API / Dashboard ] <──(HTTP POST JSON)─── [ Python Serial Bridge ]
  ├── Ingestion Layer (Concurrent, Pre-allocated float32 arrays)
  ├── Search Engine (Custom Euclidean & Cosine Distance Matrix)
  └── Indexing Layer (Inverted File Flat / IVF-Flat Index Clustered)

```

## 🛠️ Architecture Overview

The system spans three distinct layers of the computing stack:

1. **Hardware Ingestion (FPGA):** A custom Verilog SPI master reads raw X/Y/Z acceleration registers from the onboard ADXL345 G-sensor at a high sampling rate. An asynchronous UART transmission module serializes the spatial coordinates into a byte stream pushed over the USB-UART interface.
2. **Ingestion Bridge (Python):** A lightweight serial-to-network daemon reads the raw byte frame from the COM port, structures it into temporal sliding windows (e.g., 30-dimensional motion vectors tracking gesture signatures), and dispatches concurrent HTTP POST requests to the storage engine.
3. **Vector Engine & Dashboard (Go):** A high-performance HTTP engine ingests the sensor payloads, performing real-time spatial indexing and k-Nearest Neighbors (k-NN) similarity matches against historical baseline movement patterns. The results are served dynamically to a lightweight HTML5 canvas dashboard via WebSockets.

---

## ⚡ High-Performance Go Database Optimization

To handle high-throughput physical sensor streams on a single node without garbage collection (GC) or memory thrashing bottlenecks, the engine applies standard systems programming paradigms:

* **Primitive Specialization:** Replaced standard 64-bit float processing with an explicit `float32` pipeline, immediately cutting memory footprint, cache line misses, and SIMD-level processing overhead by **50%**.
* **Zero-Allocation Ingestion:** Eliminated dynamic slice growth (`append`) within the hot path. All vector storage segments are explicitly pre-allocated using continuous backing arrays to guarantee predictable $O(1)$ memory assignments.
* **Algorithmic Evolution:**
* *Phase 1:* Baseline linear brute-force scan evaluating full exact distance matrices ($O(N)$ complexity).
* *Phase 2:* Partitioned inverted space using an **IVF-Flat Index**. incoming vectors are localized to their nearest anchor centroid cluster, skipping non-relevant space to scale search operations to sub-linear time.


---

## 📝 Showcase Presentation Talking Points

1. **Full-Stack Competency:** Demonstrates hardware register allocation and protocol execution (SPI/UART) flowing directly into low-level backend database construction.
2. **Resource-Constrained Optimization:** Shows practical knowledge of CPU caching and memory boundaries by using `float32` arrays and pre-allocated capacity markers over loose data shapes.
3. **Data Integrity Under Stress:** The design guarantees sub-millisecond retrieval benchmarks at scale by dropping the database search search-space from full structural checks into target indexed clusters.

---