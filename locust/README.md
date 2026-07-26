Since your virtual environment `(venv)` is already active, here is the minimal workflow to start your load test:

1. **Install dependencies** (if you haven't already):
```bash
pip install -r requirements.txt

```


2. **Start the Locust server:**
*Because your file is named `locust.py` instead of the default `locustfile.py`, you must use the `-f` flag.*
```bash
locust -f locust.py

```


3. **Launch the attack:**
Open `http://localhost:8089` in your browser, set your user count/spawn rate, and point the host to `