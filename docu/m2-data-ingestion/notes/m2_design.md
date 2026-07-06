# M2: Data Ingestion Engine – Design & Implementation Plan

**Date:** 2026‑07‑06  
**Student:** Saula Yusuf Owolabi  
**Milestone:** M2 – Python Asynchronous Client Pipeline

---

## 1. Objective

Build an asynchronous Python client that:
- Reads the three Kaggle CSV datasets (Smart Logistics, Cold‑Chain Sensor, BDT‑MBA Digital Twin State).
- Transforms each row into a JSON transaction proposal.
- Streams the proposals into the Hyperledger Fabric network via the Fabric Gateway SDK.
- Routes sensor telemetry to the `slave-twin` chaincode.
- Routes logistics milestones to the `master-logistics` chaincode (once developed).

---

## 2. Dataset Mapping

| Dataset | Source | Target Chaincode | Key Fields |
|---------|--------|------------------|------------|
| Smart Logistics Supply Chain | `smart_logistics_dataset.csv` | `master-logistics` (to be created) | Shipment_ID, Origin_Port, Destination_Port, Current_Status, Timestamp |
| Cold‑Chain Silent Failure | `shipment-sensor-dataset.csv` | `slave-twin` (already deployed) | Sensor_ID, Temp_Max_C, Humidity_Pct, Timestamp |
| BDT‑MBA Supply Chain | `bdt_mba_supplychain_dataset_2024.csv` | `slave-twin` (twin metadata) | Asset_ID, Condition_Score, Maintenance_Logs, Operational_Lifecycle |

The third dataset will be used to initialise or update twin metadata.

---

## 3. Technical Design

- **Language:** Python 3.8+
- **SDK:** `fabric-sdk-py` (v2.2 or later)
- **Asynchronous:** `asyncio` for high‑throughput streaming.
- **CSV parsing:** `pandas` or `csv` module – I’ll use `csv` for memory efficiency.
- **Transaction proposal:** Each row → JSON object → `invoke` call to the respective chaincode function.

**Routing Logic:**
- If row contains `Sensor_ID` → route to `slave-twin.RecordTelemetry`.
- If row contains `Shipment_ID` → route to `master-logistics.RecordHandover` (to be implemented).

**Concurrency:** I will create a fixed‑size worker pool (e.g., 10 async tasks) to submit transactions concurrently, simulating high‑frequency IoT.

---

## 4. Challenges Anticipated

- **Rate limiting:** Fabric’s orderer may throttle if TPS exceeds its capacity. I will implement back‑pressure with a queue.
- **Error handling:** Network disconnections, endorsement failures – must retry with exponential backoff.
- **Performance logging:** I will capture submit‑time and commit‑time for each transaction to compute latency.

---

## 5. Implementation Steps

1. Set up Python virtual environment with `fabric-sdk-py`, `asyncio`, `pandas`.
2. Write a connection profile (YAML) to connect to the Fabric network.
3. Implement a `CSVStreamer` class that yields rows asynchronously.
4. Write two client methods: `submit_sensor_data(row)` and `submit_logistics(row)`.
5. Test with a small subset of the datasets.
6. Run a full‑scale ingestion and monitor Docker stats.

---

## 6. Current Status

- **In progress** – design phase.
- The `slave-twin` chaincode is ready.
- The `master-logistics` chaincode is yet to be developed (M4).

**Next:** Build the client and test with sensor data first.
```

---

### 4. Update the Pivot Rationale

Make sure `architecture_pivot_rationale.md` includes the **successful chaincode test** as proof that the new approach works. I’ll add a section:

```markdown
## 8. Verification of the Pivot

On 2026‑07‑06, I deployed the `slave-twin` chaincode on `mychannel` and successfully invoked `RecordTelemetry` and queried the twin. The threshold engine performed as expected. This confirms that the logical segregation approach is fully functional and ready for performance evaluation.
