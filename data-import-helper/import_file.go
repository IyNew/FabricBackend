package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	RecordID  string `json:"recordID"`
	DroneID   string `json:"droneID"`
	Zip       string `json:"zip"`
	FlyTime   string `json:"flyTime"`
	FlyRecord string `json:"flyRecord"`
	Reserved  string `json:"reserved"`
}

const filePath = "./ds1.csv"

func importFromFile() ([]Record, error) {

	// Open the CSV file
	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer csvFile.Close()

	// Parse the CSV file
	reader := csv.NewReader(csvFile)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var droneRecords []Record

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
		droneRecords = append(droneRecords, droneRecord)
		// fmt.Println(i, droneRecord)
		// createDroneRecordAsync(droneRecord.DroneID, droneRecord.Zip, droneRecord.FlyTime, droneRecord.FlyRecord, droneRecord.Reserved)
	}
	return droneRecords, err
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

const key = "thisis32bitlongpassphraseimusing"

// For Testing purposes
func Encrypt(plaintext string) (string, error) {
	return encryptWithKey(plaintext, key)
}

func Decrypt(ciphertext string) (string, error) {
	return decryptWithKey(ciphertext, key)
}

func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

func encryptWithKey(plainText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	data := PKCS7Padding([]byte(plainText), block.BlockSize())
	ciphertext := make([]byte, block.BlockSize()+len(data))
	iv := ciphertext[:block.BlockSize()]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[block.BlockSize():], data)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptWithKey(cryptoText, key string) (string, error) {
	ciphertext, _ := base64.StdEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < block.BlockSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:block.BlockSize()]
	ciphertext = ciphertext[block.BlockSize():]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	ciphertext = PKCS7UnPadding(ciphertext)
	return string(ciphertext), nil
}

func main() {
	records, _ := importFromFile()

	for _, record := range records {
		fmt.Println(record)
	}
}
