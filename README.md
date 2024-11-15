# Real-time Data Stream API

## Overview
This project is a Real-Time Data Streaming API developed in Golang that streams data to and from Redpanda (Kafka) with Prometheus for monitoring and Grafana for visualization. It includes rate-limiting, unit and integration tests, and performance benchmarking.

## Features
- Real-time data streaming using Golang and Redpanda (Kafka)
- Prometheus integration for metrics
- Rate limiting middleware
- Unit and integration tests for the API
- Performance benchmarking with detailed throughput and latency data

## Prerequisites
- **Go** (1.20+)
- **Redpanda** (Kafka)
- **Prometheus**
- **Grafana**
- **GitHub CLI** (`gh`) (for repository management)
- **wrk** (for load testing)

## Installation and Setup

### Clone the Repository
```bash
git clone https://github.com/your-username/Real-time-data-stream-API.git
cd Real-time-data-stream-API
```

### Setting up Redpanda (Kafka)

1. **Download Redpanda** (example for macOS):
   ```bash
   brew install redpanda
   ```

2. **Start Redpanda**:
   ```bash
   redpanda start --overprovisioned --smp 1 --memory 1G --reserve-memory 0M --node-id 0 --check=false
   ```

3. **Create the Topic**:
   ```bash
   rpk topic create stream-topic
   ```

### Setting up Prometheus

1. **Install Prometheus** (example for macOS):
   ```bash
   brew install prometheus
   ```

2. **Configure Prometheus** by adding the following to `prometheus.yml`:
   ```yaml
   global:
     scrape_interval: 15s

   scrape_configs:
     - job_name: 'golang-app'
       static_configs:
         - targets: ['localhost:8080']
   ```

3. **Start Prometheus**:
   ```bash
   prometheus --config.file=prometheus.yml
   ```

### Setting up Grafana

1. **Install Grafana** (example for macOS):
   ```bash
   brew install grafana
   ```

2. **Start Grafana**:
   ```bash
   brew services start grafana
   ```

3. **Add Data Source in Grafana**:
   - Go to `http://localhost:3000`, log in (default is admin/admin)
   - Add Prometheus as a data source (URL: `http://localhost:9090`).

4. **Import Dashboard**:
   - Create a dashboard or import JSON to monitor metrics like `http_requests_total` and `go_gc_duration_seconds`.

## Running the Application

1. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

2. **Run the Application**:
   ```bash
   go run main.go
   ```
   The server should start on port `8080`.

## Testing the Application

### Unit Tests
Run unit tests to ensure the functionality of individual components.
```bash
go test -v ./tests/unit_tests.go
```

### Integration Tests
Run integration tests to test the end-to-end flow.
```bash
go test -v ./tests/integration_tests.go
```

### Performance Benchmarking
Use `wrk` to load test the `/stream/start` endpoint.
```bash
wrk -t10 -c1000 -d30s -H "X-API-Key: your-secure-api-key" http://localhost:8080/stream/start
```



## API Endpoints

### Start a Stream
```bash
curl -X POST http://localhost:8080/stream/start -H "X-API-Key: your-secure-api-key"
```

### Send Data to a Stream
Replace `<stream_id>` with the ID returned from the `/stream/start` endpoint.
```bash
curl -X POST http://localhost:8080/stream/<stream_id>/send -d "sample data" -H "X-API-Key: your-secure-api-key"
```

### Get Results via WebSocket
Connect to the results using `wscat`.
```bash
wscat -c ws://localhost:8080/stream/<stream_id>/results
```

### Metrics Endpoint
View Prometheus metrics:
```bash
curl http://localhost:8080/metrics
```

## Additional Setup for GitHub Actions CI/CD

1. **Add GitHub Actions Workflow**: Create a `.github/workflows/go.yml` file with:
   ```yaml
   name: Go

   on:
     push:
       branches: [ "main" ]
     pull_request:
       branches: [ "main" ]

   jobs:
     build:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v2
         - name: Set up Go
           uses: actions/setup-go@v2
           with:
             go-version: '^1.20'
         - name: Build
           run: go build -v ./...
         - name: Test
           run: go test -v ./...
   ```

2. **Push Changes**:
   ```bash
   git add .
   git commit -m "Initial commit with README and GitHub Actions setup"
   git push origin main
   ```

---

This README provides the complete setup, testing, and usage instructions, ensuring the application meets all criteria. Adjust paths, usernames, or project-specific details as needed.
