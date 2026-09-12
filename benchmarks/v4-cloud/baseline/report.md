## Performance Load Test Report: VectorEngine Cluster

### Executive Summary

A 5-minute distributed load test was conducted against the **VectorEngine** Kubernetes deployment using Locust. The scenario simulated **500 concurrent users** performing a mixed workload of **75% vector searches** and **25% vector insertions**.

The test successfully validated the cluster's **Horizontal Pod Autoscaling (HPA)** mechanism. Under peak load, the Python API gateway CPU utilization surged to **126%**, triggering the HPA to automatically scale the API deployment from **2 to 4 replicas**. This expansion reduced average CPU load per pod to **65%**, fully stabilizing the service.

---

## 1. Test Overview & Metrics

| Parameter | Value |
| --- | --- |
| **Test Duration** | 4 minutes, 49 seconds (2:59:30 PM – 3:04:19 PM) |
| **Target Host** | `http://localhost:8000` (FastAPI Gateway) |
| **Total Requests** | 28,670 |
| **Total Failures** | 860 (3.00% Overall Failure Rate) |
| **Throughput** | 99.02 req/sec |
| **Workload Ratio** | 75% `POST /search` | 25% `POST /insert` |

---

## 2. Request Breakdown

| Endpoint | Total Requests | Failures | Failure Rate | Throughput (RPS) | Avg Latency | Min Latency | Max Latency |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `POST /insert` | 7,132 | 0 | **0.00%** | 24.63 req/s | 4,588 ms | 85 ms | 21,169 ms |
| `POST /search` | 21,538 | 860 | **3.99%** | 74.39 req/s | 4,971 ms | 127 ms | 21,210 ms |
| **Aggregated** | **28,670** | **860** | **3.00%** | **99.02 req/s** | **4,876 ms** | **85 ms** | **21,210 ms** |

---

## 3. Latency Quantiles

| Endpoint | 50%ile | 70%ile | 90%ile | 95%ile | 99%ile | Max (100%) |
| --- | --- | --- | --- | --- | --- | --- |
| `POST /insert` | 4,300 ms | 5,000 ms | 6,300 ms | 6,800 ms | 9,600 ms | 12,000 ms |
| `POST /search` | 4,700 ms | 5,500 ms | 6,600 ms | 7,100 ms | 10,000 ms | 12,000 ms |
| **Aggregated** | **4,600 ms** | **5,300 ms** | **6,600 ms** | **7,000 ms** | **9,800 ms** | **12,000 ms** |

---

## 4. Autoscaling & System Stabilization Analysis

The cluster demonstrated dynamic recovery and self-healing under traffic pressure:

```
[500 Users Swarm]
       │
       ▼
[vector-api (2 Pods)] ──► CPU Spikes to 126% (Threshold: 70%)
       │
       ▼ (HPA Event Triggered)
[HPA Scales vector-api] ──► Pods increase from 2 ──► 4
       │
       ▼
[Load Distributed Across 4 Pods] ──► CPU Stabilizes at 65%

```

### HPA State Transition

1. **Initial State (2 Replicas):** The API gateway started with 2 pods. When 500 virtual users began sending requests, Python's synchronous/async HTTP processing and JSON serialization bound the CPU, driving usage to **122%–126%**.
2. **Autoscaling Event:** `vector-api-hpa` recognized that CPU utilization exceeded the 70% target and issued a scale command to `vector-api-deployment`, increasing the replica count from **2 to 4 pods**.
3. **Stabilization Phase:** Once the 2 new API pods passed readiness probes and began accepting traffic, request concurrency was evenly distributed. CPU utilization settled at **65%**, stabilizing latency and throughput.
4. **Backend Worker Resilience:** Throughout the entire test, `vector-worker-hpa` stayed at **15% CPU utilization** across its 3 Go nodes. The Go gRPC workers easily handled the search query volume without needing to scale up.

---

## 5. Diagnostic Findings

> **Root Cause of the 860 Failures (`500 Internal Server Error`)**
> All 860 failures occurred on `POST /search` during the initial traffic spike (starting at 2:59:42 PM) while the API gateway was running with only 2 pods.
> Because the metrics server scraping interval and container spin-up time create a ~60–90 second delay before new pods are ready, the initial 2 API pods suffered from **connection queue saturation** and **gRPC channel timeout exhaustion**. Once the 2 new pods became `Running`, zero additional failures were recorded.

---

## 6. Recommendations for System Hardening

1. **Implement Connection Retries in Python:** Wrap the gRPC client calls in `vector-api` with an exponential backoff retry mechanism so transient channel saturation does not surface as HTTP 500 errors.
2. **Adjust HPA Scaling Target:** Lower the target CPU utilization threshold on `vector-api-hpa` from `70%` to `50%` or set `minReplicas: 3` to give the gateway more headroom before traffic spikes.
3. **Configure Connection Pooling:** Ensure HTTP and gRPC connection pools on the FastAPI layer are configured to handle concurrent requests up to container resource limits.