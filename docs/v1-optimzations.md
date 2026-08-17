## VectorEngine — Optimization Summary

The VectorEngine was optimized from a straightforward brute-force implementation into a **cache-friendly, SIMD-accelerated search engine**.

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

* Converted `float64` vectors to `int32` during insertion.
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

### 6. SIMD Iteration Experiment

* Tested increasing the amount of data processed per loop to further reduce loop iterations.
* The larger SIMD iteration did not produce a consistent performance improvement, so the original NEON implementation was retained.

### Final Architecture

```text
Query
  ↓
Quantize once
  ↓
Contiguous int32 vector storage
  ↓
ARM64 NEON distance kernel
  ↓
Squared distance
  ↓
Top-k selection
```

The main optimization path is therefore **memory layout → fewer allocations → quantization → NEON SIMD distance computation**, while keeping the underlying search algorithm as exact brute-force search.
