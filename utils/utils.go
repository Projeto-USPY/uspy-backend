package utils

import (
	"encoding/json"
	"io/ioutil"
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

	ioutil.WriteFile(folder+filename, bytes, 0644)
	return nil
}
