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

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}
	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
