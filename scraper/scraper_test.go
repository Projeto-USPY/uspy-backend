package scraper

import (
	"fmt"
	"sync"
	"testing"
)

func TestScrapeDepartments(t *testing.T) {
	result := ScrapeDepartments()
	GenerateJSON(result, "../data/", "professors.json")

	fmt.Println(result)
}

func TestScrapeSubject(t *testing.T) {
	wg := &sync.WaitGroup{}
	c := make(chan Subject, 200)

	wg.Add(1)
	go scrapeSubject(`/obterDisciplina?sgldis=SME0130&codcur=55041&codhab=0`, `55041`, true, c, wg)

	subj := <-c

	fmt.Printf("%+v", subj)
}

func TestScrapeICMC(t *testing.T) {
	courses, err := ScrapeICMC()

	if err != nil {
		t.Fail()
	}

	GenerateJSON(courses, "../data/", "courses.json")
}
