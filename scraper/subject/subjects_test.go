package subject

import (
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/entity"
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
