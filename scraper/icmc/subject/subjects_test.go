// package subject contains useful functions require to scrape icmc subject data from jupiterweb
package subject

import (
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/utils"
	"sync"
	"testing"
)

func TestScrapeSubject(t *testing.T) {
	wg := &sync.WaitGroup{}
	c := make(chan entity.Subject, 200)

	wg.Add(1)
	go scrapeSubject(`/obterDisciplina?sgldis=SME0130&codcur=55041&codhab=0`, `55041`, true, c, wg)

	subj := <-c

	fmt.Printf("%+v", subj)
}

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
