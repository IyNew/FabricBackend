package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

const filePath = "./ds1.csv"
const apiURL = "http://localhost:6999/api/record/create"

type Record struct {
	RecordID  string `json:"recordID"`
	DroneID   string `json:"droneID"`
	Zip       string `json:"zip"`
	FlyTime   string `json:"flyTime"`
	FlyRecord string `json:"flyRecord"`
	Reserved  string `json:"reserved"`
}

type RawRecord struct {
	Drone                    string `json:"drone"`
	Story                    string `json:"story"`
	Zip                      string `json:"zip"`
	Datetime                 string `json:"datetime"`
	Temperature              string `json:"temperature"`
	Wind                     string `json:"wind"`
	Gust                     string `json:"gust"`
	Timesincelastmaintenance string `json:"timesincelastmaintenance"`
	Flighthours              string `json:"flighthours"`
	Pitch                    string `json:"pitch"`
	Roll                     string `json:"roll"`
	Yaw                      string `json:"yaw"`
	Vibex                    string `json:"vibex"`
	Vibey                    string `json:"vibey"`
	Vibez                    string `json:"vibez"`
	Nsat                     string `json:"nsat"`
	Noise                    string `json:"noise"`
	Currentslope             string `json:"currentslope"`
	Brownout                 string `json:"brownout"`
	Batterylevel             string `json:"batterylevel"`
	Crash                    string `json:"crash"`
}

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return
	}

	// concurrency control
	var wg sync.WaitGroup
	concurrency := 5
	semaphore := make(chan struct{}, concurrency)

	for i, record := range records {
		// skip the first line
		if i == 0 {
			continue
		}
		wg.Add(1)
		semaphore <- struct{}{}

		go func(record []string) {
			defer wg.Done()
			defer func() {
				<-semaphore
			}()
			fmt.Println(record)
			// recordLine := strings.Split(record, ",")
			// parse the recordLine to RawRecord JSON
			rawRecord := RawRecord{
				Drone:                    record[0],
				Story:                    record[1],
				Zip:                      record[2],
				Datetime:                 record[3],
				Temperature:              record[4],
				Wind:                     record[5],
				Gust:                     record[6],
				Timesincelastmaintenance: record[7],
				Flighthours:              record[8],
				Pitch:                    record[9],
				Roll:                     record[10],
				Yaw:                      record[11],
				Vibex:                    record[12],
				Vibey:                    record[13],
				Vibez:                    record[14],
				Nsat:                     record[15],
				Noise:                    record[16],
				Currentslope:             record[17],
				Brownout:                 record[18],
				Batterylevel:             record[19],
				Crash:                    record[20],
			}

			rawRecordsJSON, _ := json.Marshal(rawRecord)
			resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(rawRecordsJSON))
			if err != nil {
				fmt.Println("Error posting record: ", err)
				return
			}
			defer resp.Body.Close()

			fmt.Println("Record posted successfully")
		}(record)
		// POST the rawRecordsJSON to the server
		// POSTRequest(rawRecordsJSON)
	}
	wg.Wait()
}
