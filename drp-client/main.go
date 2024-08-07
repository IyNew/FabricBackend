package main

import (
	"log"
	"net/http"

	// "drp-client/handlers"
	// "encoding/json"
	// "bytes"
	// "context"
	// "crypto/x509"
	// "errors"
	// "fmt"
	// "os"
	// "path"
	// "time"

	"github.com/IyNew/FabricBackend/drp-client/utils"

	"github.com/gorilla/mux"
	// "github.com/hyperledger/fabric-gateway/pkg/client"
	// "github.com/hyperledger/fabric-gateway/pkg/identity"
	// "github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/status"
)

func main() {
	r := mux.NewRouter()

	utils.InitGateway()
	defer utils.ClientConn.Close()
	defer utils.GatewayConn.Close()

	// initialize the ledger with data from the CSV file
	// utils.ImportFromFile("./ds1.csv")
	// utils.GetAllRecords()

	// Define routes
	r.HandleFunc("/api/hello", utils.HelloHandler).Methods("GET")
	// r.HandleFunc("/api/test", utils.TestAPI).Methods("GET")
	// r.HandleFunc("/api/init", utils.InitLedger)
	r.HandleFunc("/api/record/all", utils.GetAllRecords).Methods("GET")
	r.HandleFunc("/api/record/{droneID}", utils.GetAllRecordsForOneDrone).Methods("GET")
	r.HandleFunc("/api/record/create", utils.CreateRecord).Methods("POST")

	// Test api for querying records with selector
	r.HandleFunc("/api/record/testSelector", utils.TestQueryRecordsWithSelector).Methods("GET")

	// api for querying records with selector, with string selector in JSON format
	r.HandleFunc("/api/record/{selectorString}", utils.QueryRecordsWithSelectorJSON).Methods("GET")

	// r.HandleFunc("/api/users", utils.GetUsers).Methods("GET")
	// r.HandleFunc("/api/users/{id}", utils.GetUser).Methods("GET")
	// r.HandleFunc("/api/users", utils.CreateUser).Methods("POST")

	// Start server
	log.Println("Server listening on port 6999")
	if err := http.ListenAndServe(":6999", r); err != nil {
		log.Fatal(err)
	}

}
