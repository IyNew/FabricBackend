package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// type Record struct {
// 	drone                    string
// 	story                    string
// 	zip                      string
// 	datetime                 string
// 	temperature              string
// 	wind                     string
// 	gust                     string
// 	timesincelastmaintenance string
// 	flighthours              string
// 	pitch                    string
// 	roll                     string
// 	yaw                      string
// 	vibex                    string
// 	vibey                    string
// 	vibez                    string
// 	nsat                     string
// 	noise                    string
// 	currentslope             string
// 	brownout                 string
// 	batterylevel             string
// 	crash                    string
// }

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
