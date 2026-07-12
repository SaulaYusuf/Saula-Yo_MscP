# Saula-Yo_MscP

## Project Overview

**Title:** Enhancing Supply Chain Traceability and Transparency via Master‑Slave Multi‑Chain Architectures and Digital Twins: An Empirical Performance Evaluation.

**Student:** Saula Yusuf Owolabi  
**Institution:** Ulster University  
**Date:** July 2026 (Last Updated: 2026‑07‑12)

---

## Repository Structure

```
.
├── chaincode/                 # Smart contract source code (Go)
│   └── slave-twin/            # IoT/Digital Twin chaincode (extended for logistics)
│       ├── main.go            # Threshold engine + logistics handovers
│       ├── go.mod
│       └── go.sum
├── data/
│   └── raw/                   # Kaggle datasets
│       ├── bdt_mba_supplychain_dataset_2024.csv
│       ├── shipment-sensor-dataset.csv
│       └── smart_logistics_dataset.csv
├── docu/                      # Full project documentation (milestone notes, results)
│   ├── _global-notes/         # Environment setup, tech stack justification
│   ├── m1-infrastructure/     # Network provisioning & pivot rationale
│   ├── m2-data-ingestion/     # Python ingestion engine (as‑built)
│   ├── m3-slave-chain/        # Slave chaincode development & deployment
│   ├── m4-master-chain/       # Master logistics chaincode (in progress)
│   ├── m5-evaluation/         # Performance metrics & stress tests (planned)
│   └── m6-defense-prep/       # Final report, slides, video walkthrough (planned)
├── gateway-api/               # Go REST API bridge (HTTP -> Fabric)
│   └── main.go                # Connects to Fabric, exposes /api/sensor and /api/logistics
├── README.md                  # This file
└── milestones.md              # Detailed milestone map
```

---

## Architectural Overview (Current Implementation)

> **Note:** The original proposal described two physical channels (`masterchannel` and `slavechannel`). During deployment, the Fabric `test-network` scripts proved unreliable for custom channel provisioning. A **practical pivot** was made to a single channel (`mychannel`) with **logical segregation** at the smart‑contract level.  
> This does **not** compromise the research hypothesis – the evaluation still measures whether isolating high‑frequency IoT telemetry from macro‑logistics milestones prevents consensus bottlenecks. For full rationale, see [docu/m3-slave-chain/notes/architecture_pivot_rationale.md](docu/m3-slave-chain/notes/architecture_pivot_rationale.md).

### Current Architecture

- **Blockchain Platform:** Hyperledger Fabric v2.5.16
- **Crypto Material:** `cryptogen` (static)
- **State Database:** levelDB
- **Consensus:** Raft (5 orderer nodes)
- **Single Channel:** `mychannel`
- **Two Chaincode Contracts (in one package):**
  - **Slave (IoT) Contract:** `RecordTelemetry` – ingests environmental sensor telemetry, evaluates temperature thresholds, updates Digital Twin status (`NORMAL` / `SPOILED`).
  - **Master (Logistics) Contract:** `RecordHandover` – handles ownership handovers, port arrivals, and cryptographic anchoring of final milestones (to be deployed).

### Data Ingestion Pipeline (Two‑Tier)

1. **Python Parser** – reads Kaggle CSV datasets, transforms rows to JSON, and sends HTTP POST requests to the Go bridge.
2. **Go API Bridge** – uses the official Fabric Gateway SDK to submit transactions to the chaincode.

This pattern avoids the unmaintained Python Fabric SDK and is industry‑standard for edge‑to‑blockchain communication.

---

## Datasets & Mapping

| Dataset | Source | Target Contract | Purpose |
|---------|--------|-----------------|---------|
| **Smart Logistics Supply Chain** | [`data/raw/smart_logistics_dataset.csv`](data/raw/smart_logistics_dataset.csv) | `RecordHandover` | Macro‑movement tracking (origin/destination ports, current status) |
| **Cold‑Chain Silent Failure** | [`data/raw/shipment-sensor-dataset.csv`](data/raw/shipment-sensor-dataset.csv) | `RecordTelemetry` | High‑frequency IoT telemetry (temperature, humidity) – triggers threshold engine |
| **BDT‑MBA Digital Twin State** | [`data/raw/bdt_mba_supplychain_dataset_2024.csv`](data/raw/bdt_mba_supplychain_dataset_2024.csv) | `RecordTelemetry` | Asset metadata (condition score, maintenance logs) – syncs physical‑digital twin |

