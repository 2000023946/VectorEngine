# VectorEngine — API & Durability Architecture

## API Layer

A general-purpose API was added as the interface to the `VectorEngine`.

```text
API
 ├── Insert()
 ├── Search()
 └── NewAPI()
       └── NewVectorEngine()
             └── Reboot()
```

* `Insert()` adds vectors to the engine.
* `Search()` performs vector similarity search.
* `NewAPI()` creates the engine and synchronously calls `Reboot()` before accepting requests.

---

## Asynchronous Persistence

The engine uses a persistence queue and background worker to keep disk I/O off the insert path.

```text
Insert
  ↓
RAM
  ↓
Persistence Queue
  ↓
Persistence Worker
  ↓
Batch Requests
  ↓
Snapshot → Disk
```

* **Queue:** decouples inserts from persistence.
* **Worker:** limits concurrent disk operations.
* **Batching:** combines pending persistence requests to avoid repeated full snapshot writes.
* **Snapshot:** persists vector count, IDs, and quantized vector data.
* **File synchronization:** ensures the completed snapshot is flushed to the OS.

`Insert()` returns after updating RAM and queuing the persistence request; it does not wait for disk I/O.

---

## Concurrency Control

The engine uses an `RWMutex`.

* **Insert / Reset:** exclusive write lock.
* **Search:** shared read lock.
* **Persistence:** briefly acquires a read lock to create a consistent snapshot, then releases it before disk I/O.

This allows searches to safely read the in-memory dataset while preventing modifications from occurring during the snapshot.

---

## Reset & Persistence Race

Because persistence is asynchronous, `Reset()` could otherwise race with an outstanding disk write.

The engine therefore tracks pending persistence operations.

```text
Insert
  ↓
Register persistence
  ↓
Queue request
  ↓
Worker saves
  ↓
Persistence complete
```

When `Reset()` is received:

```text
Reset
  ↓
Block new inserts
  ↓
Wait for pending persistence
  ↓
Reset RAM
  ↓
Remove persisted state
  ↓
Allow inserts again
```

This prevents Reset from invalidating the state while a background persistence operation is still active.

---

## Synchronous Recovery

Recovery occurs during API initialization.

```text
Application Start
       ↓
     NewAPI()
       ↓
     Reboot()
       ↓
   Load snapshot
       ↓
      RAM
       ↓
   API Ready
```

The API is not available until the persisted dataset has been successfully restored into RAM.

---

## Overall Architecture

```text
                    API
                     │
              ┌──────┴──────┐
              │             │
           Insert         Search
              │             │
              ▼             ▼
             RAM ◄─────── Read
              │
              ▼
       Persistence Queue
              │
              ▼
     Background Worker
              │
              ▼
         Batched Snapshot
              │
              ▼
             Disk

Recovery:
Disk → Snapshot → RAM → API Ready
```

**Result:** the VectorEngine provides asynchronous, batched persistence, controlled concurrency, coordinated reset operations, and synchronous recovery while keeping disk I/O off the normal insert critical path.
