package scraper

import (
	"testing"

	"github.com/tpreischadt/ProjetoJupiter/utils"
)

func TestScrapeICMC(t *testing.T) {
	courses, err := ScrapeICMCCourses()

	if err != nil {
		t.Fail()
	}

	err = utils.GenerateJSON(courses, "../data/", "courses.json")

	if err != nil {
		t.Fatal(err)
	}
}
