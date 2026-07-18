from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List
import grpc

# Import the generated gRPC code
import vectordb_pb2
import vectordb_pb2_grpc


from fastapi.middleware.cors import CORSMiddleware


# Global dictionary to hold our gRPC stub
grpc_client = {}

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup: Connect to the Go gRPC server
    channel = grpc.insecure_channel('localhost:50051')
    grpc_client["stub"] = vectordb_pb2_grpc.VectorServiceStub(channel)
    print("Connected to Go Vector Engine on :50051")
    yield
    # Shutdown: Close the channel gracefully
    channel.close()

app = FastAPI(title="Vector DB API Gateway", lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"], # Or specify ["http://localhost:5173"]
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Pydantic model ensures incoming JSON has the correct shape
class VectorInput(BaseModel):
    vector: List[float]

@app.post("/insert")
def insert_vector(data: VectorInput):
    try:
        req = vectordb_pb2.InsertRequest(vector=data.vector)
        resp = grpc_client["stub"].Insert(req)
        return {"success": resp.success}
    except grpc.RpcError as e:
        raise HTTPException(status_code=500, detail=f"gRPC Error: {e.details()}")

@app.post("/search")
def search_vector(data: VectorInput):
    try:
        req = vectordb_pb2.SearchRequest(query=data.vector)
        resp = grpc_client["stub"].Search(req)
        return {
            "nearest_vector": list(resp.nearest_vector),
            "distance": resp.distance
        }
    except grpc.RpcError as e:
        raise HTTPException(status_code=500, detail=f"gRPC Error: {e.details()}")