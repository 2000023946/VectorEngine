Updating centroids synchronously is the exact bottleneck that kills database throughput. If you run Lloyd's algorithm every time a new vector is inserted, your gRPC endpoints will freeze up, and your ingestion latency will spike.

To keep the pipeline moving at cloud scale, you have to decouple the read/write paths from the clustering logic. Here is how you architect an asynchronous, event-driven update loop in your Go backend.

## The Asynchronous Centroid Pipeline

Instead of updating the inverted file index immediately, you buffer the changes and let a background worker handle the heavy mathematical lifting.

1. **The Event Queue (Write-Ahead Log):**
When a new vector arrives via your FastAPI client to the Go backend, don't recluster. Simply append the raw vector to a fast, isolated asynchronous queue (like an in-memory buffer or a durable write-ahead log) and return a `200 OK` success immediately. The primary database continues serving searches using the *existing* centroids.


2. **The Background Worker (Threshold Trigger):**
A dedicated Go goroutine constantly monitors that queue. It waits for a specific threshold—either a time limit (e.g., every 5 minutes) or a volume limit (e.g., after 10,000 new vectors are queued). Once triggered, this worker wakes up and pulls the batch of new vectors.


3. **Out-of-Band Reclustering:**
The background worker takes a snapshot of the current centroids and runs Lloyd's algorithm against the new batch in an isolated memory space. Because this is happening in a background goroutine, your main web API remains 100% responsive to incoming search queries.


4. **The Atomic Swap:**
Once the background worker finishes calculating the new, optimized centroids, it needs to update the live database. Instead of locking the whole database (which blocks reads), you use an atomic pointer swap. You instantly flip the memory reference from the old centroid map to the new one.


5. **The Hardware Sync:**
Immediately after the Go backend swaps to the new centroids, it fires an event down your hardware interface to push the updated centroid coordinates into the FPGA's local SRAM. The hardware router now begins routing incoming physical data streams based on the newly optimized clusters.


---
