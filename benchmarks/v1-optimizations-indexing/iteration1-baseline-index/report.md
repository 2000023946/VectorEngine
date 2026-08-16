# v1 Performance + IVF — Results

## Summary

VectorEngine v1 introduces an IVF (Inverted File) index using Lloyd's algorithm to reduce the number of vectors examined during search.

**Status:** IVF functional, but recall requires tuning.

The initial implementation demonstrates a major search-performance improvement compared with the v0 brute-force baseline, but the initial `nprobe=5` configuration produces only **21.50% recall@10**.

## Accuracy

| Dataset | Queries |  K | Recall@10 | Status |
| ------: | ------: | -: | --------: | :----: |
|     10K |   1,000 | 10 |    21.50% |  FAIL  |

The IVF implementation successfully returns results, but only 2,150 of the 10,000 true top-10 neighbors were recovered.

The low recall is expected to be investigated through `nprobe` tuning and improved clustering.

## Performance

| Operation      | Dataset |   Time/op | Memory/op | Allocs/op |
| -------------- | ------: | --------: | --------: | --------: |
| Initialization |       — | 0.5174 ns |       0 B |         0 |
| Insert         |     10K |    1.10 s |  12.57 MB |    11,435 |
| Insert         |    100K |    2.99 s | 123.29 MB |   101,809 |
| Insert         |      1M |   20.74 s |   1.23 GB | 1,002,510 |
| Search         |     10K |  0.179 ms |  31.90 KB |        18 |
| Search         |    100K |   2.34 ms | 367.78 KB |        24 |
| Search         |      1M |  52.71 ms |   4.37 MB |        33 |

## V0 vs V1 Search

| Dataset | V0 Brute Force |   V1 IVF | Improvement |
| ------: | -------------: | -------: | ----------: |
|     10K |        3.72 ms | 0.179 ms |      ~20.8× |
|    100K |       41.53 ms |  2.34 ms |      ~17.7× |
|      1M |      483.73 ms | 52.71 ms |       ~9.2× |

At 1M vectors, IVF reduces search latency from approximately **483.73 ms to 52.71 ms**, an approximately **89% reduction** in search latency.

## Findings

* IVF substantially reduces search latency.
* Search memory usage is significantly lower than the brute-force baseline.
* Search allocations remain very low.
* The initial `nprobe=5` configuration produces insufficient recall.
* IVF index construction introduces a significant insertion cost when the index is initially built.
* The current implementation establishes that IVF can provide substantial speed improvements, but the speed/recall tradeoff must be tuned.

## Next Experiments

The next V1 experiment will vary `nprobe`:

```text
nprobe = 1
nprobe = 5
nprobe = 10
nprobe = 20
nprobe = 50
nprobe = 100
```

For each configuration, measure:

* Recall@10
* Search latency
* Memory
* Allocations

The goal is to identify a configuration that provides **high recall with substantially lower latency than brute force**.

After the basic IVF behavior is understood, the implementation can be optimized using improved memory layout and incremental centroid/bucket updates.

## Conclusion

v1 demonstrates that IVF can dramatically reduce search cost, achieving approximately **9.2× faster search at 1M vectors** in the initial implementation.

However, the current **21.50% recall@10** is not sufficient. The next step is therefore not to optimize the implementation yet, but to understand and tune the **recall versus search-speed tradeoff** using different `nprobe` values.
