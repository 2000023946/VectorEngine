### SIMD Distance Benchmark

We tested a 128-dimensional squared-distance calculation over **10 million iterations**.

| Implementation  |  Per call | Speedup vs Go |
| --------------- | --------: | ------------: |
| Go scalar       | **96 ns** |         1.00× |
| Go + ARM64 NEON | **49 ns** |     **1.95×** |
| C/ARM64 INT32   | **58 ns** |     **1.63×** |
| C/ARM64 INT8    | **56 ns** |     **1.57×** |

### Conclusion

The **Go + ARM64 NEON implementation is fastest at 49 ns**.

The C/ARM64 versions are slightly slower because each call crosses the **cgo boundary**, adding overhead.

Therefore, calling cgo once per vector during a 1M-vector search would be inefficient:

```text
1M vectors
   ↓
1M cgo calls
   ↓
large accumulated overhead
```

Instead, if we use C/ARM64, the entire `searchRange` should be one native call:

```text
Go goroutine
     ↓
  one cgo call
     ↓
search_range()
 ├── SIMD distance
 ├── iterate vectors
 └── local top-k
     ↓
  return results
```

**Decision:** Keep the current Go + NEON distance kernel. If we move more work to native ARM64, move the **whole search range**, not individual distance calculations.
