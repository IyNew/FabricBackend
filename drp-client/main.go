package main

import (
	"log"
	"net/http"
	// "drp-client/handlers"
	"encoding/json"
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/IyNew/FabricBackend/drp-client/utils"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)


// type User struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// }

// var users = []User{
// 	{ID: "1", Name: "John Doe"},
// 	{ID: "2", Name: "Jane Doe"},
// }

// func HelloHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Hello, world!"))
// }

// func GetUsers(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(users)
// }

// func GetUser(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	for _, user := range users {
// 		if user.ID == params["id"] {
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusOK)
// 			json.NewEncoder(w).Encode(user)
// 			return
// 		}
// 	}
// 	w.WriteHeader(http.StatusNotFound)
// 	w.Write([]byte("User not found"))
// }

// func CreateUser(w http.ResponseWriter, r *http.Request) {
// 	var newUser User
// 	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	users = append(users, newUser)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(newUser)
// }


func main() {
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/api/hello", HelloHandler).Methods("GET")
	r.HandleFunc("/api/users", GetUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/api/users", CreateUser).Methods("POST")

	// Start server
	log.Println("Server listening on port 6999")
	if err := http.ListenAndServe(":6999", r); err != nil {
		log.Fatal(err)
	}
}