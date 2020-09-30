package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
)

func TestScrapeDepartments(t *testing.T) {
	result := ScrapeDepartments()
	fmt.Println(result)
}

func TestScrapeSubject(t *testing.T) {
	wg := &sync.WaitGroup{}
	c := make(chan Subject, 200)

	wg.Add(1)
	go scrapeSubject(`/obterDisciplina?sgldis=SME0130&codcur=55041&codhab=0`, `55041`, c, wg)

	subj := <-c

	fmt.Printf("%+v", subj)
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
