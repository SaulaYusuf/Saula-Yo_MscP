# Saula-Yo_MscP

## Project Overview

**Title:** Enhancing Supply Chain Traceability and Transparency via Master‑Slave Multi‑Chain Architectures and Digital Twins: An Empirical Performance Evaluation.

**Student:** Saula Yusuf Owolabi  
**Institution:** Ulster University  
**Date:** July 2026 (Last Updated: 2026‑07‑21)

---

## Repository Structure

```
.
├── chaincode/                 # Smart contract source code (Go)
│   └── slave-twin/            # Unified chaincode (IoT + Logistics + Metadata)
│       ├── main.go            # RecordTelemetry, RecordHandover, RecordMetadata
│       ├── go.mod
│       └── go.sum
├── data/
│   └── raw/                   # Kaggle datasets
│       ├── bdt_mba_supplychain_dataset_2024.csv   (500 records)
│       ├── shipment-sensor-dataset.csv            (8,000 records)
│       └── smart_logistics_dataset.csv            (1,000 records)
├── docu/                      # Full project documentation (milestone notes, results)
│   ├── _global-notes/         # Environment setup, tech stack justification
│   ├── m1-infrastructure/     # Network provisioning & pivot rationale
│   ├── m2-data-ingestion/     # Python ingestion engine (as‑built, charts)
│   ├── m3-slave-chain/        # Slave chaincode development & deployment
│   ├── m4-master-chain/       # Master logistics chaincode (as‑built)
│   ├── m5-evaluation/         # Performance metrics, benchmark results, graphs
│   └── m6-defense-prep/       # Final report, slides, video walkthrough (planned)
├── gateway-api/               # Go REST API bridge (HTTP → Fabric)
│   └── main.go                # /api/sensor, /api/logistics, /api/metadata
├── README.md                  # This file
├── milestones.md              # Detailed milestone map
├── benchmark_plots.png        # Throughput & latency visualisation (M5)
└── benchmark_results.json     # Raw performance data (M5)
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
- **Unified Chaincode (`slave-twin`) containing three contracts:**
  - **Slave (IoT) Contract:** `RecordTelemetry` – ingests environmental sensor telemetry, evaluates temperature thresholds, updates Digital Twin status (`NORMAL` / `SPOILED`).
  - **Master (Logistics) Contract:** `RecordHandover` – handles ownership handovers, port arrivals, and cryptographic anchoring of final milestones.
  - **Metadata Contract:** `RecordMetadata` – stores asset condition, maintenance logs, and efficiency labels.

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
| **BDT‑MBA Digital Twin State** | [`data/raw/bdt_mba_supplychain_dataset_2024.csv`](data/raw/bdt_mba_supplychain_dataset_2024.csv) | `RecordMetadata` | Asset metadata (condition score, maintenance logs, efficiency label) |

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

## Key Performance Results (M5)

The system was benchmarked with 8,000 sensor records at varying concurrency levels. The key findings are:

| Concurrency | TPS (peak) | p95 Latency | Success Rate |
|-------------|------------|-------------|--------------|
| 10          | 121        | 100 ms      | 100%         |
| 50          | **159**    | 414 ms      | 100%         |
| 100         | 154        | 780 ms      | 100%         |
| 200         | 144        | 904 ms      | 93.75%       |

**Raft Crash Recovery:** 5.264 seconds (leader re‑election after container termination).

**Graphical Summary:** See [`benchmark_plots.png`](benchmark_plots.png) for throughput and latency visualisations.
![`benchmark_plots.png`](benchmark_plots.png)

These results validate that the architecture sustains **~150 TPS** with sub‑second latency for critical logistics milestones, and recovers automatically from node failures.

---

## Current Status (as of 2026‑07‑21)

| Milestone | Status | Notes |
|-----------|--------|-------|
| **M1: Infrastructure** | ✅ Complete | Network up, `mychannel` created, peers joined. Pivot documented. |
| **M2: Data Ingestion** | ✅ Complete | All three datasets ingested (8,000 + 1,000 + 500 records) with 100% success. See [M2 as‑built](docu/m2-data-ingestion/notes/m2_as_built.md). |
| **M3: Slave Chain** | ✅ Complete | `RecordTelemetry` and `ReadTwin` deployed and tested. |
| **M4: Master Chain** | ✅ Complete | `RecordHandover`, `ReadHandover`, `RecordMetadata`, `ReadMetadata` deployed and tested. See [M4 as‑built](docu/m4-master-chain/notes/m4_as_built.md). |
| **M5: Evaluation** | ✅ Complete | Benchmarking, latency analysis, resource monitoring, and crash fault tolerance tested. See [M5 as‑built](docu/m5-evaluation/notes/m5_as_built.md). |
| **M6: Defense Prep** | ⏳ Pending | Final report, slides, video walkthrough. |

See the [milestones.md](milestones.md) file for the complete task breakdown.

---

## Documentation

All detailed notes, as‑built reports, and design documents are stored in the `docu/` folder. Key references:

- **[Global Notes](docu/_global-notes/)** – environment setup, tech stack justification.
- **[M1 Infrastructure](docu/m1-infrastructure/notes/m1_as_built.md)** – full narrative of network provisioning, version mismatches, CouchDB issues, CA vs. cryptogen, custom channel failures, and the pivot decision.
- **[M2 As‑Built](docu/m2-data-ingestion/notes/m2_as_built.md)** – ingestion engine implementation, performance results (all three datasets), and comparison chart.
- **[M3 Slave Chain](docu/m3-slave-chain/notes/m3_as_built.md)** – development and deployment of the `slave-twin` chaincode, including the successful invoke/query test.
- **[Architectural Pivot Rationale](docu/m3-slave-chain/notes/architecture_pivot_rationale.md)** – detailed justification of the single‑channel, logical‑segregation approach.
- **[M4 Master Chain](docu/m4-master-chain/notes/m4_as_built.md)** – development and deployment of logistics handover functions, including path correction, version increments, and successful verification.
- **[M5 As‑Built](docu/m5-evaluation/notes/m5_as_built.md)** – comprehensive performance evaluation: throughput, latency, resource utilisation, and crash recovery.

---

## How to Run This Project (for replication)

### Prerequisites

- Docker & Docker Compose
- Go (1.20+)
- Python 3.12+ with `pip`
- Hyperledger Fabric binaries (download `fabric-samples`)

### Step 1: Network Setup

```bash
cd fabric-samples/test-network
./network.sh up
./network.sh createChannel -c mychannel
```

### Step 2: Deploy Chaincode

```bash
./network.sh deployCC -ccn slave-twin -ccp ../chaincode/slave-twin/go -ccl go -c mychannel -ccv 2.7 -ccs 4
```

### Step 3: Start the Go API Bridge

```bash
cd gateway-api
go run main.go
```

### Step 4: Run Ingestion (All Datasets)

```bash
# Sensor (8,000 records)
python docu/m2-data-ingestion/code/ingest_sensors.py

# Logistics (1,000 records)
python docu/m2-data-ingestion/code/ingest_logistics.py

# Metadata (500 records)
python docu/m2-data-ingestion/code/ingest_metadata.py
```

### Step 5: Run Benchmarks (M5)

```bash
cd docu/m5-evaluation/code
python benchmark_runner.py
python plot_benchmark.py
bash crash_test.sh
```

### Step 6: Query Chaincode

```bash
# Set environment variables (see M3 notes)
peer chaincode query -C mychannel -n slave-twin -c '{"function":"ReadTwin","Args":["sensor-001"]}'
peer chaincode query -C mychannel -n slave-twin -c '{"function":"ReadHandover","Args":["shipment-001"]}'
peer chaincode query -C mychannel -n slave-twin -c '{"function":"ReadMetadata","Args":["A0001"]}'
```

---

## License & Acknowledgements

This project is part of a Master’s research dissertation. Datasets are sourced from Kaggle and are publicly available.

The Hyperledger Fabric network is built on the official `fabric-samples` test network, adapted for multi‑contract deployment.

---

**Last Updated:** 2026‑07‑21