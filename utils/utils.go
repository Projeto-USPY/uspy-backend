// package utils contain useful helper functions
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	JupiterURL = "https://uspdigital.usp.br/jupiterweb/"
)

// LoadJSON loads json file into data interface
func LoadJSON(filename string, into interface{}) (err error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(bytes), into)

	if err != nil {
		return err
	}

	return nil
}

// GenerateJSON creates json file inside given folder from data struct
func GenerateJSON(data interface{}, folder string, filename string) error {
	bytes, err := json.MarshalIndent(&data, "", "\t")

	if err != nil {
		return err
	}

	_ = ioutil.WriteFile(folder+filename, bytes, 0644)
	return nil
}

func CheckPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckResponse(res *http.Response) {
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Status code error: %d %s\n", res.StatusCode, res.Status)
	}
}

// Returns HTTP response and io.Reader from http.Get, which should substitute http.Body, so characters are read with UTF-8 encoding
// Already panics if error, remember to close response.Body
func HTTPGetWithUTF8(url string) (*http.Response, io.Reader) {
	resp, err := http.Get(url)

	CheckPanic(err)
	CheckResponse(resp)

	reader, err := charset.NewReader(resp.Body, resp.Header["Content-Type"][0])

	CheckPanic(err)

	return resp, reader
}

// Returns HTTP response and io.Reader from http.Post, which should substitute http.Body, so characters are read with UTF-8 encoding
// Already panics if error, remember to close response.Body
func HTTPPostWithUTF8(url string, values url.Values) (*http.Response, io.Reader) {
	resp, err := http.PostForm(url, values)

	CheckPanic(err)
	CheckResponse(resp)

	reader, err := charset.NewReader(resp.Body, resp.Header["Content-Type"][0])

	CheckPanic(err)

	return resp, reader
}

// Encrypts using AES. keyString must be 128, 196 or 256 bits.
func AESEncrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

// Encrypts using AES. keyString must be 128, 196 or 256 bits.
func AESDecrypt(encryptedString string, keyString string) (decryptedString string) {
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
