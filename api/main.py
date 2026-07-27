import asyncio
import socket
import random
import itertools
from contextlib import asynccontextmanager
from typing import List, Dict
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import grpc
from grpc import aio

import vectordb_pb2
import vectordb_pb2_grpc

HEADLESS_DNS = "vector-worker-headless.default.svc.cluster.local"
GRPC_PORT = 50051

# --- MULTIPLEXING CONFIGURATION ---
# 5 channels per node ensures we don't hit the HTTP/2 max_concurrent_streams ceiling
POOL_SIZE_PER_NODE = 5  

# Global State
CHANNEL_POOL: Dict[str, List[aio.Channel]] = {}
POOL_ITERATORS: Dict[str, itertools.cycle] = {}
WORKER_STATS: Dict[str, int] = {} 
ACTIVE_IPS: List[str] = []

# gRPC C-Core tuning for high-throughput and keep-alive stability
GRPC_OPTIONS = [
    ('grpc.max_concurrent_streams', 1000),
    ('grpc.keepalive_time_ms', 10000),
    ('grpc.keepalive_timeout_ms', 5000),
    ('grpc.keepalive_permit_without_calls', True),
    ('grpc.http2.max_pings_without_data', 0),
]

async def sync_worker_stats():
    """Background task to resolve DNS and fetch sizes every 5 seconds."""
    global WORKER_STATS, ACTIVE_IPS
    loop = asyncio.get_running_loop()
    
    while True:
        try:
            _, _, ips = await loop.run_in_executor(None, socket.gethostbyname_ex, HEADLESS_DNS)
            ACTIVE_IPS = ips
        except socket.gaierror:
            ACTIVE_IPS = ["localhost"]

        new_stats = {}
        for ip in ACTIVE_IPS:
            try:
                stub = get_grpc_stub(ip)
                req = vectordb_pb2.StatsRequest()
                resp = await stub.GetStats(req, timeout=1.0)
                new_stats[ip] = resp.count
            except grpc.RpcError:
                pass
                
        if new_stats:
            WORKER_STATS = new_stats
            
        await asyncio.sleep(5)


def get_grpc_stub(ip: str) -> vectordb_pb2_grpc.VectorServiceStub:
    """Fetches a stub using a Round-Robin channel pool to prevent stream starvation."""
    target = f"{ip}:{GRPC_PORT}"
    
    if target not in CHANNEL_POOL:
        # Initialize a pool of channels for this newly discovered node
        CHANNEL_POOL[target] = [
            aio.insecure_channel(target, options=GRPC_OPTIONS) 
            for _ in range(POOL_SIZE_PER_NODE)
        ]
        # Create an infinite iterator that cycles through indices 0 to POOL_SIZE-1
        POOL_ITERATORS[target] = itertools.cycle(range(POOL_SIZE_PER_NODE))
        
    # Pick the next channel in the rotation
    idx = next(POOL_ITERATORS[target])
    channel = CHANNEL_POOL[target][idx]
    
    return vectordb_pb2_grpc.VectorServiceStub(channel)


@asynccontextmanager
async def lifespan(app: FastAPI):
    stats_task = asyncio.create_task(sync_worker_stats())
    yield
    stats_task.cancel()
    # Clean up all channels across all pools
    for pool in CHANNEL_POOL.values():
        for channel in pool:
            await channel.close()

app = FastAPI(title="Distributed Vector DB Gateway", lifespan=lifespan)
app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_credentials=True, allow_methods=["*"], allow_headers=["*"])

class VectorInput(BaseModel):
    vector: List[float]


def pick_insert_target() -> str:
    if not ACTIVE_IPS:
        return ""

    current_stats = {ip: WORKER_STATS.get(ip, 0) for ip in ACTIVE_IPS}
    targets = list(current_stats.keys())
    
    weights = [1.0 / max(count, 0.1) for count in current_stats.values()]
    total_weight = sum(weights)
    probabilities = [w / total_weight for w in weights]
    
    return random.choices(targets, weights=probabilities, k=1)[0]


@app.post("/insert")
async def insert_vector(data: VectorInput):
    target_ip = pick_insert_target()
    if not target_ip:
        raise HTTPException(status_code=503, detail="No active database nodes found.")

    stub = get_grpc_stub(target_ip)
    try:
        req = vectordb_pb2.InsertRequest(vector=data.vector)
        resp = await stub.Insert(req)
        return {"success": resp.success, "node_inserted": target_ip}
    except grpc.RpcError as e:
        raise HTTPException(status_code=500, detail=f"gRPC Error: {e.details()}")

async def search_single_node(ip: str, vector: List[float]):
    stub = get_grpc_stub(ip)
    req = vectordb_pb2.SearchRequest(query=vector)
    try:
        # Bumped timeout to 2.0s. Under extreme Locust load, Python's 
        # event loop queues tasks for ~500ms before they even hit the network.
        resp = await stub.Search(req, timeout=2.0)
        return resp
    except grpc.RpcError:
        return None


@app.post("/search")
async def search_vector(data: VectorInput):
    if not ACTIVE_IPS:
        raise HTTPException(status_code=503, detail="No active database nodes found.")

    # Target all active pods to find the global nearest neighbor
    targets = ACTIVE_IPS

    async with asyncio.TaskGroup() as tg:
        tasks = [tg.create_task(search_single_node(ip, data.vector)) for ip in targets]

    valid_responses = [task.result() for task in tasks if task.result() is not None]

    if not valid_responses:
        raise HTTPException(status_code=500, detail="All nodes failed to process the search or timed out.")

    best_match = min(valid_responses, key=lambda resp: resp.distance)

    return {
        "nearest_vector": list(best_match.nearest_vector),
        "distance": best_match.distance,
        "nodes_searched": len(targets)
    }