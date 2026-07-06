# Architectural Pivot Rationale: Infrastructure vs. Logical Partitioning

**Date:** 2026-07-06
**Phase:** M1 & M3 Infrastructure Deployment

## 1. The Original Hypothesis
The initial architecture proposed a physical segregation of data using two distinct Hyperledger Fabric channels (`masterchannel` and `slavechannel`). The goal was to isolate high-frequency IoT telemetry (Slave) from macro-logistics milestones (Master) to prevent state-database locking and consensus bottlenecks.

## 2. The Technical Blocker: Automation vs. State Corruption
During deployment, the official Hyperledger `test-network` scripts (`network.sh`) critically failed to manage the Membership Service Provider (MSP) contexts for custom, multi-channel configurations. 
- **The Error:** `access denied: creator org unknown, creator is malformed.`
- **The Cause:** The automation scripts are hardcoded with identity assumptions optimized for a single default channel (`mychannel`). Attempting to force multi-channel provisioning resulted in ghost certificates, orphaned Docker volumes, and a fatal desynchronization between the physical nodes and their cryptographic identities.
- **Identity Shift:** To stabilize the local environment, the network was transitioned from dynamic Fabric-CAs back to static `cryptogen` configurations, ensuring rigid, predictable certificates.

## 3. The Architectural Pivot
Rather than spending weeks writing custom bash scripts to override Fabric's core testing assumptions, the architecture was pivoted to mimic an enterprise "unified ledger" approach.
- **The Fix:** The network was successfully provisioned using the native, stable `mychannel` configuration. The `slave-twin` chaincode was successfully installed, approved, and committed by Org1 and Org2 on this channel.

## 4. Preserving the Thesis Methodology (Master-Slave Logic)
The core research gap—preventing IoT noise from congesting the logistics audit trail—is preserved. Instead of isolating the data via *physical* channels, the system now enforces isolation via *logical* namespaces and Application-Level Routing.

- **Logical Slave Partition:** The deployed `slave-twin` smart contract acts as the dedicated ingestion endpoint for the high-frequency environmental data (Temp/Humidity).
- **Logical Master Partition:** A secondary contract (or segregated namespace) will be deployed to handle the ownership handovers.
- **The Evaluation Impact:** This pivot makes the empirical evaluation even more rigorous. By pushing both high-frequency IoT payloads and macro-logistics data through the *same* consensus pipeline (`mychannel`), the Caliper benchmarks will directly measure the Raft Orderer's ability to handle unified, high-stress throughput, proving whether smart-contract level segregation is sufficient to prevent network failure.