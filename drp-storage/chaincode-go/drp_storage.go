package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

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
// if RecordID = droneID, the record is a meta record for the drone, FlyRecord holds the number of records for the drone, the rest are set empty for now, including DroneID
// if RecordID = droneID_(digit), the record is a flying record for the drone, FlyRecord holds the encrypted flying record, the rest are set empty for now
type Record struct {
	RecordID  string `json:"recordID"`
	DroneID   string `json:"droneID"`
	Zip       string `json:"zip"`
	FlyTime   int64  `json:"flyTime"`
	FlyRecord string `json:"flyRecord"`
	Reserved  string `json:"reserved"`
}

type PaginatedQueryResult struct {
	Records             []Record `json:"records"`
	FetchedRecordsCount int32    `json:"fetchedRecordsCount"`
	Bookmark            string   `json:"bookmark"`
}

// Hello returns a greeting message to check if the chaincode is alive
func (s *SimpleChaincode) Hello(ctx contractapi.TransactionContextInterface) string {
	return "Hello from fabric, the service is running!"
}

// ReadRecord returns the record with the given recordID
func (s *SimpleChaincode) ReadRecord(ctx contractapi.TransactionContextInterface, recordID string) (*Record, error) {
	recordJSON, err := ctx.GetStub().GetState(recordID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if recordJSON == nil {
		return nil, fmt.Errorf("the record %s does not exist", recordID)
	}

	var record Record
	err = json.Unmarshal(recordJSON, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// // CreateMetaRecord adds a new meta rocord to the drone with droneID
// func (s *SimpleChaincode) CreateMetaRecord(ctx contractapi.TransactionContextInterface, droneID string) error {

// 	metaRecordJSON, _ := json.Marshal(metaRcord)

// 	return ctx.GetStub().PutState(droneID, metaRecordJSON)
// }

// GetMetaRecord returns the meta record for the drone with droneID
func (s *SimpleChaincode) GetMetaRecord(ctx contractapi.TransactionContextInterface, droneID string) (*Record, error) {
	// check if the meta record exists
	exists, _ := s.RecordExists(ctx, droneID)
	var metaRecord Record
	if !exists {
		metaRecord = Record{
			RecordID:  droneID,
			DroneID:   "",
			Zip:       "",
			FlyTime:   -1,
			FlyRecord: "0",
			Reserved:  "",
		}
	} else {
		metaRecordJSON, err := ctx.GetStub().GetState(droneID)
		if err != nil {
			return nil, fmt.Errorf("failed to get meta record for %s: %v", droneID, err)
		}
		fmt.Println("found MetaRecordJSON: ", metaRecordJSON)
		err = json.Unmarshal(metaRecordJSON, &metaRecord)
		if err != nil {
			return nil, err
		}
	}

	return &metaRecord, nil
}

func MD5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// CreateRecord adds a new record to the world state with given details
func (s *SimpleChaincode) CreateRecord(ctx contractapi.TransactionContextInterface, droneID string, zip string, flytime string, flyrecord string, reserved string) error {

	// ========================================
	// // Remove the meta record function to avoid read conflict
	// var nextFlag int
	// var metaRecord *Record
	// metaRecord, _ = s.GetMetaRecord(ctx, droneID)
	// fmt.Println("Got MetaRecord", metaRecord)

	// nextFlag, _ = strconv.Atoi(metaRecord.FlyRecord)
	// nextFlag++

	// // the recordID is a combination of the droneID_flytime
	// recordID := droneID + "_" + strconv.Itoa(nextFlag)
	// ========================================

	recordID := droneID + "_" + MD5Hash(flyrecord)
	// Check if the record already exists
	exists, err := s.RecordExists(ctx, recordID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if exists {
		return fmt.Errorf("the record %s already exists", recordID)
	}

	flyTimeInt, err := strconv.ParseInt(flytime, 10, 64)
	if err != nil {
		return err
	}

	record := Record{
		RecordID:  recordID,
		DroneID:   droneID,
		Zip:       zip,
		FlyTime:   flyTimeInt,
		FlyRecord: flyrecord,
		Reserved:  reserved,
	}
	// print the record to be created
	fmt.Println("Record to be created", record)

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// create the record
	err = ctx.GetStub().PutState(recordID, recordJSON)
	if err != nil {
		fmt.Println("Error in creating record for", recordID, err)
		return err
	}

	// ========================================
	// Remove the meta record function to avoid read conflict
	// metaRecord.FlyRecord = strconv.Itoa(nextFlag)
	// metaRecordJSON, _ := json.Marshal(metaRecord)
	// err = ctx.GetStub().PutState(droneID, metaRecordJSON)
	// if err != nil {
	// 	fmt.Println("Error in updating meta record for", droneID, err)
	// 	return err
	// }
	// ========================================

	return nil
}

// InitLedger adds a base set of records to the ledger, not used in the current implementation
func (s *SimpleChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	records := []Record{
		{DroneID: "drone1", Zip: "10001", FlyTime: 100, FlyRecord: "record1", Reserved: "reserved1"},
		{DroneID: "drone2", Zip: "10002", FlyTime: 200, FlyRecord: "record2", Reserved: "reserved2"},
		{DroneID: "drone3", Zip: "10003", FlyTime: 300, FlyRecord: "record3", Reserved: "reserved3"},
		{DroneID: "drone4", Zip: "10004", FlyTime: 400, FlyRecord: "record4", Reserved: "reserved4"},
		{DroneID: "drone5", Zip: "10005", FlyTime: 500, FlyRecord: "record5", Reserved: "reserved5"},
	}

	// records, _ := importFromFile()

	for _, record := range records {
		err := s.CreateRecord(ctx, record.DroneID, record.Zip, strconv.FormatInt(record.FlyTime, 10), record.FlyRecord, record.Reserved)
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

// getQueryResultForQueryString queries for records based on a passed in query string.
// This is only supported for couchdb
func (s *SimpleChaincode) getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Record, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// QueryRecords uses a query string to perform a query for records.
func (s *SimpleChaincode) QueryRecords(ctx contractapi.TransactionContextInterface, queryString string) ([]*Record, error) {
	return s.getQueryResultForQueryString(ctx, queryString)
}

// QueryRecordsByDroneID queries for records based on a passed in droneID.
func (s *SimpleChaincode) QueryRecordsByDroneID(ctx contractapi.TransactionContextInterface, droneID string) ([]*Record, error) {
	queryString := fmt.Sprintf(`{"selector":{"droneID":"%s"}}`, droneID)
	return s.getQueryResultForQueryString(ctx, queryString)
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset chaincode: %v", err)
	}
}
