package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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

type Record struct {
	DroneID   string `json:"droneID"`
	Zip       string `json:"zip"`
	FlyTime   string `json:"flyTime"`
	FlyRecord string `json:"flyRecord"`
	Reserved  string `json:"reserved"`
}

func ConvertToRFC3339(unixtime string) string {
	unixtimeInt, _ := strconv.ParseInt(unixtime, 10, 64)
	t := time.Unix(unixtimeInt, 0)
	return t.Format(time.RFC3339)
}

func ConvertToUnixTime(datetime string) string {
	t, _ := time.Parse(time.RFC3339, datetime)
	return strconv.FormatInt(t.Unix(), 10)
}

func DecompressRecord(returnedRecord string) string {
	// return returnedRecord
	// Parse the JSON string to
	var records []Record
	var rawRecord RawRecord
	var rawRecords []RawRecord
	json.Unmarshal([]byte(returnedRecord), &records)

	for _, record := range records {
		// skip if record.FlyTime == "-1"
		if record.DroneID == "" {
			continue
		}
		// fmt.Println(record)
		decryptedFlyRecord, _ := Decrypt(record.FlyRecord)
		// fmt.Println(decryptedFlyRecord)
		// split the decryptedFlyRecord by ","
		stringSlice := strings.Split(decryptedFlyRecord, ",")
		// fmt.Println(stringSlice)
		rawRecord = RawRecord{
			Drone:                    stringSlice[0],
			Story:                    stringSlice[1],
			Zip:                      stringSlice[2],
			Datetime:                 stringSlice[3],
			Temperature:              stringSlice[4],
			Wind:                     stringSlice[5],
			Gust:                     stringSlice[6],
			Timesincelastmaintenance: stringSlice[7],
			Flighthours:              stringSlice[8],
			Pitch:                    stringSlice[9],
			Roll:                     stringSlice[10],
			Yaw:                      stringSlice[11],
			Vibex:                    stringSlice[12],
			Vibey:                    stringSlice[13],
			Vibez:                    stringSlice[14],
			Nsat:                     stringSlice[15],
			Noise:                    stringSlice[16],
			Currentslope:             stringSlice[17],
			Brownout:                 stringSlice[18],
			Batterylevel:             stringSlice[19],
			Crash:                    stringSlice[20],
		}

		rawRecords = append(rawRecords, rawRecord)
	}

	// fmt.Println(rawRecords)

	rawRecordsJSON, _ := json.Marshal(rawRecords)
	return string(rawRecordsJSON)
}

func CompressRecord(rawRecord *RawRecord) *Record {
	recordString := strings.Join([]string{
		rawRecord.Drone,
		rawRecord.Story,
		rawRecord.Zip,
		rawRecord.Datetime,
		rawRecord.Temperature,
		rawRecord.Wind,
		rawRecord.Gust,
		rawRecord.Timesincelastmaintenance,
		rawRecord.Flighthours,
		rawRecord.Pitch,
		rawRecord.Roll,
		rawRecord.Yaw,
		rawRecord.Vibex,
		rawRecord.Vibey,
		rawRecord.Vibez,
		rawRecord.Nsat,
		rawRecord.Noise,
		rawRecord.Currentslope,
		rawRecord.Brownout,
		rawRecord.Batterylevel,
		rawRecord.Crash}, ",")
	encryptedFlyRecord, _ := Encrypt(recordString)
	droneRecord := Record{
		DroneID:   rawRecord.Drone,
		Zip:       rawRecord.Zip,
		FlyTime:   ConvertToUnixTime(rawRecord.Datetime),
		FlyRecord: encryptedFlyRecord,
		Reserved:  "",
	}
	return &droneRecord
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

	// Convert the CSV records
	for i, record := range records {
		if i == 0 {
			continue
		}

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
