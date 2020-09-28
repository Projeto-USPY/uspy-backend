package scraper

import (
	"fmt"
	"testing"
)

func TestScrapeDepartments(t *testing.T) {
	result := ScrapeDepartments()
	fmt.Println(result)
}
func TestScrapeICMC(t *testing.T) {
	courses, err := ScrapeICMC()

	if err == nil {
		for _, c := range courses {
			fmt.Println(c.name, len(c.subjects))
		}
	} else {
		t.Fail()
	}
}
