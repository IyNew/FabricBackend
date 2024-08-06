package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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

type Record struct {
	DroneID   string `json:"droneID"`
	Zip       string `json:"zip"`
	FlyTime   string `json:"flyTime"`
	FlyRecord string `json:"flyRecord"`
	Reserved  string `json:"reserved"`
}

func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(getAllDroneRecords()))
	// json.NewEncoder(w).Encode(records)
}

func InitLedger(w http.ResponseWriter, r *http.Request) {
	initDroneRecords()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ledger initialized"))
}

func ImportFromFile(filePath string) {
	// Open the CSV file
	csvFile, err := os.Open("./ds1.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer csvFile.Close()

	// Parse the CSV file
	reader := csv.NewReader(csvFile)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// var droneRecords []Record

	// Convert the CSV records
	for i, record := range records {
		if i == 0 {
			continue
		}
		// concat record with "," to get the flyRecord
		flyRecord := strings.Join(record, ",")
		encryptedFlyRecord, _ := Encrypt(flyRecord)

		droneRecord := Record{
			DroneID:   record[0],
			Zip:       record[2],
			FlyTime:   ConvertToUnixTime(record[3]),
			FlyRecord: encryptedFlyRecord,
			Reserved:  "",
		}
		fmt.Println(i, droneRecord)
		createDroneRecord(droneRecord.DroneID, droneRecord.Zip, droneRecord.FlyTime, droneRecord.FlyRecord, droneRecord.Reserved)
	}
}
