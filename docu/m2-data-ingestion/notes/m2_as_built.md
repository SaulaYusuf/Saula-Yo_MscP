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

### 2.3 Implemented Scripts

| Script | Dataset | Records | Endpoint | Chaincode Function |
|--------|---------|---------|----------|---------------------|
| `ingest_sensors.py` | `shipment-sensor-dataset.csv` | 8,000 | `/api/sensor` | `RecordTelemetry` |
| `ingest_logistics.py` | `smart_logistics_dataset.csv` | 1,000 | `/api/logistics` | `RecordHandover` |
| `ingest_metadata.py` | `bdt_mba_supplychain_dataset_2024.csv` | 500 | `/api/metadata` | `RecordMetadata` |

---
---

## 3. Performance Results

| Dataset | Records | Time (s) | Throughput (TPS) |
|---------|---------|----------|------------------|
| Sensor (IoT) | 8,000 | 123.8 | 64.6 |
| Logistics | 1,000 | 8.7 | 114.9 |
| Metadata | 500 | 4.7 | 106.2 |

> **Note:** TPS varies due to different payload sizes and chaincode complexity. The sensor dataset includes threshold evaluation, which adds computational overhead.

### Ingestion Comparison Chart

![Ingestion Comparison](docu/m2-data-ingestion/code/ingestion_comparison.png)

The chart shows that the sensor dataset is the largest and has the lowest TPS, while logistics and metadata have higher throughput due to simpler transactions.

> **Note:** The TPS is limited by the Go bridge’s single‑threaded processing and the Fabric orderer’s block configuration. Future stress tests (M5) will measure the absolute maximum throughput.

---

## 4. Verification

After each ingestion, queries confirmed that data was stored correctly. Example:

```bash
# Sensor twin
peer chaincode query -C mychannel -n slave-twin -c '{"function":"ReadTwin","Args":["sensor-001"]}'
# Logistics handover
peer chaincode query -C mychannel -n slave-twin -c '{"function":"ReadHandover","Args":["Truck_7"]}'
# Metadata
peer chaincode query -C mychannel -n slave-twin -c '{"function":"ReadMetadata","Args":["A0001"]}'
```

---

## 5. Challenges & Resolutions

- **SDK Unavailability:** The Python Fabric SDK does not exist. We pivoted to the two‑tier architecture with Go bridge.
- **API Version Mismatch:** The Fabric Gateway Go SDK changed its API. We used `WithClientConnection` (v1.11.0) instead of the deprecated endpoint functions.
- **TLS Cert Paths:** Ensured the bridge could locate the peer and orderer TLS CA certificates by using relative paths from `gateway-api/`.

---

## 6. Next Steps

- **M5:** Run stress tests, measure TPS, latency, resource usage, and simulate Raft crash recovery.

---