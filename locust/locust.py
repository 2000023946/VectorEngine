import random
from locust import HttpUser, task, between

# Helper function to generate a random 4D vector
def generate_random_vector(dim=4):
    return [round(random.uniform(0.0, 1.0), 4) for _ in range(dim)]

class VectorEngineUser(HttpUser):
    # Short wait time to generate heavy load quickly
    wait_time = between(0.01, 0.1)

    # @task(1) means this runs 25% of the time
    @task(1)
    def insert_vector(self):
        # Adjust the endpoint "/insert" and JSON structure to match your exact FastAPI schema
        vector = generate_random_vector()
        self.client.post("/insert", json={
            "id": str(random.randint(1, 1000000)), # Optional: random ID if your DB requires one
            "vector": vector
        })

    # @task(3) means this runs 75% of the time (3 times as often as inserts)
    @task(3)
    def search_vector(self):
        vector = generate_random_vector()
        self.client.post("/search", json={
            "vector": vector,
            "k": 5
        })