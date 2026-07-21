# M2: Data Ingestion Engine – Design & Implementation Plan

**Date:** 2026‑07‑06 (Last updated: 2026‑07‑21)  
**Student:** Saula Yusuf Owolabi  
**Milestone:** M2 – Python Asynchronous Client Pipeline

---

## 1. Objective

Build an asynchronous Python client that:
- Reads the three Kaggle CSV datasets (Smart Logistics, Cold‑Chain Sensor, BDT‑MBA Digital Twin State).
- Transforms each row into a JSON transaction proposal.
- Streams the proposals to the **Go API Bridge**, which submits them to the Fabric network via the official Gateway SDK.
- Routes sensor telemetry to `slave-twin.RecordTelemetry`.
- Routes logistics milestones to `slave-twin.RecordHandover`.
- Routes asset metadata to `slave-twin.RecordMetadata`.

This design avoids the unmaintained Python Fabric SDK and uses a production‑grade two‑tier microservice architecture.

---

## 2. Dataset Mapping

| Dataset | Source | Target Chaincode Function | Key Fields |
|---------|--------|---------------------------|------------|
| Smart Logistics Supply Chain | `smart_logistics_dataset.csv` | `RecordHandover` | `Asset_ID`, `Shipment_Status`, `Timestamp`, `Latitude`, `Longitude` |
| Cold‑Chain Silent Failure | `shipment-sensor-dataset.csv` | `RecordTelemetry` | `sensor_id`, `temp_max_c`, `temp_min_c`, `rh_mean`, `timestamp` |
| BDT‑MBA Supply Chain | `bdt_mba_supplychain_dataset_2024.csv` | `RecordMetadata` | `Asset_ID`, `Location`, `Temperature`, `Condition_Score`, `SupplyChain_Efficiency_Label`, etc. |

---

## 3. Technical Architecture (Two‑Tier)

### 3.1 Python Client (Edge Layer)

- **Language:** Python 3.12+
- **Libraries:** `aiohttp` (async HTTP), `csv` (parsing), `asyncio`
- **Concurrency:** Configurable worker pool (default 10) for parallel POST requests.
- **Endpoints:**
  - `POST /api/sensor` → sensor data
  - `POST /api/logistics` → logistics data
  - `POST /api/metadata` → asset metadata
- **Payload:** JSON object matching the chaincode function arguments.

The script reads the CSV, queues rows, and workers send HTTP requests to the Go bridge. It logs successes, failures, and timing.

### 3.2 Go API Bridge (Gateway Layer)

- **Framework:** Go with official `fabric-gateway` SDK v1.11.0
- **Identity:** Org1 Admin (cryptogen‑generated, TLS enabled)
- **TLS:** Uses the peer’s TLS CA certificate (`peer0.org1.example.com/tls/ca.crt`) for secure gRPC.
- **Endpoints:** Exposes `/api/sensor`, `/api/logistics`, `/api/metadata`; each invokes the corresponding chaincode function.
- **Error Handling:** Returns HTTP 500 on chaincode failure; 200 on success.

This bridge is stateless and can be scaled horizontally if needed.

---

## 4. Implementation Steps (Completed)

1. ✅ Set up Python virtual environment with `aiohttp` and `csv`.
2. ✅ Write `ingest_sensors.py` for the sensor dataset.
3. ✅ Write `ingest_logistics.py` for the logistics dataset.
4. ✅ Write `ingest_metadata.py` for the metadata dataset.
5. ✅ Extend Go bridge to support all three endpoints.
6. ✅ Test each ingestion script with 100% success.
7. ✅ Generate ingestion comparison chart.

---

## 5. Challenges & Resolutions

- **SDK Unavailability:** Fabric’s Python SDK is deprecated; adopted the two‑tier microservice pattern.
- **API Version Mismatch:** The Go SDK’s `client.Connect` changed; we used `WithClientConnection` with gRPC connections.
- **TLS Cert Paths:** Initially used wrong CA files; corrected to peer-specific TLS CA certificates.
- **Chaincode Deployment:** Multiple attempts due to packaging path and sequence mismatches; eventually deployed version 2.7 with all functions.

---

## 6. Performance Targets for M5

- **Maximum TPS:** Determine saturation point by varying concurrency (10, 50, 100, 200 workers).
- **Latency Percentiles:** Measure p50, p95, p99 commit latency.
- **Resource Usage:** CPU and memory consumption under load.
- **Crash Recovery:** Time for Raft leader re‑election after orderer failure.