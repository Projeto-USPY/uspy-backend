package subject

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func scrapeSubjectNames(doc *goquery.Document) (code, name string, e error) {
	defer func() {
		if r := recover(); r != nil {
			code, name, e = "", "", fmt.Errorf("Error getting subjects name or code: %v", r)
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
	// Returns error telling that subjects is invalid or not yet activated
	if doc.Find("#web_mensagem").Length() > 0 {
		return "", fmt.Errorf("Wasn't able to find subjects")
	}

	// To parse subjects description, get <b> element with content "Objetivos" and subjects description will be on next <tr>
	var objetivosNode *goquery.Selection = nil
	doc.Find("b").Each(func(i int, s *goquery.Selection) {
		html, err := s.Html() // get inner html
		utils.CheckPanic(err)

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
			err = fmt.Errorf("Couldnt get subjects stats: %v", r)
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
			err = fmt.Errorf("Couldn't get subjects requirements: %v", r)
		}
	}()

	if doc.Find("#web_mensagem").Length() > 0 { // If a subjects has no requirements from any course, a message with text "Esta disciplina nao tem requisitos" will appear
		return []string{}, nil
	}

	var trReq *goquery.Selection = nil
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		regexCode := regexp.MustCompile(`Curso:\s+(\d+)`) // regex to get course code
		codeMatches := regexCode.FindStringSubmatch(s.Text())
		if codeMatches != nil {
			if codeMatches[1] == courseCode {
				trReq = s.Closest("tr") // Found section where subjects requirements are
			}
		}
	})

	if trReq == nil { // if didn't find section with course code, the subjects has no requirements from this course
		return []string{}, nil
	}

	seen := map[string]bool{} // map used to avoid repeated subjects in slice of requirements
	var answer []string

	for {
		trReq = trReq.Next()
		re := regexp.MustCompile(`([A-Z]{3}\d{4})`) // regex to get subjects code
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

func scrapeSubject(subjectURL string, courseCode string, isOptional bool, results chan<- entity.Subject, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, body := utils.HTTPGetWithUTF8(utils.JupiterURL + subjectURL)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	utils.CheckPanic(err)

	// subjects has to have a name / code otherwise panic
	subCode, subName, err := scrapeSubjectNames(doc)
	utils.CheckPanic(err)

	// Get subjects description text
	subDesc, err := scrapeSubjectDescription(doc)

	if err != nil {
		log.Printf("Error getting %v description\n", subCode)
	}

	// Get subjects stats, such as class credits, work credits etc
	subClass, subAssign, subTotal, err := scrapeSubjectStats(doc)

	if err != nil {
		log.Printf("Error getting %v stats\n", subCode)
	}

	// Get requirements of subjects
	requirementsURL := "https://uspdigital.usp.br/jupiterweb/listarCursosRequisitos?coddis=%v"
	reqURL := fmt.Sprintf(requirementsURL, subCode)

	reqResp, body := utils.HTTPGetWithUTF8(reqURL)
	defer reqResp.Body.Close()

	reqDoc, err := goquery.NewDocumentFromReader(body)
	utils.CheckPanic(err)

	subRequirements, err := scrapeSubjectRequirements(reqDoc, subCode, courseCode)

	if err != nil {
		log.Printf("Error getting %v requirements\n", subCode)
	}

	subject := entity.Subject{
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
func GetSubjects(courseURL string, courseCode string) ([]entity.Subject, error) {
	resp, body := utils.HTTPGetWithUTF8(courseURL)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	utils.CheckPanic(err)

	optional := false

	sections := doc.Find("tr[bgcolor='#658CCF']") // Finds section "Disciplinas Obrigatórias"

	if sections.Length() == 0 {
		return []entity.Subject{}, fmt.Errorf("Invalid courseURL")
	}

	c := make(chan entity.Subject, 200)
	wg := &sync.WaitGroup{}

	sections.Each(func(i int, s *goquery.Selection) {
		subjects := s.NextUntil("tr[bgcolor='#658CCF']").Find("td > .link_gray")

		subjects.Each(func(i int, s *goquery.Selection) {
			subjectURL, exists := s.Attr("href")

			if !exists {
				log.Printf("%s has no subjects page", strings.TrimSpace(s.Text()))
			}

			wg.Add(1)
			go scrapeSubject(subjectURL, courseCode, optional, c, wg)
		})

		optional = true // after the first section, all subjects are optional
	})

	var results []entity.Subject
	wg.Wait()
	close(c)

	for subj := range c {
		results = append(results, subj)
	}

	return results, nil
}
