package depscraper

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func checkPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkResponse(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s\n", res.StatusCode, res.Status)
	}
}

func getProfessors(dep *string, page int) []string {
	icmcURL := "https://www.icmc.usp.br/templates/icmc2015/php/pessoas.php"
	formData := url.Values{
		"grupo":  {"Docente"},
		"depto":  {*dep},
		"nome":   {""},
		"pagina": {strconv.Itoa(page)},
	}

	response, err := http.PostForm(icmcURL, formData)
	checkPanic(err)
	checkResponse(response)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(response)
	checkPanic(err)

	results := make([]string, 0, 100)

	doc.Find(".caption").Each(func(i int, s *goquery.Selection) {
		prof := strings.TrimSpace(s.Text())
		results = append(results, prof)
	})

	return results
}

// Scrape scrapes the professors page
func Scrape() *map[string][]string {
	deps := []string{"SCC", "SMA", "SME", "SSC"}
	results := make(map[string][]string)

	for _, dep := range deps {
		i := 1
		for {
			profs := getProfessors(&dep, i)

			if len(profs) == 0 {
				break
			} else {
				results[dep] = append(results[dep], profs...)
				i++
			}
		}
	}

	return &results
}
