package scraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

const jupiterURL = "https://uspdigital.usp.br/jupiterweb/"

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

// Returns HTTP response and io.Reader from http.Get, which should substitute http.Body, so characters are read with UTF-8 encoding
// Already panics if error, remember to close response.Body
func httpGetWithUTF8(url string) (*http.Response, io.Reader) {
	resp, err := http.Get(url)

	checkPanic(err)
	checkResponse(resp)

	reader, err := charset.NewReader(resp.Body, resp.Header["Content-Type"][0])

	checkPanic(err)

	return resp, reader
}

// Returns HTTP response and io.Reader from http.Post, which should substitute http.Body, so characters are read with UTF-8 encoding
// Already panics if error, remember to close response.Body
func httpPostWithUTF8(url string, values url.Values) (*http.Response, io.Reader) {
	resp, err := http.PostForm(url, values)

	checkPanic(err)
	checkResponse(resp)

	reader, err := charset.NewReader(resp.Body, resp.Header["Content-Type"][0])

	checkPanic(err)

	return resp, reader
}

func getProfessors(dep *string, page int) []string {
	icmcURL := "https://www.icmc.usp.br/templates/icmc2015/php/pessoas.php"
	formData := url.Values{
		"grupo":  {"Docente"},
		"depto":  {*dep},
		"nome":   {""},
		"pagina": {strconv.Itoa(page)},
	}

	response, body := httpPostWithUTF8(icmcURL, formData)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
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

// Subject describes a subject (example: SMA0356 - Cálculo IV)
type Subject struct {
	Code          string
	Name          string
	Description   string
	ClassCredits  int
	AssignCredits int
	TotalHours    string
	Requirements  []string
	Optional      bool
}

// Course represents a course/major (example: BCC)
type Course struct {
	Name     string
	Code     string
	Subjects []Subject
}

func scrapeSubjectNames(doc *goquery.Document) (code, name string, e error) {
	defer func() {
		if r := recover(); r != nil {
			code, name, e = "", "", fmt.Errorf("Error getting subject name or code: %v", r)
		}
	}()

	doc.Find("b").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())

		if strings.HasPrefix(text, "Disciplina:") {
			names := strings.Split(text, "-")

			code = strings.TrimSpace(names[0])
			name = strings.TrimSpace(names[1])

			// Remove "Disciplina:"
			code = strings.TrimSpace(strings.Split(code, ":")[1])
			e = nil

			return
		}
	})

	return code, name, e
}

func scrapeSubjectDescription(doc *goquery.Document) (string, error) {
	// Returns error telling that subject is invalid or not yet activated
	if doc.Find("#web_mensagem").Length() > 0 {
		return "", fmt.Errorf("Wasn't able to find subject")
	}

	// To parse subject description, get <b> element with content "Objetivos" and subject description will be on next <tr>
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

func scrapeSubjectStats(doc *goquery.Document) (class, assign int, total string, err error) {
	defer func() {
		if r := recover(); r != nil {
			class, assign, total = -1, -1, ""
			err = fmt.Errorf("Couldnt get subject stats: %v", r)
		}
	}()

	/* This is a really bad way of getting these (getting first 3 matches), but I dont think
	this terrible website will ever change its terrible design, so it will probably
	continue to work, if the stats break, fix this please.
	*/

	search := doc.Find("tr[valign=\"TOP\"][align=\"LEFT\"] > td > font > span[class=\"txt_arial_8pt_gray\"]")
	classCredits := strings.TrimSpace(search.Eq(0).Text())
	class, _ = strconv.Atoi(classCredits)

	assignCredits := strings.TrimSpace(search.Eq(1).Text())
	assign, _ = strconv.Atoi(assignCredits)

	totalHours := strings.Trim(search.Eq(2).Text(), " \n\t")
	space := regexp.MustCompile(`\s+`)
	total = space.ReplaceAllString(totalHours, " ")

	return class, assign, total, nil
}

func scrapeSubjectRequirements(doc *goquery.Document, subCode string, courseCode string) (reqs []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Couldn't get subject requirements: %v", r)
		}
	}()

	if doc.Find("#web_mensagem").Length() > 0 { // If a subject has no requirements from any course, a message with text "Esta disciplina nao tem requisitos" will appear
		return []string{}, nil
	}

	var trReq *goquery.Selection = nil
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		regexCode := regexp.MustCompile(`Curso:\s+(\d+)`) // regex to get course code
		codeMatches := regexCode.FindStringSubmatch(s.Text())
		if codeMatches != nil {
			if codeMatches[1] == courseCode {
				trReq = s.Closest("tr") // Found section where subject requirements are
			}
		}
	})

	if trReq == nil { // if didn't find section with course code, the subject has no requirements from this course
		return []string{}, nil
	}

	seen := map[string]bool{} // map used to avoid repeated subjects in slice of requirements
	var answer []string

	for {
		trReq = trReq.Next()
		re := regexp.MustCompile(`([A-Z]{3}\d{4})`) // regex to get subject code
		text := trReq.Text()

		matches := re.FindAllStringSubmatch(text, -1)

		if matches == nil { // no more requirements from this course
			break
		}

		for _, code := range matches {
			if len(code) > 0 && code[0] != subCode && seen[code[0]] == false {
				answer = append(answer, code[0])
				seen[code[0]] = true
			}
		}
	}

	return answer, nil
}

