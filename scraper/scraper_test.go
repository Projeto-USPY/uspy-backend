package scraper

import (
	"net/http"
	"testing"
)

func TestGetProfessorByCode(t *testing.T) {
	prof, status, err := getProfessorByCode(54946) // Janete
	t.Log(prof)
	if err != nil {
		t.Fatal(err, status)
	} else if status != http.StatusOK {
		t.Fatal("status", status)
	}

}

func TestGetProfessorHistory(t *testing.T) {
	results, err := GetProfessorHistory(54946, 2010)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(results)
}

func TestScrapeAllOfferings(t *testing.T) {
	ScrapeAllOfferings()
}
