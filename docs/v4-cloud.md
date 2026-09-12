# Cloud Architecture

VectorEngine runs as a distributed service on Kubernetes, separating the API gateway from the Go vector database workers.

```text
                    Locust
                 Load Generator
                      |
                     HTTP
                      v
             +------------------+
             | FastAPI Gateway  |
             |   vector-api     |
             +--------+---------+
                      |
                 gRPC / Protobuf
                      |
                      v
             +------------------+
             | Kubernetes       |
             | Headless Service |
             +--------+---------+
                      |
          +-----------+-----------+
          |           |           |
          v           v           v
       Go Worker   Go Worker   Go Worker
       Vector DB   Vector DB   Vector DB
          |           |           |
          +-----------+-----------+
                      |
                 Prometheus
                      |
                     KEDA
                      |
          Request-Driven Scaling
                 
```

## Components

* **FastAPI Gateway** — exposes the HTTP API and forwards requests to workers through gRPC.
* **Go Workers** — run the VectorEngine database and handle search/insert operations.
* **Kubernetes Headless Service** — provides service discovery for the Go worker pods.
* **Prometheus + KEDA** — provide **request-driven autoscaling** for the Go workers based on incoming gRPC workload.
* **API HPA** — scales the FastAPI gateway replicas based on CPU utilization.
* **Locust** — generates concurrent search/insert workloads for performance testing.

## Request Flow

```text
HTTP Request
     ↓
FastAPI
     ↓
gRPC / Protobuf
     ↓
Go VectorEngine Worker
     ↓
Response
```

The system does not use a separate distributed message queue between the API and database workers. Requests are handled directly through gRPC.

The cloud infrastructure is separated from the VectorEngine implementation, allowing the underlying database engine to be replaced without redesigning the deployment architecture.
