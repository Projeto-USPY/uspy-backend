package utils

import (
	"encoding/json"
	"io/ioutil"
)

// GenerateJSON creates json file inside given folder from data struct
func GenerateJSON(data interface{}, folder string, filename string) error {
	bytes, err := json.MarshalIndent(&data, "", "\t")

	if err != nil {
		return err
	}

	ioutil.WriteFile(folder+filename, bytes, 0644)
	return nil
}
