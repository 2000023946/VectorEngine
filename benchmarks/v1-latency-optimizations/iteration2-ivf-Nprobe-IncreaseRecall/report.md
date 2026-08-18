# v1 IVF Index — Results

## Summary

VectorEngine v1 introduces an IVF (Inverted File) index using Lloyd's algorithm for clustering.

The goal was to determine whether IVF could reduce search latency while maintaining high recall compared with the v0 brute-force baseline.

**Status:**  IVF was not a favorable tradeoff for the current random 128-dimensional workload.

## Experimental Results

| nprobe | Accuracy | 10K Search | 100K Search | 1M Search |
| -----: | -------: | ---------: | ----------: | --------: |
|      5 |   21.50% |    0.18 ms |     2.34 ms |  52.71 ms |
|     25 |   58.86% |    1.05 ms |    13.91 ms | 263.35 ms |
|     65 |   89.26% |    2.54 ms |    35.34 ms | 701.46 ms |

As `nprobe` increases, recall improves because more clusters are searched. However, search latency and memory usage also increase.

## Comparison With Brute Force

| Metric       | Brute Force | IVF @ nprobe 65 |
| ------------ | ----------: | --------------: |
| Accuracy     |     100.00% |          89.26% |
| 10K Search   |     3.72 ms |         2.54 ms |
| 100K Search  |    41.53 ms |        35.34 ms |
| 1M Search    |   483.73 ms |       701.46 ms |
| 1M Memory/op |    15.27 MB |        56.09 MB |
| 1M Allocs/op |           4 |              44 |

IVF provides some improvement at smaller dataset sizes, but at 1M vectors the `nprobe=65` configuration is **slower than brute force** while using substantially more memory and allocations.

## Findings

* Increasing `nprobe` improves recall.
* `nprobe=5` achieved only **21.50%** accuracy.
* `nprobe=25` improved accuracy to **58.86%**.
* `nprobe=65` reached **89.26%**, still below the 90% target.
* Higher `nprobe` substantially increases search cost.
* At 1M vectors, IVF @ 65 probes reached **701.46 ms/search**, compared with **483.73 ms** for brute force.
* The current random 128-dimensional workload does not produce a favorable IVF tradeoff.

## Conclusion

The experiment shows that **our current IVF implementation is not a good fit for VectorEngine's workload**.

This does not mean IVF is inherently ineffective. Rather, with the current random dataset and implementation, increasing the number of probes to recover recall eliminates the expected performance advantage.

The IVF implementation and measurements are retained as an engineering experiment.

**Next:** evaluate a graph-based ANN approach such as **HNSW** and compare it against both the brute-force and IVF baselines using the same benchmark methodology.
