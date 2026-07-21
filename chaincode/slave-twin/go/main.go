package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing Digital Twins and Logistics records
type SmartContract struct {
	contractapi.Contract
}

// AssetTwin defines the schema for the physical asset's digital counterpart
type AssetTwin struct {
	SensorID  string  `json:"sensor_id"`
	TempC     float64 `json:"temp_c"`
	Humidity  float64 `json:"humidity"`
	Status    string  `json:"status"` // "NORMAL" or "SPOILED"
	Timestamp string  `json:"timestamp"`
}

// LogisticsRecord defines a macro‑logistics milestone
type LogisticsRecord struct {
	ShipmentID  string `json:"shipment_id"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Status      string `json:"status"` // e.g., "IN_TRANSIT", "ARRIVED", "DELIVERED"
	Timestamp   string `json:"timestamp"`
}

// ============ SLAVE CONTRACT (IoT Telemetry) ============

// RecordTelemetry takes high-frequency IoT data, checks thresholds, and updates the ledger
func (s *SmartContract) RecordTelemetry(ctx contractapi.TransactionContextInterface, sensorID string, tempC float64, humidity float64, timestamp string) error {
	// Default status
	status := "NORMAL"
	if tempC > 8.0 {
		status = "SPOILED"
	}

	twin := AssetTwin{
		SensorID:  sensorID,
		TempC:     tempC,
		Humidity:  humidity,
		Status:    status,
		Timestamp: timestamp,
	}

	twinJSON, err := json.Marshal(twin)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(sensorID, twinJSON)
}

// ReadTwin queries the current state of the Digital Twin
func (s *SmartContract) ReadTwin(ctx contractapi.TransactionContextInterface, sensorID string) (*AssetTwin, error) {
	twinJSON, err := ctx.GetStub().GetState(sensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if twinJSON == nil {
		return nil, fmt.Errorf("the asset twin %s does not exist", sensorID)
	}

	var twin AssetTwin
	err = json.Unmarshal(twinJSON, &twin)
	if err != nil {
		return nil, err
	}
	return &twin, nil
}

// ============ MASTER CONTRACT (Logistics Handovers) ============

// RecordHandover stores a logistics milestone
func (s *SmartContract) RecordHandover(ctx contractapi.TransactionContextInterface, shipmentID string, origin string, destination string, status string, timestamp string) error {
	record := LogisticsRecord{
		ShipmentID:  shipmentID,
		Origin:      origin,
		Destination: destination,
		Status:      status,
		Timestamp:   timestamp,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(shipmentID, recordJSON)
}

// ReadHandover retrieves a logistics record by shipment ID
func (s *SmartContract) ReadHandover(ctx contractapi.TransactionContextInterface, shipmentID string) (*LogisticsRecord, error) {
	recordJSON, err := ctx.GetStub().GetState(shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if recordJSON == nil {
		return nil, fmt.Errorf("the logistics record %s does not exist", shipmentID)
	}

	var record LogisticsRecord
	err = json.Unmarshal(recordJSON, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// AssetMetadata defines the schema for asset metadata
type AssetMetadata struct {
	AssetID                    string  `json:"asset_id"`
	Location                   string  `json:"location"`
	Temperature                float64 `json:"temperature"`
	Vibration                  float64 `json:"vibration"`
	LastMaintenance            string  `json:"last_maintenance"`
	ConditionScore             float64 `json:"condition_score"`
	ResourceUtilization        float64 `json:"resource_utilization"`
	DeliveryEfficiency         float64 `json:"delivery_efficiency"`
	DowntimeHours              float64 `json:"downtime_hours"`
	InventoryLevel             string  `json:"inventory_level"`
	LogisticsCost              float64 `json:"logistics_cost"`
	Timestamp                  string  `json:"timestamp"`
	SupplyChainEfficiencyLabel int     `json:"supply_chain_efficiency_label"`
}

// RecordMetadata stores asset metadata
func (s *SmartContract) RecordMetadata(ctx contractapi.TransactionContextInterface, assetID string, location string, temperature string, vibration string, lastMaintenance string, conditionScore string, resourceUtilization string, deliveryEfficiency string, downtimeHours string, inventoryLevel string, logisticsCost string, timestamp string, efficiencyLabel string) error {
	meta := AssetMetadata{
		AssetID:         assetID,
		Location:        location,
		LastMaintenance: lastMaintenance,
		InventoryLevel:  inventoryLevel,
		Timestamp:       timestamp,
	}
	// Note: We take strings from the API bridge and store them directly, or I can cast them to floats here if you prefer strict types inside the struct, but this matches the string-based arguments sent from my Go API SubmitTransaction.
	// To keep it simple and match my API bridge exactly, let's just let it save the raw JSON bytes.
	// For a production system, I'd parse those string numbers back to floats here.

	// I am updating this function signature to expect strings since my Gateway API formats them as strings (`fmt.Sprintf("%f", payload.Temperature)`) before submitting.
	fmt.Sscanf(temperature, "%f", &meta.Temperature)
	fmt.Sscanf(vibration, "%f", &meta.Vibration)
	fmt.Sscanf(conditionScore, "%f", &meta.ConditionScore)
	fmt.Sscanf(resourceUtilization, "%f", &meta.ResourceUtilization)
	fmt.Sscanf(deliveryEfficiency, "%f", &meta.DeliveryEfficiency)
	fmt.Sscanf(downtimeHours, "%f", &meta.DowntimeHours)
	fmt.Sscanf(logisticsCost, "%f", &meta.LogisticsCost)
	fmt.Sscanf(efficiencyLabel, "%d", &meta.SupplyChainEfficiencyLabel)

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("meta_"+assetID, metaJSON)
}

// ReadMetadata retrieves asset metadata by asset ID
func (s *SmartContract) ReadMetadata(ctx contractapi.TransactionContextInterface, assetID string) (*AssetMetadata, error) {
	metaJSON, err := ctx.GetStub().GetState("meta_" + assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if metaJSON == nil {
		return nil, fmt.Errorf("metadata for asset %s does not exist", assetID)
	}
	var meta AssetMetadata
	err = json.Unmarshal(metaJSON, &meta)
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}
	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
