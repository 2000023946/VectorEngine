The **batching experiments are not worth pursuing**. The data shows that batching 2, 8, or 16 search requests at a time does not produce a meaningful or consistent throughput improvement, while adding extra complexity through channels/goroutines.

### Batching Results

| Batch Size       |              Result |
| ---------------- | ------------------: |
| **1 (baseline)** | ~27–38 searches/sec |
| **2**            | ~25–36 searches/sec |
| **8**            | ~25–33 searches/sec |
| **16**           | ~28–35 searches/sec |

The important point is that **larger batches do not consistently outperform the baseline**. Across the repeated runs, throughput moves around because the workload is dominated by scanning **1 million × 128-dimensional vectors**, not by the overhead of dispatching individual requests.

For example, in the final baseline run:

* Batch 1 / concurrency 1: **27.81 ops/sec**
* Batch 1 / concurrency 4: **32.41 ops/sec**
* Batch 1 / concurrency 8: **28.15 ops/sec**
* Batch 1 / concurrency 16: **27.68 ops/sec**
* Batch 1 / concurrency 32: **32.82 ops/sec**

There is no clear scaling trend.

### Decision

**We will not implement batching.**

The experiment shows that batching **2 vs. 8 vs. 16 vectors/requests does not provide enough benefit to justify the additional channel/goroutine complexity**. 
