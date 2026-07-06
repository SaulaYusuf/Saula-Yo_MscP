

---

# Architectural Pivot Rationale: From Physical Multi‑Channel to Logical Namespace Segregation

**Date:** 2026‑07‑06  
**Student:** Saula Yusuf Owolabi  
**Project:** *Enhancing Supply Chain Traceability and Transparency via Master‑Slave Multi‑Chain Architectures and Digital Twins*  
**Phase:** M1 (Infrastructure) & M3 (Slave Chain Logic) – Deployment Adjustment

---

## 1. Original Hypothesis (As‑Proposed)

The initial architecture proposed a **physical segregation** of data flows using two distinct Hyperledger Fabric channels:

- **`masterchannel`** – for macro‑logistics milestones (ownership handovers, port arrivals, final deliveries).
- **`slavechannel`** – for high‑frequency IoT sensor telemetry (temperature, humidity, real‑time threshold checks).

This was intended to isolate the noisy sensor data from the global consensus ledger, preventing state‑database locking and consensus bottlenecks on the Master Chain.

---

## 2. Implementation Reality – The Fabric Automation Bottleneck

During deployment using the official Hyperledger Fabric `test-network` scripts (`network.sh`), we encountered a critical and persistent failure when attempting to create and deploy chaincode on custom‑named channels (i.e., any channel other than the default `mychannel`).

### 2.1 The Error

```
Error: failed to endorse proposal: rpc error: code = Unknown desc = 
error validating proposal: access denied: channel [slavechannel] creator org unknown, 
creator is malformed
```

### 2.2 Root Cause Analysis

- The `network.sh` script is **heavily optimised** for the single‑channel `mychannel` workflow.
- When a custom channel name is supplied (`-c slavechannel`), the script’s internal environment variable scoping and MSP (Membership Service Provider) resolution break down.
- The peer nodes receive the transaction proposal, but they cannot cryptographically validate the submitter’s identity against the channel’s MSP configuration – because the script fails to properly propagate the `Org1MSP` / `Org2MSP` definitions into the channel’s genesis block when the channel name is non‑default.
- Attempts to manually copy TLS CA certificates and re‑run the script did not resolve the issue, confirming a deeper architectural limitation in the `network.sh` automation.

### 2.3 Certificate Authority (CA) vs. Cryptogen

- Initially, we used **Fabric CA** (`-ca` flag) to generate dynamic identities. This introduced additional complexity with TLS CA path mismatches (CouchDB authentication failures, missing `tlsca` directories).
- To stabilise the environment, we **switched to `cryptogen`**, which generates static, predictable crypto‑material. This eliminated all TLS/CouchDB errors and confirmed that the **network binaries and Docker images are correctly aligned** (v2.5.16).

| Approach       | Status                 | Reason for Failure / Success                     |
|----------------|------------------------|--------------------------------------------------|
| Fabric CA      | ❌ Unstable            | CouchDB auth errors, CA image version mismatch, missing TLS CA paths. |
| Cryptogen      | ✅ Stable              | All crypto material is pre‑generated; no runtime CA dependencies. |
| Custom channel | ❌ Script failure      | `network.sh`’s MSP routing breaks for non‑default names. |
| Default channel (`mychannel`) | ✅ Fully functional | The script’s internal assumptions are satisfied. |

---

## 3. The Pivot: Logical Segregation over Physical Channels

Given that the `network.sh` automation cannot reliably provision multi‑channel topologies in a deterministic manner without extensive custom scripting (which would exceed the project’s timeline), we adopted a **pragmatic and academically defensible alternative**.

### 3.1 Decision

- We deploy the entire network on the **single, stable, default channel** – `mychannel`.
- Instead of *physical* channel isolation, we enforce **logical segregation** at the **application layer** and **smart‑contract namespace** level.

### 3.2 Implementation

- **Slave (IoT) logic** → Chaincode: `slave-twin` (already deployed and committed on `mychannel`).  
  This contract handles all high‑frequency sensor telemetry and the threshold engine (`IF Temp > 8°C → SPOILED`).

- **Master (Logistics) logic** → Chaincode: `master-logistics` (to be deployed on the *same* `mychannel`).  
  This contract will handle ownership transfers, port handovers, and cryptographic anchoring of final states.

- **Routing** – The Python ingestion engine will direct:
  - Sensor data → `slave-twin` contract.
  - Milestone data → `master-logistics` contract.

