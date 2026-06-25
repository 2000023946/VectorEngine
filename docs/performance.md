
---

# 🧠 1. Real performance is about memory, not computation

Modern CPUs are very fast at math, but slow when they wait on memory.

So performance is mainly determined by:

> **how efficiently data moves through cache and RAM**

---

# 🧠 2. Cache is the key layer between CPU and RAM

* CPU has small fast memory called **cache (L1/L2/L3)**
* RAM is much slower
* CPU performs best when data is already in cache

---

# 🧠 3. Cache misses are the main performance killer

A **cache miss** happens when:

> CPU needs data but it is not in cache → must fetch from slow RAM

This causes:

* CPU stalls (idle time)
* major slowdowns in loops and graph traversal

---

# 🧠 4. Contiguous memory reduces cache misses

When data is stored sequentially:

```text
arr[0], arr[1], arr[2], ...
```

* CPU loads entire cache lines (64 bytes at a time)
* prefetching works
* access patterns are predictable

Result:

> fewer cache misses → faster execution

---

# 🧠 5. Pointer-heavy structures hurt performance

Structures like:

* `*Node`
* `[][]T` (Go slices of slices)
* graph pointers

cause:

> memory jumps instead of sequential access

This leads to:

* cache misses
* RAM stalls
* CPU idle time

---

# 🧠 6. Flattened arrays are faster than pointer-based layouts

Flattened structure:

```text
[]T (single contiguous block)
```

is faster because:

* no pointer chasing
* predictable indexing
* better cache usage
* enables SIMD/vectorization

Access becomes:

```text
i * width + j
```

---

# 🧠 7. Preallocation improves performance

Instead of allocating during inserts:

* allocate large memory upfront
* write into existing space

Benefits:

* no repeated heap allocations
* less GC pressure
* predictable runtime performance

---

# 🧠 8. Inserts should be simple writes (not allocations)

In FAISS-style systems:

> insert = write into preallocated array slot

NOT:

* create object
* allocate heap memory
* build pointer structures

---

# 🧠 9. Prefetching hides memory latency

CPU can predict sequential access and:

> load future data while processing current data

This creates a pipeline:

* compute current item
* memory loads next items in background

Works only when memory is predictable.

---

# 🧠 10. Core system design principle

Instead of thinking:

> stack vs heap = performance

You should think:

> **predictable contiguous memory > pointer-based random access**

---

# 🚀 Final one-line summary

> High-performance systems like FAISS are fast because they use contiguous memory, preallocation, and index-based access to minimize cache misses and keep the CPU continuously fed with predictable data instead of waiting on RAM.

---
