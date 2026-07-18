import requests
import random
import time

BASE_URL = "http://localhost:8000"
DIMENSIONS = 8 

def generate_vector():
    # Generates a random 8-dimensional vector mimicking sensor floats
    return [random.uniform(-10.0, 10.0) for _ in range(DIMENSIONS)]

def test_api():
    print("1. Testing Vector Ingestion (/insert)")
    print("-" * 40)
    
    # Insert 10 vectors to populate the Go engine
    for i in range(1, 11):
        vec = generate_vector()
        response = requests.post(f"{BASE_URL}/insert", json={"vector": vec})
        
        if response.status_code == 200:
            print(f"Inserted vector {i}/10 successfully.")
        else:
            print(f"Failed to insert: {response.text}")
            return

    print("\n2. Testing Nearest Neighbor Search (/search)")
    print("-" * 40)
    
    query_vec = generate_vector()
    print(f"Querying: {[round(v, 4) for v in query_vec]}")
    
    # Measure the round-trip time including JSON serialization, HTTP, and gRPC
    start_time = time.perf_counter()
    response = requests.post(f"{BASE_URL}/search", json={"vector": query_vec})
    end_time = time.perf_counter()
    
    if response.status_code == 200:
        data = response.json()
        nearest = data.get("nearest_vector", [])
        dist = data.get("distance", -1)
        
        latency_ms = (end_time - start_time) * 1000
        
        print(f"\nResult Nearest: {[round(v, 4) for v in nearest]}")
        print(f"Distance (L2 Sq): {dist:.4f}")
        print(f"End-to-End Latency: {latency_ms:.2f} ms")
    else:
        print(f"Search failed: {response.text}")

if __name__ == "__main__":
    try:
        test_api()
    except requests.exceptions.ConnectionError:
        print(f"Error: Could not connect to {BASE_URL}. Make sure uvicorn is running!")