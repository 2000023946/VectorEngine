import asyncio
import socket
import random
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

CHANNEL_POOL: Dict[str, aio.Channel] = {}
WORKER_STATS: Dict[str, int] = {} 
ACTIVE_IPS: List[str] = []  # NEW: Global cache for worker IPs


async def sync_worker_stats():
    """Background task to resolve DNS and fetch sizes every 5 seconds."""
    global WORKER_STATS, ACTIVE_IPS
    loop = asyncio.get_running_loop()
    
    while True:
        # 1. Resolve DNS asynchronously so we don't block the main thread
        try:
            _, _, ips = await loop.run_in_executor(None, socket.gethostbyname_ex, HEADLESS_DNS)
            ACTIVE_IPS = ips
        except socket.gaierror:
            ACTIVE_IPS = ["localhost"]

        # 2. Fetch stats from the discovered IPs
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
    target = f"{ip}:{GRPC_PORT}"
    if target not in CHANNEL_POOL:
        channel = aio.insecure_channel(target)
        CHANNEL_POOL[target] = channel
    return vectordb_pb2_grpc.VectorServiceStub(CHANNEL_POOL[target])


@asynccontextmanager
async def lifespan(app: FastAPI):
    stats_task = asyncio.create_task(sync_worker_stats())
    yield
    stats_task.cancel()
    for ip, channel in CHANNEL_POOL.items():
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
        # Fail fast under load so we don't back up the gateway queue
        resp = await stub.Search(req, timeout=1.0)
        return resp
    except grpc.RpcError:
        return None


@app.post("/search")
async def search_vector(data: VectorInput):
    if not ACTIVE_IPS:
        raise HTTPException(status_code=503, detail="No active database nodes found.")

    # Optional: If you have many worker pods, pick a subset (e.g., up to 2) 
    # instead of broadcasting to every single node simultaneously.
    # targets = random.sample(ACTIVE_IPS, min(len(ACTIVE_IPS), 2))
    targets = ACTIVE_IPS

    async with asyncio.TaskGroup() as tg:
        tasks = [tg.create_task(search_single_node(ip, data.vector)) for ip in targets]

    valid_responses = [task.result() for task in tasks if task.result() is not None]

    if not valid_responses:
        raise HTTPException(status_code=500, detail="All nodes failed to process the search.")

    best_match = min(valid_responses, key=lambda resp: resp.distance)

    return {
        "nearest_vector": list(best_match.nearest_vector),
        "distance": best_match.distance,
        "nodes_searched": len(targets)
    }