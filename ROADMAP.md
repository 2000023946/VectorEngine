# VectorEngine Roadmap

The goal is to build a vector database incrementally, where each architectural layer is introduced to solve a measured problem.

## 1. Brute Force

**Goal:**
Build the simplest correct vector search engine.

**Question:**
How does search performance scale as the dataset grows?

**Experiment:**
Measure initialization, insertion, search latency, memory, and allocations at 10K, 100K, and 1M vectors. Establish the baseline for all future versions.

---

## 2. Performance + Indexing

**Goal:**
Improve the single-node engine and reduce search cost using an approximate nearest-neighbor index.

**Question:**
How much can we improve performance, and how much faster can indexing make search while maintaining high recall?

**Experiment:**
Optimize memory usage, allocations, and search operations. Then compare brute force against IVF at different `nprobe` values using the same benchmark and accuracy suite.

---

## 3. Concurrency

**Goal:**
Allow multiple clients to search and modify the database concurrently while safely maintaining the index.

**Question:**
How does the engine behave under concurrent searches, inserts, and index maintenance?

**Experiment:**
Measure throughput, p50, p95, p99 latency, errors, contention, and race safety at increasing concurrency levels. Test concurrent index builds and background index updates.

---

## 4. Durability

**Goal:**
Make VectorEngine recoverable after a process crash.

**Question:**
Can the database recover committed data after losing RAM?

**Experiment:**
Introduce WAL and persistent storage. Crash the engine during writes and verify that committed data can be recovered correctly.

---

## 5. Distribution

**Goal:**
Scale VectorEngine across multiple machines.

**Question:**
Does adding nodes increase throughput and capacity?

**Experiment:**
Compare one-node, two-node, and three-node deployments while measuring throughput, latency, failures, and resource usage.

---

## Final Goal

A distributed, concurrent, durable vector database whose architecture, tradeoffs, and performance characteristics are understood through **implementation, experimentation, and measurement**.
