package scraper

import (
	"fmt"
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

// ScrapeDepartments scrapes the professors page
func ScrapeDepartments() *map[string][]string {
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

// ScrapeCourseDescription scrapes the description of a course
func ScrapeCourseDescription(courseCode string) (string, error) {
	const url string = "https://uspdigital.usp.br/jupiterweb/obterDisciplina?nomdis=&sgldis=%v"
	formattedURL := fmt.Sprintf(url, courseCode) // create course code
	resp, err := http.Get(formattedURL)
	checkResponse(resp)
	checkPanic(err)

	defer resp.Body.Close()

	// Create goquery HTML structure
	doc, err := goquery.NewDocumentFromResponse(resp)
	checkPanic(err)

	// Returns error telling that course is invalid or not yet activated
	if doc.Find("#web_mensagem").Length() > 0 {
		return "", fmt.Errorf("Wasn't able to find course named %v", courseCode)
	}

	// To parse course description, get <b> element with content "Objetivos" and course description will be on next <tr>
	var objetivosNode *goquery.Selection = nil
	doc.Find("b").Each(func(i int, s *goquery.Selection) {
		html, err := s.Html() // get inner html
		checkPanic(err)

		if strings.Trim(html, " ") == "Objetivos" { // found
			objetivosNode = s
		}
	})

	if objetivosNode == nil {
		log.Fatal("Couldn't find node with message: \"Objetivos\"")
	}

	objetivosTr := objetivosNode.Closest("tr") // get tr parent
	descriptionTr := objetivosTr.Next()        // tr with description is next <tr>

	desc := strings.Trim(descriptionTr.Text(), " \n")

	return desc, nil
}
