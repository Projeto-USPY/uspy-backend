package scraper

import (
	"encoding/json"
	"testing"
)

func TestNewInstituteScraper(t *testing.T) {
	instituteSc := NewInstituteScraper("55")
	if institute, err := instituteSc.Start(); err != nil {
		t.Fatal(err)
	} else {
		if bytes, err := json.MarshalIndent(&institute, "", "    "); err != nil {
			t.Fatal(err)
		} else {
			t.Log(string(bytes))
		}
	}
}
