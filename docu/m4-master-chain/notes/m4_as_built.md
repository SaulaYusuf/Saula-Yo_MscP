# M4: Master Chain Logic – Logistics Handovers

**Date:** 2026‑07‑12  
**Student:** Saula Yusuf Owolabi  
**Milestone:** M4 – Master Logistics Chaincode Development & Deployment

---

## 1. Objective

Extend the existing `slave-twin` chaincode to handle **macro‑logistics milestones** – ownership handovers, port arrivals, and final deliveries. This completes the **Master** side of the Master‑Slave architecture, enabling full separation of IoT telemetry (Slave) from logistics tracking (Master) within the same channel.

---

## 2. Design Decision: Extend vs. New Chaincode

I considered two approaches:

1. **Create a separate chaincode package** (`master-logistics`) – would require a full separate deployment lifecycle (package, install, approve, commit).
2. **Extend the existing `slave-twin` chaincode** – add new functions to the same contract, keep deployment unified.

I chose **Option 2** because:
- Reduces deployment complexity – only one chaincode to manage.
- Both contracts share the same channel and endorsement policies.
- Fabric supports multiple functions in a single contract – this is the recommended pattern for related business logic.

---

## 3. Implementation

### 3.1 Adding the Logistics Struct

I added a new struct to `main.go`:

```go
type LogisticsRecord struct {
    ShipmentID  string `json:"shipment_id"`
    Origin      string `json:"origin"`
    Destination string `json:"destination"`
    Status      string `json:"status"` // e.g., "IN_TRANSIT", "ARRIVED", "DELIVERED"
    Timestamp   string `json:"timestamp"`
}
```

### 3.2 Adding the Functions

I implemented two new functions:

- **`RecordHandover`** – accepts shipment ID, origin, destination, status, timestamp; stores the record using `shipmentID` as the key.
- **`ReadHandover`** – retrieves a logistics record by shipment ID.

The `SmartContract` struct now contains both IoT and logistics methods – it serves as a unified contract for the entire supply chain.

---

## 4. Deployment Challenges & Resolutions

### 4.1 The "Function Not Found" Error (Version 2.0 / 2.1)

When I first deployed version 2.0 with the new functions, the invoke failed with:

```
Error: endorsement failure during invoke. 
response: status:500 message:"Function RecordHandover not found in contract SmartContract"
```

I verified that `main.go` contained the new functions, but the chaincode container was still running the old code. The peer had cached the previous version and refused to rebuild.

**Attempt 1 – Remove Containers:** I removed the chaincode containers and images using `docker rm -f` and `docker rmi -f`. When I tried to reinstall the same package, the peer responded with:

```
Error: chaincode install failed with status: 500 - chaincode already successfully installed
```

**Attempt 2 – Version 2.1:** I incremented the version to 2.1 and sequence to 3. The package installed and approved successfully, but the invoke still failed with the same "Function not found" error.

### 4.2 The Root Cause – Wrong Source Path

I extracted the chaincode package to verify its contents:

```bash
mkdir -p /tmp/cc-verify
tar -xzf slave-twin-2.1.tar.gz
tar -xzf code.tar.gz
grep RecordHandover src/main.go
```

**The `RecordHandover` function was NOT in the packaged source.** The packaging command was using the wrong path.

I was running `peer lifecycle chaincode package` from `fabric-samples/test-network/` with `--path ../chaincode/slave-twin/go`. This pointed to `fabric-samples/chaincode/slave-twin/go`, which contained the **old** version of the code. My updated source was in `~/Projects/Saula_MscP/chaincode/slave-twin/go`.

### 4.3 The Fix

I copied the updated `main.go` to the correct location:

```bash
cp ../../chaincode/slave-twin/go/main.go ../chaincode/slave-twin/go/main.go
```

Then I created a **new package (version 2.2)** with the correct source:

```bash
peer lifecycle chaincode package slave-twin-2.2.tar.gz \
    --path ../chaincode/slave-twin/go \
    --lang golang \
    --label slave-twin_2.2
```

I verified the package:

```bash
mkdir -p /tmp/cc-verify-final
tar -xzf slave-twin-2.2.tar.gz
tar -xzf code.tar.gz
grep RecordHandover src/main.go
# Output: // RecordHandover stores a logistics milestone
```

The function was now in the package.

### 4.4 Successful Deployment

I installed version 2.2 on both peers, approved it with sequence 4, and committed it:

```bash
peer lifecycle chaincode install slave-twin-2.2.tar.gz
peer lifecycle chaincode approveformyorg ... --version 2.2 --sequence 4
peer lifecycle chaincode commit ... --version 2.2 --sequence 4
```

---

## 5. Verification

### 5.1 Invoke the Handover

```bash
peer chaincode invoke ... -c '{"function":"RecordHandover","Args":["shipment-001","Port A","Port B","IN_TRANSIT","2026-07-12T12:00:00Z"]}'
```

**Result:** `status:200` – success.

### 5.2 Query the Record

```bash
peer chaincode query -C mychannel -n slave-twin -c '{"function":"ReadHandover","Args":["shipment-001"]}'
```

**Output:**

```json
{"shipment_id":"shipment-001","origin":"Port A","destination":"Port B","status":"IN_TRANSIT","timestamp":"2026-07-12T12:00:00Z"}
```

---

## 6. Current State

- **Chaincode name:** `slave-twin`
- **Version:** 2.2
- **Sequence:** 4
- **Channel:** `mychannel`
- **Functions available:**
  - `RecordTelemetry` / `ReadTwin` – Slave (IoT)
  - `RecordHandover` / `ReadHandover` – Master (Logistics)

Both contracts coexist in the same package, achieving **logical Master‑Slave segregation** on a single channel.

---

## 7. Key Lessons Learned

1. **Path matters:** The `--path` argument in `peer lifecycle chaincode package` is relative to the current working directory. Always verify the packaged source before deployment.
2. **Version increments are necessary:** To force a rebuild, you must increment the version **and** sequence. Fabric will not rebuild a container for the same package ID.
3. **Verify before committing:** Extract the `.tar.gz` package and check the source before installing it on peers. This would have saved several deployment attempts.
4. **Sequence must be strictly increasing:** Using sequence 3 for version 2.1 and sequence 4 for version 2.2 was correct – Fabric enforces monotonic sequence numbers.

---

## 8. Next Steps (M5)

Now that both Master and Slave contracts are deployed:
- Extend the Go API bridge to support `/api/logistics`.
- Build a Python script to ingest the logistics dataset.
- Run combined stress tests with both sensor and logistics traffic.
- Measure TPS, latency, resource usage, and simulate Raft crash recovery.