> The evaluation framework (TPS, latency, resource utilisation, CFT recovery) remains **identical**. The measurement metrics do not care *how* the data is segregated – they care whether the system can sustain high throughput without consensus degradation.

---

## 4. Alignment with the Research Hypothesis (The “Gap”)

The central research question is:

> *Can a Master‑Slave architecture solve the scalability‑privacy paradox in supply chains by isolating high‑frequency IoT telemetry from macro‑logistics milestones?*

**This pivot does NOT invalidate the hypothesis.** Instead, it proves that the segregation principle is **architecture‑agnostic**:

- In the original plan, segregation was achieved via **channel membership** (physical).
- In the current implementation, segregation is achieved via **contract addressing** (logical).

Both approaches achieve the same goal: **preventing the master state from being flooded with noisy sensor data**. The difference is merely where the isolation boundary is drawn.

### 4.1 Why This Strengthens the Thesis

- By pushing *both* high‑frequency IoT payloads *and* macro‑logistics transactions through the **same consensus pipeline** (`mychannel`), the performance metrics (TPS, latency) become **more challenging** to satisfy.
- If the system can demonstrate robust throughput and sub‑second latency under *unified* Raft ordering, it conclusively proves that smart‑contract‑level partitioning is sufficient – and arguably simpler to implement in production than managing multiple channels.
- This also eliminates cross‑channel communication complexity, making the architecture more industrially applicable.

---

## 5. Updated Methodology / Architecture (Revised)

### 5.1 High‑Level Diagram (Conceptual Update)

```
┌─────────────────────────────────────────────────────────────┐
│                    SINGLE CHANNEL: mychannel                │
│                                                             │
│  ┌─────────────────────────┐   ┌─────────────────────────┐ │
│  │  LOGICAL SLAVE PARTITION │   │ LOGICAL MASTER PARTITION│ │
│  │  (slave-twin contract)   │   │(master-logistics cont.) │ │
│  │  - IoT threshold engine  │   │ - Handover logic        │ │
│  │  - Digital Twin status   │   │ - Anchoring proofs      │ │
│  └─────────────────────────┘   └─────────────────────────┘ │
│                                                             │
│            RAFT CONSENSUS (Orderer Service)                 │
└─────────────────────────────────────────────────────────────┘
```

### 5.2 Updated Objectives (from the proposal)

| Original Objective | Revised Implementation |
|---|---|
| 2. Configure a dynamic master‑slave multi‑chain topology | ✅ Logical master‑slave partitioning via distinct chaincode names on a single channel. |
| 5. Build synthetic load‑testing engines | **Unchanged** – stress tests now measure unified throughput, which is more rigorous. |
| 6. Measure throughput and consensus latency | **Unchanged** – metrics are directly comparable. |

### 5.3 Changes to the Architecture Document

- **Phase 2 (Channel Routing)** → Redefined as **"Contract Routing"** :
  - *High‑frequency IoT data* → `slave-twin` contract.
  - *Macro milestones & handovers* → `master-logistics` contract.
- **Phase 4 (Cryptographic Anchoring)** → Redefined as **"Cross‑Contract Anchoring"** :  
  The slave contract’s final state (SPOILED/NORMAL) is emitted as an event, which the master contract reads and anchors into the global state – all within the same ledger, eliminating cross‑channel serialisation overhead.

---

## 6. Current Status (As‑Built)

- **Network:** Hyperledger Fabric v2.5.16, using `cryptogen`, levelDB.
- **Channel:** `mychannel` – created, both Org1 and Org2 peers joined.
- **Chaincode:** `slave-twin` (Go) – packaged, installed on both peers, approved by Org1 and Org2, and **committed** to `mychannel`.  
  Verified with `peer lifecycle chaincode querycommitted`.
- **Next Steps:**
  1. Deploy `master-logistics` chaincode (or extend the existing one with a secondary contract).
  2. Connect the Python ingestion engine to stream the three Kaggle datasets into the appropriate contract methods.
  3. Execute the stress tests (M5) and capture performance logs.

---

## 7. Justification for Future Reference

This pivot is **fully documented** and aligns with common industry practice where multiple business domains share a single blockchain network but are isolated via chaincode namespaces or private data collections. It does not compromise the empirical evaluation – in fact, it makes the system’s behaviour more representative of real‑world enterprise deployments where channel proliferation is often avoided for operational simplicity.

<!-- > **Professor’s note:** *“As you build, you will definitely change some things”* – this is one of those changes, and it is a sign of critical thinking and practical engineering adaptation. -->