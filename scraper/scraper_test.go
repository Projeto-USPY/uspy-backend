package scraper

import (
	"fmt"
	"testing"
)

func TestScrapeDepartments(t *testing.T) {
	result := ScrapeDepartments()
	fmt.Println(result)
}

/*
func TestScrapeSubjectDescription(t *testing.T) {
	var course string = "SCC0200"
	desc, err := ScrapeSubjectDescription(course)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(desc)
	}

	course = "SCC1234"
	desc, err = ScrapeSubjectDescription(course)
	if err != nil {
		t.Log(err)
	} else {
		t.Fail()
	}
}
*/

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
