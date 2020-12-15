package utils

import (
	"encoding/json"
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
