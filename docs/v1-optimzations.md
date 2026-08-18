## VectorEngine — Optimization Summary

The VectorEngine was optimized from a straightforward brute-force implementation into a **cache-friendly, SIMD-accelerated, parallel search engine** while keeping the underlying search algorithm as exact brute-force search.

### 1. Memory Layout

* Replaced scattered/per-vector allocations with a **preallocated contiguous memory buffer**.
* Vectors are stored sequentially, improving cache locality and reducing pointer chasing.
* `Reset()` reuses the existing memory instead of reallocating it.

### 2. Allocation Reduction

* Insert writes directly into the preallocated buffer.
* Search only maintains the required top-k results.
* Removed unnecessary allocations from the hot path.

### 3. Distance Calculation

* Changed Euclidean distance to **squared Euclidean distance**, eliminating the unnecessary square-root operation.
* The ordering of nearest neighbors remains unchanged.

### 4. Quantization

* Converted `float64` vectors to **int32** during insertion.
* This was necessary because the ARM64 Go assembler support available for the implementation did not provide the required floating-point SIMD arithmetic.
* The quantized representation also reduces the memory required per vector.

```text
float64 input
    ↓
quantization
    ↓
contiguous int32 storage
    ↓
NEON SIMD distance calculation
```

### 5. ARM64 NEON SIMD

* Replaced scalar distance computation with a custom **ARM64 NEON assembly kernel**.
* Multiple vector dimensions are processed simultaneously rather than one dimension at a time.
* The main SIMD implementation processes four `int32` values per loop iteration.
* Experiments with smaller integer representations showed that a smaller datatype does not automatically produce better performance because of additional widening/conversion work required by the arithmetic pipeline.

### 6. SIMD Iteration Experiment

* Tested increasing the amount of data processed per loop to further reduce loop iterations.
* The larger SIMD iteration did not produce a consistent performance improvement, so the original NEON implementation was retained.

### 7. Parallel Dataset Search

The final optimization was **parallelizing the brute-force dataset scan across CPU cores**.

Instead of one goroutine scanning the entire dataset:

```text
                 Dataset
                   │
                   ▼
              ┌─────────┐
              │ Search  │
              │ entire  │
              │ dataset │
              └─────────┘
```

the dataset is divided into independent partitions:

```text
                    Dataset
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
      Partition 1  Partition 2  Partition 3
          │            │            │
       Worker 1     Worker 2     Worker 3
          │            │            │
        Top-K        Top-K        Top-K
          │            │            │
          └────────────┼────────────┘
                       ▼
                    Merge
                       │
                       ▼
                  Final Top-K
```

Each worker:

1. Receives an independent range of vectors.
2. Computes distances using the shared read-only query.
3. Maintains its own local Top-K.
4. Returns its local results.

The local Top-K results are then merged into the final Top-K.

This provides **parallelism without changing the underlying brute-force algorithm**. Every vector is still examined, so search remains exact.

### Final Architecture

```text
                         Query
                           ↓
                    Quantize once
                           ↓
              Contiguous int32 storage
                           ↓
                 Partition dataset
                           ↓
        ┌──────────┬──────────┬──────────┐
        ▼          ▼          ▼          ▼
     Worker 1   Worker 2   Worker 3   ... Worker N
        │          │          │             │
      NEON       NEON       NEON          NEON
        │          │          │             │
     Local K     Local K    Local K       Local K
        └──────────┴──────────┴─────────────┘
                           ↓
                       Merge Top-K
                           ↓
                    Final Top-K
```

### Overall Optimization Path

**Memory layout → allocation reduction → squared distance → int32 quantization → NEON SIMD → CPU parallelism**

The important point is that **the algorithm never changed**: it is still brute-force exact search. We improved how efficiently the computer executes that algorithm.
