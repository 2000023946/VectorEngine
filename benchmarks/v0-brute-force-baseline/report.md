# v0 Brute Force — Results

## Summary

VectorEngine v0 uses brute-force search: every query is compared against every stored vector.

**Status:** Baseline established. The 10K accuracy test passes; the 100K accuracy test currently times out while generating exact ground truth.

## Accuracy

| Dataset | Queries |  K | Accuracy |  Status |
| ------: | ------: | -: | -------: | :-----: |
|     10K |   1,000 | 10 |  100.00% |   PASS  |
|    100K |   1,000 | 10 |        — | TIMEOUT |

The 10K test achieved **100% accuracy**. The 100K test exceeded Go's 10-minute test timeout while calculating the exact brute-force ground truth.

The timeout occurs in the test's `bruteForce()` ground-truth calculation, not in the VectorEngine search implementation.

## Performance

| Operation | Dataset |   Time/op | Memory/op | Allocs/op |
| --------- | ------: | --------: | --------: | --------: |
| Insert    |     10K |  14.59 ms |  11.71 MB |    10,019 |
| Insert    |    100K | 145.40 ms | 120.16 MB |   100,029 |
| Insert    |      1M |    1.57 s |   1.20 GB | 1,000,040 |
| Search    |     10K |   3.72 ms | 163.93 KB |         4 |
| Search    |    100K |  41.53 ms |   1.53 MB |         4 |
| Search    |      1M | 483.73 ms |  15.27 MB |         4 |

## Findings

* Search latency increases roughly linearly with dataset size.
* 1M-vector search requires approximately **484 ms/query**.
* Search uses only **4 allocations/op**, despite scanning the entire dataset.
* Insert memory and allocation costs grow substantially with dataset size.
* Brute-force search provides exact results, but its scaling is the primary limitation.
* The baseline demonstrates that the engine is functionally correct at 10K while exposing clear scaling limitations at larger datasets.

## Conclusion

v0 successfully establishes the **brute-force performance baseline**.

The next phase will combine **performance optimization and IVF indexing**. The existing benchmark and accuracy framework will be retained so improvements can be measured directly against this baseline.

Key baseline:

**1M vectors → 483.73 ms/search, 16.06 MB/op, 4 allocations/op.**
