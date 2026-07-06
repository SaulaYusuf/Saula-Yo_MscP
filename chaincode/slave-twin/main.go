package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a Digital Twin
type SmartContract struct {
	contractapi.Contract
}

// AssetTwin defines the schema for the physical asset's digital counterpart
type AssetTwin struct {
	SensorID  string  `json:"sensor_id"`
	TempC     float64 `json:"temp_c"`
	Humidity  float64 `json:"humidity"`
	Status    string  `json:"status"` // e.g., "NORMAL" or "SPOILED"
	Timestamp string  `json:"timestamp"`
}

// RecordTelemetry takes high-frequency IoT data, checks thresholds, and updates the ledger
func (s *SmartContract) RecordTelemetry(ctx contractapi.TransactionContextInterface, sensorID string, tempC float64, humidity float64, timestamp string) error {
	
	// Default status
	status := "NORMAL"

	// THE THRESHOLD ENGINE: If temp crosses 8.0C, flag as spoiled.
	if tempC > 8.0 {
		status = "SPOILED"
	}

	// Create the twin object
	twin := AssetTwin{
		SensorID:  sensorID,
		TempC:     tempC,
		Humidity:  humidity,
		Status:    status,
		Timestamp: timestamp,
	}

	// Convert to JSON and save to the Slave Ledger
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

func main() {
	twinChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating slave-twin chaincode: %v", err)
	}

	if err := twinChaincode.Start(); err != nil {
		log.Panicf("Error starting slave-twin chaincode: %v", err)
	}
}
