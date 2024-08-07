package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// ========================================
// For testing purposes
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = []User{
	{ID: "1", Name: "John Doe"},
	{ID: "2", Name: "Jane Doe"},
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, world!"))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, user := range users {
		if user.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("User not found"))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	users = append(users, newUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// End of testing purposes
// ========================================

// func TestAPI(w http.ResponseWriter, r *http.Request) {
// 	initLedger()
// 	fmt.Println("Test initLedger API completed.")
// 	getAllAssets()
// 	fmt.Println("Test getAllAssets API completed.")
// 	createAsset()
// 	fmt.Println("Test createAsset API completed.")
// 	getAllAssets()
// 	fmt.Println("Test getAllAssets API completed.")
// }

func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(DecompressRecord(getAllRecords())))
	// json.NewEncoder(w).Encode(records)
}

func InitLedger(w http.ResponseWriter, r *http.Request) {
	initDroneRecords()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ledger initialized"))
}

func GetAllRecordsForOneDrone(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(DecompressRecord(getAllRecordsForOneDrone(params["droneID"]))))
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	// Post request body
	var rawRecord RawRecord
	if err := json.NewDecoder(r.Body).Decode(&rawRecord); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newRecord := CompressRecord(&rawRecord)
	createDroneRecord(newRecord.DroneID, newRecord.Zip, newRecord.FlyTime, newRecord.FlyRecord, newRecord.Reserved)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(newRecord)
}

func TestQueryRecordsWithSelector(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	selector := `{"selector":{"droneID":"blue"}, "use_index":["indexDroneTimeDoc","indexDroneTime"], "sort":[{"flyTime":"desc"}]}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(DecompressRecord(getRecordWithSelector(selector))))
}

// QueryRecordsWithSelectorJSON is a handler for querying records with selector in JSON format
// The selector is passed as a URL parameter
func QueryRecordsWithSelectorJSON(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// selector params["slectorBody"] in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(DecompressRecord(getRecordWithSelector(params["selectorString"]))))
}