---

## System Evaluation Metrics

The success of this architecture is empirically evaluated using the following metrics (formulas and success criteria remain as originally designed).

### 3.1 Transaction Throughput (TPS)

**Formula:**  
`TPS = Σ Tx_committed / (T_end - T_start)`

**Success Criteria:** The isolated Slave logic must sustain high TPS without degrading Master logic performance.

### 3.2 Transaction Committal Latency

**Formula:**  
`L = T_commit - T_submit`

**Success Criteria:** Sub‑second latency for Master milestones, even under heavy IoT load.

### 3.3 Hardware Resource Utilization

**Formula:**  
`Resource_Utilization = (CPU_used/CPU_allocated , RAM_used/RAM_allocated)`

**Success Criteria:** No Out‑Of‑Memory (OOM) failures during spikes.

### 3.4 Crash Fault Tolerance (CFT) Recovery Time

**Formula:**  
`T_recovery = T_new_leader_elected - T_node_failure`

**Success Criteria:** Rapid Raft leader re‑election after orderer node crash.

---

## Current Status (as of 2026‑07‑12)

| Milestone | Status | Notes |
|-----------|--------|-------|
| **M1: Infrastructure** | ✅ Complete | Network up, `mychannel` created, peers joined. Pivot documented. |
| **M2: Data Ingestion** | ✅ Complete | Python script ingested 8,000 sensor records with 100% success through Go bridge. See [M2 as‑built](docu/m2-data-ingestion/notes/m2_as_built.md). |
| **M3: Slave Chain** | ✅ Complete | `slave-twin` deployed and tested (invoke/query successful). |
| **M4: Master Chain** | ✅ Complete | Extended `slave-twin` with `RecordHandover` and `ReadHandover` – tested with `shipment-001`. See [M4 as‑built](docu/m4-master-chain/notes/m4_as_built.md). |
| **M5: Evaluation** | ⏳ Next | Stress tests, metrics collection, crash simulation. |
| **M6: Defense Prep** | ⏳ Pending | Final report, slides, video walkthrough. |

See the [milestones.md](milestones.md) file for the complete task breakdown.

---

## Documentation

All detailed notes, as‑built reports, and design documents are stored in the `docu/` folder. Key references:

- **[Global Notes](docu/_global-notes/)** – environment setup, tech stack justification.
- **[M1 Infrastructure](docu/m1-infrastructure/notes/m1_as_built.md)** – full narrative of network provisioning, version mismatches, CouchDB issues, CA vs. cryptogen, custom channel failures, and the pivot decision.
- **[M2 As‑Built](docu/m2-data-ingestion/notes/m2_as_built.md)** – ingestion engine implementation, performance results (8k records, 65 TPS).
- **[M3 Slave Chain](docu/m3-slave-chain/notes/m3_as_built.md)** – development and deployment of the `slave-twin` chaincode, including the successful invoke/query test.
- **[Architectural Pivot Rationale](docu/m3-slave-chain/notes/architecture_pivot_rationale.md)** – detailed justification of the single‑channel, logical‑segregation approach.
- **[M4 Master Chain](docu/m4-master-chain/notes/m4_as_built.md)** – development and deployment of logistics handover functions, including path correction, version increments, and successful verification.

---

## How to Run This Project (for replication)

1. Clone this repository.
2. Ensure Docker, Docker Compose, and Go are installed (see [environment setup](docu/_global-notes/environment_setup_guide.md)).
3. Navigate to `fabric-samples/test-network/` (you may need to download Fabric samples separately – see the global notes).
4. Bring up the network:
   ```bash
   ./network.sh up
   ./network.sh createChannel -c mychannel
   ```
5. Deploy the `slave-twin` chaincode:
   ```bash
   ./network.sh deployCC -ccn slave-twin -ccp ../chaincode/slave-twin/go -ccl go -c mychannel -ccv 1.0 -ccs 1
   ```
6. Start the Go API bridge (from `gateway-api/`):
   ```bash
   go run main.go
   ```
7. Run the Python ingestion script (from `docu/m2-data-ingestion/code/ingest_sensors.py`):
   ```bash
   python ingest_sensors.py
   ```
8. Query the chaincode to verify twins (commands in M3 notes).

---

## License & Acknowledgements

This project is part of a Master’s research dissertation. Datasets are sourced from Kaggle and are publicly available.