func scrapeSubject(subjectURL string, courseCode string, isOptional bool, results chan<- Subject, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, body := httpGetWithUTF8(jupiterURL + subjectURL)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	checkPanic(err)

	// subject has to have a name / code otherwise panic
	subCode, subName, err := scrapeSubjectNames(doc)
	checkPanic(err)

	// Get subject description text
	subDesc, err := scrapeSubjectDescription(doc)

	if err != nil {
		log.Printf("Error getting %v description\n", subCode)
	}

	// Get subject stats, such as class credits, work credits etc
	subClass, subAssign, subTotal, err := scrapeSubjectStats(doc)

	if err != nil {
		log.Printf("Error getting %v stats\n", subCode)
	}

	// Get requirements of subject
	requirementsURL := "https://uspdigital.usp.br/jupiterweb/listarCursosRequisitos?coddis=%v"
	reqURL := fmt.Sprintf(requirementsURL, subCode)

	reqResp, body := httpGetWithUTF8(reqURL)
	defer reqResp.Body.Close()

	reqDoc, err := goquery.NewDocumentFromReader(body)
	checkPanic(err)

	subRequirements, err := scrapeSubjectRequirements(reqDoc, subCode, courseCode)

	if err != nil {
		log.Printf("Error getting %v requirements\n", subCode)
	}

	subject := Subject{
		Code:          subCode,
		Name:          subName,
		Description:   subDesc,
		ClassCredits:  subClass,
		AssignCredits: subAssign,
		TotalHours:    subTotal,
		Requirements:  subRequirements,
		Optional:      isOptional,
	}

	results <- subject
}

// GetSubjects scrapes all subjects from a course page
func GetSubjects(courseURL string, courseCode string) ([]Subject, error) {
	resp, body := httpGetWithUTF8(courseURL)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	checkPanic(err)

	optional := false

	sections := doc.Find("tr[bgcolor='#658CCF']") // Finds section "Disciplinas Obrigatórias"

	if sections.Length() == 0 {
		return []Subject{}, fmt.Errorf("Invalid courseURL")
	}

	c := make(chan Subject, 200)
	wg := &sync.WaitGroup{}

	sections.Each(func(i int, s *goquery.Selection) {
		subjects := s.NextUntil("tr[bgcolor='#658CCF']").Find("td > .link_gray")

		subjects.Each(func(i int, s *goquery.Selection) {
			subjectURL, exists := s.Attr("href")

			if !exists {
				log.Printf("%s has no subject page", strings.TrimSpace(s.Text()))
			}

			wg.Add(1)
			go scrapeSubject(subjectURL, courseCode, optional, c, wg)
		})

		optional = true // after the first section, all subjects are optional
	})

	var results []Subject
	wg.Wait()
	close(c)

	for subj := range c {
		results = append(results, subj)
	}

	return results, nil
}

// ScrapeICMC scrapes the whole institute (every course)
func ScrapeICMC() (courses []Course, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error scraping ICMC courses: %v", r)
		}
	}()

	allCoursesURL := jupiterURL + "jupCursoLista?codcg=55&tipo=N"
	resp, body := httpGetWithUTF8(allCoursesURL)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	checkPanic(err)

	doc.Find("td[valign=\"top\"] a.link_gray").Each(func(i int, s *goquery.Selection) {
		courseURL, exists := s.Attr("href")

		if !exists {
			panic("Couldnt fetch course")
		}

		courseName := strings.TrimSpace(s.Text())

		regexCode := regexp.MustCompile(`codcur=(\d+)`)
		courseCodeMatches := regexCode.FindStringSubmatch(courseURL)
		if courseCodeMatches == nil {
			panic("Couldn't find course code of %v" + courseName)
		}

		courseCode := courseCodeMatches[1]
		subjects, err := GetSubjects(jupiterURL+courseURL, courseCode)
		checkPanic(err)

		courseObj := Course{
			Name:     courseName,
			Code:     courseCode,
			Subjects: subjects,
		}

		courses = append(courses, courseObj)
	})

	return courses, nil
}
