# VectorEngine Roadmap

The goal is to build a vector database incrementally, where each
architectural layer is introduced to solve a measured problem.

## 1. Brute Force

Goal:
Build the simplest correct vector search engine.

Question:
How does search performance scale as the dataset grows?

Experiment:
Measure search latency at 10K, 100K, and 1M vectors.

---

## 2. Indexing

Goal:
Reduce search cost using an approximate nearest-neighbor index.

Question:
How much faster can we make search while maintaining high recall?

Experiment:
Compare brute force against IVF at different nprobe values.

---

## 3. Concurrency

Goal:
Allow multiple clients to search and modify the database concurrently.

Question:
How does the engine scale under concurrent workloads?

Experiment:
Measure throughput, p50, p95, p99 latency, and errors
at increasing concurrency levels.

---

## 4. Durability

Goal:
Make VectorEngine recoverable after a process crash.

Question:
Can the database recover committed data after losing RAM?

Experiment:
Crash the engine during writes and verify recovery using WAL.

---

## 5. Distribution

Goal:
Scale VectorEngine across multiple machines.

Question:
Does adding nodes increase throughput and capacity?

Experiment:
Compare one-node, two-node, and three-node deployments.

---

## Final Goal

A distributed, concurrent, durable vector database whose
architecture and performance characteristics are understood
through implementation and measurement.