package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// SensorPayload matches JSON from the sensor ingestion script
type SensorPayload struct {
	SensorID  string  `json:"sensor_id"`
	TempC     float64 `json:"temp_c"`
	Humidity  float64 `json:"humidity"`
	Timestamp string  `json:"timestamp"`
}

// LogisticsPayload matches JSON from the logistics ingestion script
type LogisticsPayload struct {
	ShipmentID  string `json:"shipment_id"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Status      string `json:"status"`
	Timestamp   string `json:"timestamp"`
}

var (
	channelName   = "mychannel"
	chaincodeName = "slave-twin"
)

func loadCertificate(path string) (*x509.Certificate, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	return x509.ParseCertificate(block.Bytes)
}

func loadPrivateKey(path string) (identity.Sign, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	}
	return identity.NewPrivateKeySign(privateKey)
}

func newGateway() (*client.Gateway, *grpc.ClientConn, error) {
	certPath := "wallet/org1/admin/Admin@org1.example.com-cert.pem"
	keyPath := "wallet/org1/admin/priv_sk"
	peerTLSCACert := "../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"

	// Load identity
	cert, err := loadCertificate(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load certificate: %w", err)
	}
	id, err := identity.NewX509Identity("Org1MSP", cert)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create identity: %w", err)
	}

	sign, err := loadPrivateKey(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Create TLS credentials for the peer connection
	peerCertPool := x509.NewCertPool()
	peerCertPEM, err := os.ReadFile(peerTLSCACert)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read peer TLS CA: %w", err)
	}
	if !peerCertPool.AppendCertsFromPEM(peerCertPEM) {
		return nil, nil, fmt.Errorf("failed to parse peer TLS CA PEM")
	}

	peerEndpoint := "localhost:7051"

	// Create a gRPC connection to the peer
	peerConn, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(peerCertPool, "")))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial peer: %w", err)
	}

	// Create gateway connecting ONLY to the peer
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(peerConn),
	)
	if err != nil {
		peerConn.Close()
		return nil, nil, fmt.Errorf("failed to create gateway: %w", err)
	}
	return gw, peerConn, nil
}

func handleSensor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	var payload SensorPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	gw, conn, err := newGateway()
	if err != nil {
		log.Printf("Gateway error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer gw.Close()
	defer conn.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	_, err = contract.SubmitTransaction(
		"RecordTelemetry",
		payload.SensorID,
		fmt.Sprintf("%.2f", payload.TempC),
		fmt.Sprintf("%.2f", payload.Humidity),
		payload.Timestamp,
	)
	if err != nil {
		log.Printf("Submit failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"committed"}`))
}

func handleLogistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	var payload LogisticsPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	gw, conn, err := newGateway()
	if err != nil {
		log.Printf("Gateway error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer gw.Close()
	defer conn.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	_, err = contract.SubmitTransaction(
		"RecordHandover",
		payload.ShipmentID,
		payload.Origin,
		payload.Destination,
		payload.Status,
		payload.Timestamp,
	)
	if err != nil {
		log.Printf("Submit failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"committed"}`))
}

func main() {
	http.HandleFunc("/api/sensor", handleSensor)
	http.HandleFunc("/api/logistics", handleLogistics)
	log.Println("Gateway API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
