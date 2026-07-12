# M2: Data Ingestion Engine – As‑Built Report

**Date:** 2026‑07‑12  
**Student:** Saula Yusuf Owolabi  
**Milestone:** M2 – Python Asynchronous Ingestion Pipeline

---

## 1. Overview

The ingestion engine streams real‑world Kaggle datasets into the Hyperledger Fabric network. It uses a **two‑tier microservice architecture**:

1. **Python CSV Parser** – reads the `shipment-sensor-dataset.csv` (8,000 records), transforms each row to JSON, and sends HTTP POST requests to the Go API bridge.
2. **Go API Bridge** – uses the official Fabric Gateway SDK to submit transactions to the `slave-twin` chaincode on `mychannel`.

This design avoids the dead `fabric-sdk-py` library and mirrors real‑world edge‑to‑blockchain patterns.

---

## 2. Implementation Details

### 2.1 Python Client (`ingest_sensors.py`)

- **Language:** Python 3.12+
- **Libraries:** `aiohttp` (async HTTP), `csv` (parsing)
- **Concurrency:** 10 concurrent workers
- **Endpoint:** `http://localhost:8080/api/sensor`
- **Payload:** JSON with `sensor_id`, `temp_c`, `humidity`, `timestamp`

The script reads the CSV row‑by‑row, places each row into an async queue, and workers send POST requests concurrently. It logs successes and failures.

### 2.2 Go API Bridge (`gateway-api/main.go`)

- **Framework:** Go with `fabric-gateway` SDK v1.11.0
- **Identity:** Org1 Admin (cryptogen‑generated)
- **TLS:** Enabled; uses the peer and orderer TLS CA certificates from the `test-network`.
- **Endpoint:** `POST /api/sensor` → invokes `slave-twin.RecordTelemetry`

The bridge creates a gRPC connection to the peer, submits the transaction, and returns `{"status":"committed"}` on success.

---

## 3. Performance Results

| Metric | Value |
|--------|-------|
| Total records | 8,000 |
| Successful commits | 8,000 |
| Failed commits | 0 |
| Total time | 123.8 seconds |
| Throughput (TPS) | ~65 TPS |
| Average per‑record latency | ~15.5 ms (excluding network overhead) |

> **Note:** The TPS is limited by the Go bridge’s single‑threaded processing and the Fabric orderer’s block configuration. Future stress tests (M5) will measure the absolute maximum throughput.

---

## 4. Verification

After ingestion, a query confirmed that the twin for `sensor-001` was stored correctly:

```bash
peer chaincode query -C mychannel -n slave-twin -c '{"function":"ReadTwin","Args":["sensor-001"]}'
```

Output:

```json
{"sensor_id":"sensor-001","temp_c":7.5,"humidity":45,"status":"NORMAL","timestamp":"2026-07-06T12:00:00Z"}
```

---

## 5. Challenges & Resolutions

- **SDK Unavailability:** The Python Fabric SDK does not exist. We pivoted to the two‑tier architecture with Go bridge.
- **API Version Mismatch:** The Fabric Gateway Go SDK changed its API. We used `WithClientConnection` (v1.11.0) instead of the deprecated endpoint functions.
- **TLS Cert Paths:** Ensured the bridge could locate the peer and orderer TLS CA certificates by using relative paths from `gateway-api/`.

---

## 6. Next Steps

- **M4:** Develop the `master-logistics` chaincode to handle the logistics dataset (ownership handovers, port arrivals).
- **M5:** Run stress tests, measure TPS, latency, resource usage, and simulate Raft crash recovery.

---

## 7. Files

- Python script: `docu/m2-data-ingestion/code/ingest_sensors.py`
- Go bridge: `gateway-api/main.go`
- Dataset: `data/raw/shipment-sensor-dataset.csv`