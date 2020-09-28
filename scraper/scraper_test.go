package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestScrapeDepartments(t *testing.T) {
	result := ScrapeDepartments()
	fmt.Println(result)
}

func _UnescapeUnicodeCharactersInJSON(_jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(_jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
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
		fixedBytes, _ := _UnescapeUnicodeCharactersInJSON(bytes)
		ioutil.WriteFile("courses.json", fixedBytes, 6444)
	}
}
