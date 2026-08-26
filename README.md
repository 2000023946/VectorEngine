
# VectorEngine

VectorEngine is a custom vector-search database implemented in Go. The project explores how low-level storage, SIMD, concurrency, throughput, and durability techniques affect vector-search performance.

## Architecture

```text
                    VectorEngine
                         │
              ┌──────────┴──────────┐
              │                     │
           Insert                 Search
              │                     │
        Quantization          Parallel Search
              │                     │
              ▼                     ▼
       Contiguous RAM        SIMD Distance
              │                     │
              └──────────┬──────────┘
                         │
                  Async Persistence
                         │
                         ▼
                        Disk
```

The engine stores quantized vectors in contiguous memory and uses parallel workers with an ARM64 NEON distance kernel for search.

## Key Features

* Flat contiguous vector storage
* Int32 vector quantization
* ARM64 NEON SIMD distance computation
* Parallel search using independent partitions
* Concurrent read/write synchronization with `sync.RWMutex`
* Asynchronous disk persistence
* Batched persistence through a bounded worker queue
* Synchronous database recovery during startup
* Reset coordination with asynchronous persistence
* Accuracy, performance, throughput, and durability tests

## Project Evolution

### V0 — Brute-Force Baseline

Established the initial vector-search implementation and benchmark baseline.

### V1 — Latency Optimization

Explored:

* IVF indexing
* N-probe tuning
* Graph-based indexing
* Memory layout optimizations
* Quantization
* SIMD distance computation
* Parallel search

### V2 — Throughput Optimization

Focused on improving throughput under concurrent workloads through:

* ARM64 optimization
* Concurrent request processing
* Search parallelism
* Batched workloads

### V3 — Durability

Added persistent storage so the database can recover its in-memory state after shutdown or restart.

Persistence is performed asynchronously through a worker queue with batched disk writes, allowing inserts to update RAM without waiting for disk I/O.

Startup recovery is synchronous: the database loads the persisted state before becoming available for requests.

## Testing

The project includes separate test suites for:

* Unit correctness
* Search accuracy
* Performance
* Concurrent throughput
* Durability and recovery

Detailed benchmark results and implementation notes are available under:

```text
benchmarks/
docs/
```

## Project Structure

```text
src/          Core VectorEngine implementation
tests/        Unit, accuracy, performance, throughput, durability tests
benchmarks/   Benchmark results and experimental iterations
docs/         Architecture and implementation documentation
tools/        Benchmark/testing utilities
```

## Status

VectorEngine is an experimental vector database focused on understanding and measuring systems-level performance tradeoffs in storage, computation, concurrency, and persistence.

## Future Work

### Distributed VectorEngine

The next stage is to extend VectorEngine from a single-node database into a distributed vector-search system.

Planned work includes:

- Sharding vectors across multiple nodes
- Routing queries to the appropriate shards
- Parallel distributed search and result merging
- Replication for improved fault tolerance
- Node health detection and recovery
- Distributed persistence and recovery