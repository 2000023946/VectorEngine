Here is a high-level architectural write-up of our cloud-native distributed search pipeline. You can use this as the foundation for your `README.md`, an enterprise architectural diagram, or project documentation.

This description focuses on the strict domain boundaries and asynchronous routing mechanisms that make the cluster resilient.

---

## Architecture Overview

The system is designed as a highly available, **cloud-native distributed search pipeline**. By strictly separating the stateless routing logic from the stateful database processing, the architecture enables independent horizontal scaling across different resource profiles. The platform leverages event-driven asynchronous queues for data ingestion and parallel scatter-gather orchestration for retrieval.

## Component Breakdown

### 1. Stateless API Gateway (Python / FastAPI)

The entry point of the system is a lightweight, fully stateless API tier. It acts as the cluster's orchestration layer, handling external HTTP requests and managing internal gRPC communication.

* **Concurrency:** Utilizes Python 3.11’s `asyncio.TaskGroup` and `grpc.aio` to maintain sub-millisecond network fan-out without blocking the main event loop.
* **Dynamic Discovery:** Resolves active worker nodes on the fly by querying internal cluster DNS, eliminating the need for a centralized service registry.
* **Role:** Enforces API contracts, validates payloads, and merges distributed mathematical results (min-heaps) before returning the final response to the client.

### 2. Stateful Vector Engine Workers (Go / gRPC)

The data layer consists of high-performance, polyglot microservices written in Go. These worker nodes are stateful—they hold the partitioned vector spaces and execute the heavy computational clustering (Lloyd’s Algorithm).

* **Asynchronous Queues:** Ingestion is decoupled from processing. Incoming vectors are dropped into a hardware-like channel FIFO queue, allowing the gRPC response to return immediately while a background worker safely commits the data.
* **Lock-Free Computation:** Heavy indexing calculations are performed on immutable snapshot copies of the data, ensuring that reads (searches) are never blocked by writes (inserts) or re-indexing operations.

## Network & Routing Topography

To ensure even data distribution and rapid query execution, the cluster relies on two distinct routing patterns:

* **Scatter-Gather (Read Path):** When a search query hits the Gateway, the request is broadcasted perfectly in parallel to every active Go worker node. Each node calculates its local nearest neighbor, and the Gateway aggregates these responses to find the absolute global minimum distance.
* **Randomized Load Balancing (Write Path):** Insertions are routed using a randomized node-selection strategy (acting as a foundational step toward Power of Two Choices). This prevents any single worker node from becoming a hot spot during massive data ingestion.

## Scalability & Deployment Strategy

The deployment leverages Kubernetes primitives to treat the infrastructure as immutable, auto-scaling domains:

* **Headless Services:** The Go database workers are exposed via a K8s Headless Service (`ClusterIP: None`). This bypasses standard K8s load balancing, passing the raw Pod IPs directly to the Python Gateway for precise, node-aware gRPC routing.
* **Standard Load Balancers:** The Python Gateways sit behind a standard external LoadBalancer, distributing user traffic in a simple round-robin fashion since the Gateways hold no local state.
* **Dual-Axis Autoscaling:** Two independent Horizontal Pod Autoscalers (HPAs) govern the cluster. A spike in raw web traffic triggers the provisioning of more Python pods, while a spike in computational load or memory utilization scales the Go worker nodes.