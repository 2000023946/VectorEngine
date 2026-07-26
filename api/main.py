import asyncio
import socket
import random
from contextlib import asynccontextmanager
from typing import List, Dict
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel

# Import standard gRPC and the AIO (async) module
import grpc
from grpc import aio

import vectordb_pb2
import vectordb_pb2_grpc

# The internal Kubernetes DNS for your headless service
HEADLESS_DNS = "vector-worker-headless.default.svc.cluster.local"
GRPC_PORT = 50051

# Global dictionary to hold dynamic gRPC channels mapped by IP
CHANNEL_POOL: Dict[str, aio.Channel] = {}


@asynccontextmanager
async def lifespan(app: FastAPI):
    yield
    # Shutdown: Close all cached dynamic channels gracefully
    for ip, channel in CHANNEL_POOL.items():
        await channel.close()
        print(f"Closed channel to {ip}")


app = FastAPI(title="Distributed Vector DB Gateway", lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class VectorInput(BaseModel):
    vector: List[float]


def get_active_worker_ips() -> List[str]:
    """Resolves K8s Headless Service DNS to get all active Go Pod IPs."""
    try:
        _, _, ips = socket.gethostbyname_ex(HEADLESS_DNS)
        return ips
    except socket.gaierror:
        # Fallback for local testing without K8s
        return ["localhost"]


def get_grpc_stub(ip: str) -> vectordb_pb2_grpc.VectorServiceStub:
    """Retrieves or creates an async gRPC stub for a specific IP."""
    target = f"{ip}:{GRPC_PORT}"
    if target not in CHANNEL_POOL:
        # Create an asynchronous channel
        channel = aio.insecure_channel(target)
        CHANNEL_POOL[target] = channel
    return vectordb_pb2_grpc.VectorServiceStub(CHANNEL_POOL[target])


@app.post("/insert")
async def insert_vector(data: VectorInput):
    worker_ips = get_active_worker_ips()
    if not worker_ips:
        raise HTTPException(status_code=503, detail="No active database nodes found.")

    # Load Balancing: Pick a random node for the insert
    # (If you add a GetCount RPC later, you can upgrade this to Power of Two Choices)
    target_ip = random.choice(worker_ips)
    stub = get_grpc_stub(target_ip)
    
    try:
        req = vectordb_pb2.InsertRequest(vector=data.vector)
        # MUST use 'await' with aio stubs
        resp = await stub.Insert(req)
        return {"success": resp.success, "node_inserted": target_ip}
    except grpc.RpcError as e:
        raise HTTPException(status_code=500, detail=f"gRPC Error: {e.details()}")


async def search_single_node(ip: str, vector: List[float]):
    """Helper to search a single node asynchronously."""
    stub = get_grpc_stub(ip)
    req = vectordb_pb2.SearchRequest(query=vector)
    try:
        # 500ms timeout so one frozen pod doesn't break the whole cluster
        resp = await stub.Search(req, timeout=0.5)
        return resp
    except grpc.RpcError:
        return None


@app.post("/search")
async def search_vector(data: VectorInput):
    worker_ips = get_active_worker_ips()
    if not worker_ips:
        raise HTTPException(status_code=503, detail="No active database nodes found.")

    # 1. SCATTER: Query all active nodes perfectly in parallel
    async with asyncio.TaskGroup() as tg:
        tasks = [
            tg.create_task(search_single_node(ip, data.vector))
            for ip in worker_ips
        ]

    # 2. GATHER: Collect successful responses, dropping any nodes that timed out
    valid_responses = [task.result() for task in tasks if task.result() is not None]

    if not valid_responses:
        raise HTTPException(status_code=500, detail="All nodes failed to process the search.")

    # 3. COMBINE: Find the absolute best distance across all responding nodes
    best_match = min(valid_responses, key=lambda resp: resp.distance)

    return {
        "nearest_vector": list(best_match.nearest_vector),
        "distance": best_match.distance,
        "nodes_searched": len(worker_ips)
    }