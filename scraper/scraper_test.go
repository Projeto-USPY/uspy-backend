package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestScrapeDepartments(t *testing.T) {
	result := ScrapeDepartments()
	fmt.Println(result)
}

func TestScrapeICMC(t *testing.T) {
	courses, err := ScrapeICMC()

	if err != nil {
		t.Fail()
	}

	bytes, err := json.MarshalIndent(&courses, "", "\t")

	if err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		ioutil.WriteFile("courses.json", bytes, 0644)
	}
}
