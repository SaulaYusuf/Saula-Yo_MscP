
# M3: Slave Chain Logic – Digital Twin Threshold Engine

**Date:** 2026‑07‑06  
**Student:** Saula Yusuf Owolabi  
**Milestone:** M3 – Smart Contract Development & Deployment

---

## 1. Goal

Develop a smart contract that:
- Defines a **Digital Twin** schema for an asset (sensor ID, temperature, humidity, status, timestamp).
- Implements a **threshold engine** – if temperature exceeds 8.0°C, automatically set status to `"SPOILED"`.
- Records incoming telemetry as **state objects** on the ledger.
- Allows querying the current twin state.

This contract is the **Slave** side of the Master‑Slave architecture – it ingests high‑frequency IoT data and performs automated arbitration.

---

## 2. Implementation

I chose **Go** as the chaincode language because the Fabric documentation and community heavily support it. I scaffolded the project:

```bash
mkdir -p ../chaincode/slave-twin/go
cd ../chaincode/slave-twin/go
go mod init slave-twin
```

I wrote `main.go` (see the attached code) with:

- `AssetTwin` struct – fields for sensor ID, temperature, humidity, status, timestamp.
- `RecordTelemetry` function – takes sensor ID, temperature, humidity, timestamp; evaluates the threshold; stores the twin.
- `ReadTwin` function – retrieves the twin by sensor ID.

The threshold logic is simple:

```go
if tempC > 8.0 {
    status = "SPOILED"
} else {
    status = "NORMAL"
}
```

I used the `fabric-contract-api-go` SDK v1.2.2 (resolved via `go mod tidy`).

---

## 3. Packaging & Installation

I had the chaincode source in `../chaincode/slave-twin/go`. I used the `deployCC` script from `test-network` (but it was failing for custom channels, as explained in M1). After the pivot to `mychannel`, I deployed on the default channel:

```bash
./network.sh deployCC -ccn slave-twin -ccp ../chaincode/slave-twin/go -ccl go -c mychannel -ccv 1.0 -ccs 1
```

The script:
- Packaged the Go code into `slave-twin.tar.gz`.
- Installed it on both peers (Org1 and Org2).
- Approved the chaincode definition for both organisations.
- Committed the definition to `mychannel`.

I saw the successful output:

```
Chaincode definition committed on channel 'mychannel'
Version: 1.0, Sequence: 1, Approvals: [Org1MSP: true, Org2MSP: true]
```

---

## 4. Verification – Invoke & Query

To ensure the chaincode was working, I ran an invoke with sample telemetry:

```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n slave-twin --peerAddresses localhost:7051 --tlsRootCertFiles ... --peerAddresses localhost:9051 --tlsRootCertFiles ... -c '{"function":"RecordTelemetry","Args":["sensor-001","7.5","45.0","2026-07-06T12:00:00Z"]}'
```

I received `status:200` – success.

Then I queried the twin:

```bash
peer chaincode query -C mychannel -n slave-twin -c '{"function":"ReadTwin","Args":["sensor-001"]}'
```

The output:

```json
{"sensor_id":"sensor-001","temp_c":7.5,"humidity":45,"status":"NORMAL","timestamp":"2026-07-06T12:00:00Z"}
```

**Perfect.** The threshold engine correctly kept the status as `"NORMAL"` because 7.5 < 8.0.

---

## 5. Testing the Threshold

I also invoked with a temperature above the threshold (e.g., 9.0°C) and verified the status changed to `"SPOILED"` – the automated arbitration works.

---

## 6. Challenges & Decisions

- **Go module dependencies** – initially I had trouble with the Fabric SDK path. Running `go mod tidy` inside the chaincode directory resolved it.
- **Package ID** – each deployment generates a unique package ID. I had to use the correct ID when approving manually (though the script handled that automatically).
- **The pivot** – I deployed on `mychannel` instead of the original `slavechannel`. This was a necessary workaround as explained in M1, but the chaincode logic is identical and functional.

---

## 7. Current State (End of M3)

- **Chaincode name:** `slave-twin`
- **Channel:** `mychannel`
- **Status:** Installed, approved, committed, and verified with invoke/query.
- **Digital Twin** – fully functional threshold engine.

**Next Steps:** Develop the `master-logistics` chaincode to handle ownership handovers, and then build the Python ingestion engine (M2 and M4).
