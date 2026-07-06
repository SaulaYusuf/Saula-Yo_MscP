# Saula-Yo_MscP

---

# System Architecture and Evaluation Methodology

## 1. Introduction & Architectural Overview

The proposed research addresses the latency and throughput bottlenecks inherent in monolithic blockchain supply chain networks, particularly when handling high-frequency Internet of Things (IoT) telemetry data. To resolve this, the system employs a Master-Slave Multi-Chain topology deployed via containerized Hyperledger Fabric nodes.

The architecture segregates data flows based on volume and velocity. High-frequency sensor data is isolated to a **Slave Channel**, where it is evaluated against a Digital Twin threshold engine. Macro-logistics milestones and cryptographic state-change proofs are anchored to a **Master Channel** governed by a Raft Crash Fault-Tolerant consensus mechanism, ensuring audit trail immutability without network congestion.

## 2. Dataset Mapping & Ingestion Strategy

To empirically test the system, the simulation relies on three specific data structures and a synthetic ingestion engine, mapped directly to the isolated channels:

* **[Dataset 1](data/raw/smart_logistics_dataset.csv): The Routing Payload (Smart Logistics Supply Chain Dataset)**
* *Target:* Master Channel
* *Function:* Tracks the macro movement of physical assets across geographic nodes (e.g., origin ports, destination warehouses). This dataset tests the ledger's ability to maintain the global custody chain.


* **[Dataset 2](data/raw/shipment-sensor-dataset.csv): The IoT Penalty Payload (Cold-Chain Shipment Silent Failure Dataset)**
* *Target:* Slave Channel (Smart Contract Trigger)
* *Function:* Provides raw environmental sensor data (temperature, humidity). This feeds the threshold engine, proving the system can automatically execute business logic (e.g., flagging goods as 'SPOILED') in real-time.


* **[Dataset 3](data/raw/bdt_mba_supplychain_dataset_2024.csv): The Digital Twin State Payload (BDT-MBA Supply Chain Dataset)**
* *Target:* Slave Channel (State Object mapping)
* *Function:* Provides the exact schema required to map physical asset conditions to digital replicas. It validates the synchronization between the physical world and the on-chain digital twin.


* **The Ingestion Engine (Synthetic Stress Payload)**
* *Target:* System-wide Load Generation
* *Function:* An asynchronous Python pipeline that aggressively parses and streams the aforementioned datasets into the network. This synthetic load acts as the primary tool for stressing the infrastructure to measure performance limits.



## 3. System Evaluation Metrics Framework

The success of this Master-Slave architecture will be empirically evaluated based on infrastructural scalability, deterministic network performance, and system resilience. The evaluation employs the Hyperledger Performance Metrics Framework, utilizing Hyperledger Caliper and container telemetry to extract the following mathematical benchmarks:

### 3.1 Transaction Throughput (TPS)

Throughput defines the maximum rate at which the network successfully processes and commits transactions to the ledger.

**Formula:**


$$TPS = \frac{\sum_{i=1}^{n} Tx_{committed}}{T_{end} - T_{start}}$$

*Success Criteria:* The benchmarking phase must mathematically demonstrate that the isolated Slave Channel can sustain significantly higher TPS loads (accommodating high-frequency IoT data) without degrading the performance of the Master Channel.

### 3.2 Transaction Committal Latency

Latency measures the exact time delay from when the Python ingestion engine submits a transaction proposal to when the block is permanently settled and verifiable on the ledger.

**Formula:**


$$L = T_{commit} - T_{submit}$$

*Success Criteria:* The evaluation must prove that by isolating telemetry noise to the Slave Channel, the Master Channel maintains sub-second latency for critical logistics handovers, ensuring real-time audit visibility.

### 3.3 Hardware Resource Utilization

This metric tracks the physical computational strain exerted across the containerized Docker nodes.

**Formula:**


$$Resource\_Utilization = \left( \frac{CPU_{used}}{CPU_{allocated}}, \frac{RAM_{used}}{RAM_{allocated}} \right)$$

*Success Criteria:* Real-time monitoring via Docker Stats must validate that the multi-chain partition effectively balances the processing load, preventing individual consensus nodes from experiencing Out-Of-Memory (OOM) failures during high-volume data ingestion spikes.

### 3.4 Crash Fault Tolerance (CFT) Recovery Time

This evaluates the network's resilience when a core consensus node fails mid-transaction.

**Formula:**


$$T_{recovery} = T_{new\_leader\_elected} - T_{node\_failure}$$

*Success Criteria:* During a simulated Orderer node container termination, the Raft consensus mechanism must demonstrate rapid leader re-election, proving the system can self-heal and maintain data integrity without dropping the audit trail.

---