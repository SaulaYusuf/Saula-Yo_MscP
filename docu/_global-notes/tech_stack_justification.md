# 1. Documenting the Tech Stack

I have updated the `environment_setup_guide.md` file (which lives in th `docu/_global-notes/` folder) to include this justification section. 


| Tool | Role | Justification |
| --- | --- | --- |
| **Hyperledger Fabric** | Distributed Ledger | Provides permissioned, modular consensus (Raft) and private channels, which is essential for segregating IoT noise from Master audit trails. |
| **Go (Golang)** | Smart Contracts | Fabric's native language; provides high performance, concurrency, and low memory overhead, critical for IoT threshold engines. |
| **Python (3.10+)** | Ingestion Engine | Unmatched library support (`pandas`, `asyncio`) for parsing and streaming high-velocity Kaggle datasets into the ledger. |
| **Docker** | Containerization | Ensures 1:1 environmental consistency between development and production, essential for reproducible research. |
| **Caliper** | Benchmarking | The official Hyperledger tool for standardized TPS, latency, and resource stress-testing. |

---