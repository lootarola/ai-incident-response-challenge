# AI Incident Response Challenge

A hands-on workshop where you use AI tools and live observability data to find and fix bugs in a production-like Go service.

## How it works

The repository ships with some deliberately introduced bugs. Each bug produces a distinct signal in Grafana when its simulation script is run. Your job is to:

1. Run a simulation script
2. Open Grafana and find the anomaly
3. Use an AI tool (with the Grafana MCP) to investigate the signal and identify the root cause
4. Reproduce the failure locally
5. Apply a fix and confirm the signal disappears

## Prerequisites

- Docker and Docker Compose
- `curl` and `bash` (for simulation scripts)
- Claude with the Grafana MCP configured (for the AI-assisted investigation)

## Setup

```bash
# 1. Clone the repository and enter it
git clone <repo-url>
cd ai-incident-response-challenge

# 2. Copy the environment file
cp .env.example .env

# 3. Start the full stack
docker compose up --build -d

# 4. Wait ~30 seconds for all services to be healthy, then verify
curl http://localhost:8080/healthz
```

A `traffic-generator` container starts automatically and sends a low-rate mix of requests to every endpoint, so the Grafana dashboard is never empty. Your bug's anomaly appears as a clear deviation from this baseline rather than a spike from zero. To pause it for a clean-room run:

```bash
docker compose stop traffic-generator
# ... run your script ...
docker compose start traffic-generator
```

Grafana will be available at **http://localhost:3000** — no login required.
The pre-built dashboard is at **Dashboards → Incident Response — Workshop**.

## Running a simulation

```bash
bash scripts/locust.sh
bash scripts/moth.sh
bash scripts/aphid.sh
bash scripts/slug.sh
bash scripts/tick.sh
```

You can override the base URL if the API is running elsewhere:

```bash
BASE_URL=http://my-host:8080 bash scripts/locust.sh
```

## Investigating with AI

Open an AI tool of your choice and point it at your Grafana instance using the Grafana MCP. A good starting prompt:

> "I am running a Go e-commerce API and something is wrong. Here is my Grafana dashboard. Can you investigate the signals and tell me what you think is happening?"

Let the telemetry guide the conversation. Claude can query Prometheus metrics, read Loki logs, and inspect Tempo traces, the root cause is visible in all three.

## Verifying your fix

After applying a fix:

```bash
# Rebuild and restart the API container only
docker compose up --build -d api

# Re-run your assigned script
bash scripts/<your-script>.sh

# Confirm the anomalous signal is gone in Grafana
```

You can also run the test suite to confirm you haven't broken anything else:

```bash
go test ./internal/...
```

## Tearing down

```bash
docker compose down -v
```

The `-v` flag removes the MongoDB and observability volumes so the next run starts clean.

## Stack

| Service | URL |
|---------|-----|
| API | http://localhost:8080 |
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |
| Loki | http://localhost:3100 |
| Tempo | http://localhost:3200 |
| MongoDB | localhost:27017 |
