## Complete Milestone Map: Start to Finish

**M1: INFRASTRUCTURE & FABRIC SANDBOX**
**Goal:** Provision the base Docker containers, initialize the Hyperledger Fabric test network, and configure the isolated Master and Slave channels.

| ID | Sub-Milestone | Micro-Milestones |
| --- | --- | --- |
| M1.1 | Documentation scaffold | Create docu folder structure; initialize cross-OS READMEs (Linux & Windows WSL2). |
| M1.2 | Base Network Provisioning | Configure `docker-compose.yaml` for Orderers, Certificate Authorities, and Peer nodes. |
| M1.3 | Channel Creation | Spin up the network; create `masterchannel` and `slavechannel`; join respective peers. |
| M1.4 | Gateway Setup | Configure the Fabric Gateway SDK to allow external client connections to the network. |
| M1.5 | Network API Access | Map container ports to localhost to allow API interactions from the host OS. |
| M1.6 | Infrastructure verification | Run network health checks; log container resource usage (RAM/CPU) during idle state. |

**M2: DATA INGESTION ENGINE**
**Goal:** Build the asynchronous Python client to parse Kaggle CSV datasets and stream JSON payloads to the network.

| ID | Sub-Milestone | Micro-Milestones |
| --- | --- | --- |
| M2.1 | Environment setup | Initialize Python `venv`; install `pandas`, `asyncio`, `fabric-sdk-py`. |
| M2.2 | Dataset preparation | Download Logistics route dataset and Cold-Chain IoT dataset; clean missing values. |
| M2.3 | Python Parser script | Write script to read CSV rows sequentially and convert them to standard JSON dictionaries. |
| M2.4 | Asynchronous streaming | Implement `asyncio` to simulate live, high-frequency edge sensor data pushes. |
| M2.5 | Payload Routing Logic | Build the logic to send IoT data to the Slave Channel and Logistics data to the Master Channel. |
| M2.6 | Ingestion verification | Print JSON payloads to terminal; measure rows-per-second ingestion rate. |

**M3: SLAVE CHAIN & DIGITAL TWIN LOGIC**
**Goal:** Develop and deploy the smart contract on the Slave channel to monitor IoT thresholds and mutate Digital Twin states.

| ID | Sub-Milestone | Micro-Milestones |
| --- | --- | --- |
| M3.1 | Smart Contract Scaffold | Initialize Go module (`go mod init`); import Fabric contract API. |
| M3.2 | Digital Twin Schema | Define Go structs for the Asset Twin, including `Sensor_ID`, `Temp`, and `Status`. |
| M3.3 | Threshold Engine | Implement logic: `IF Temp > Max_C THEN Status = SPOILED`. |
| M3.4 | Slave Chain Deployment | Package, install, and approve the Go chaincode specifically on the Slave peers. |
| M3.5 | Local Integration Test | Fire manual JSON payloads via CLI; verify Digital Twin status changes based on temperature. |

**M4: MASTER CHAIN & CROSS-CHAIN ANCHORING**
**Goal:** Develop the Master chain smart contract for macro milestones and implement the cryptographic bridge from the Slave chain.

| ID | Sub-Milestone | Micro-Milestones |
| --- | --- | --- |
| M4.1 | Master Contract Scaffold | Initialize separate Go module for Master logistics logic. |
| M4.2 | Handover Logic | Implement functions for ownership transfers and shipping port arrivals. |
| M4.3 | Cross-Chain Anchor | Build the event-listener: when Slave chain marks `SPOILED`, push cryptographic hash to Master chain. |
| M4.4 | Master Chain Deployment | Package, install, and approve the Go chaincode on Master peers. |
| M4.5 | Full Pipeline Test | Run Python ingestion -> Slave Chain Threshold -> Master Chain Anchor. Verify full lifecycle. |

**M5: SYSTEM EVALUATION & METRICS**
**Goal:** Run the full simulation, stress test the network, measure consensus latency under load, and generate performance graphs.

| ID | Sub-Milestone | Micro-Milestones |
| --- | --- | --- |
| M5.1 | Benchmark Suite Design | Define metrics: Transaction Throughput (TPS), Read Latency, CPU/Memory load per container. |
| M5.2 | High-Load Stress Test | Blast the Python ingestion engine at maximum speed; monitor Docker container limits. |
| M5.3 | Raft Crash Simulation | Manually kill an Orderer node mid-transaction; document Raft fault tolerance and recovery time. |
| M5.4 | Data Parsing & Graphing | Write Python `matplotlib` scripts to plot TPS, Latency, and RAM usage over time. |
| M5.5 | The Gap Proof | Compare benchmark data against monolithic network assumptions; prove the Master-Slave efficiency. |

**M6: DOCUMENTATION & DEFENSE PREP**
**Goal:** Finalize the academic report, export results, and build the defense slides.

| ID | Sub-Milestone | Micro-Milestones |
| --- | --- | --- |
| M6.1 | Result Formatting | Compile all graphs, node crash logs, and TPS metrics into formal tables. |
| M6.2 | Video Walkthrough | Screen-record the full execution, side-by-side with Docker stats and Python output. |
| M6.3 | Slide Deck Creation | Map the flowchart, Excel ledger, and final graphs to presentation slides. |
| M6.4 | Final Code Cleanup | Add explicit comments, clean up `docker-compose.yaml`, ensure GitHub repo is pristine. |
