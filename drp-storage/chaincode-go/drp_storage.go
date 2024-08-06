package main

import (
	"encoding/json"
	"fmt"
	"log"

	// Add this import statement
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// the index is
const index = "id~time"

type SimpleChaincode struct {
	contractapi.Contract
}

// Structure of the flying record, only the droneID, droneZip and flytime are stored in the ledger for query
// The rest are serialized and encrypted in the flyrecord field
// The reserved field is reserved for future use
type Record struct {
	RecordID  string `json:"recordID"`
	DroneID   string `json:"droneID"`
	Zip       string `json:"zip"`
	FlyTime   string `json:"flyTime"`
	FlyRecord string `json:"flyRecord"`
	Reserved  string `json:"reserved"`
}

type PaginatedQueryResult struct {
	Records             []Record `json:"records"`
	FetchedRecordsCount int32    `json:"fetchedRecordsCount"`
	Bookmark            string   `json:"bookmark"`
}

// CreateRecord adds a new record to the world state with given details
func (s *SimpleChaincode) CreateRecord(ctx contractapi.TransactionContextInterface, droneID string, zip string, flytime string, flyrecord string, reserved string) error {

	// the recordID is a combination of the droneID_flytime
	recordID := droneID + "_" + flytime

	// Check if the record already exists
	exists, err := s.RecordExists(ctx, recordID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if exists {
		return fmt.Errorf("the record %s already exists", recordID)
	}

	record := Record{
		RecordID:  recordID,
		DroneID:   droneID,
		Zip:       zip,
		FlyTime:   flytime,
		FlyRecord: flyrecord,
		Reserved:  reserved,
	}
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(recordID, recordJSON)
}

// InitLedger adds a base set of records to the ledger
func (s *SimpleChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	records := []Record{
		{DroneID: "drone1", Zip: "10001", FlyTime: "100", FlyRecord: "record1", Reserved: "reserved1"},
		{DroneID: "drone2", Zip: "10002", FlyTime: "200", FlyRecord: "record2", Reserved: "reserved2"},
		{DroneID: "drone3", Zip: "10003", FlyTime: "300", FlyRecord: "record3", Reserved: "reserved3"},
		{DroneID: "drone4", Zip: "10004", FlyTime: "400", FlyRecord: "record4", Reserved: "reserved4"},
		{DroneID: "drone5", Zip: "10005", FlyTime: "500", FlyRecord: "record5", Reserved: "reserved5"},
	}

	// records, _ := importFromFile()

	for _, record := range records {
		err := s.CreateRecord(ctx, record.DroneID, record.Zip, record.FlyTime, record.FlyRecord, record.Reserved)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAllRecords returns all records found in world state
func (s *SimpleChaincode) GetAllRecords(ctx contractapi.TransactionContextInterface) ([]Record, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []Record
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var record Record
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// RecordExists returns true when record with given ID exists in world state
func (s *SimpleChaincode) RecordExists(ctx contractapi.TransactionContextInterface, recordID string) (bool, error) {
	recordJSON, err := ctx.GetStub().GetState(recordID)
	if err != nil {
		return false, fmt.Errorf("failed to read record %s from world state: %v", recordID, err)
	}

	return recordJSON != nil, nil
}

// constructQueryResponseFromIterator constructs a slices of Records from QueryResultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Record, error) {
	var records []*Record
	for resultsIterator.HasNext() {
		recordResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var record Record
		err = json.Unmarshal(recordResponse.Value, &record)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}

// GetRecordByRange performs a range query based on the start and end keys provided.
func (s *SimpleChaincode) GetRecordByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Record, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

//  getQueryResultForQueryString queries for records based on a passed in query string.
//  This is only supported for couchdb
// func (s *SimpleChaincode) getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Record, error) {
// 	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	return constructQueryResponseFromIterator(resultsIterator)
// }

func main() {
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset chaincode: %v", err)
	}
}
