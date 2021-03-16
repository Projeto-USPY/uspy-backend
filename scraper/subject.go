package scraper

import (
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
	"github.com/PuerkitoBio/goquery"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type SubjectScraper struct {
	URLMask        string
	Code           string
	CourseCode     string
	Specialization string
}

func NewSubjectScraper(subject, course, spec string) SubjectScraper {
	return SubjectScraper{
		URLMask:        DefaultSubjectURLMask,
		Code:           subject,
		CourseCode:     course,
		Specialization: spec,
	}
}

func (sc SubjectScraper) Start() (db.Manager, error) {
	URL := fmt.Sprintf(sc.URLMask, sc.Code, sc.CourseCode, sc.Specialization)
	return Start(sc, URL)
}

func (sc SubjectScraper) Scrape(reader io.Reader) (db.Manager, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	fullName := doc.Find("span.txt_arial_10pt_black > b").Text()
	fields := strings.Split(fullName, "-")
	name := strings.TrimSpace(fields[len(fields)-1])

	subject := entity.Subject{
		Code:           sc.Code,
		CourseCode:     sc.CourseCode,
		Specialization: sc.Specialization,
		Name:           name,
		Stats: map[string]int{
			"total":    0,
			"worth_it": 0,
		},
	}

	if description, err := getDescription(doc); err == nil {
		subject.Description = description
	} else {
		return nil, err
	}

	search := doc.Find("tr[valign=\"TOP\"][align=\"LEFT\"] > td > font > span[class=\"txt_arial_8pt_gray\"]")
	if class, err := getClassCredits(search); err == nil {
		subject.ClassCredits = class
	} else {
		return nil, err
	}

	if assign, err := getAssignCredits(search); err == nil {
		subject.AssignCredits = assign
	} else {
		return nil, err
	}

	if total, err := getTotalHours(search); err == nil {
		subject.TotalHours = total
	} else {
		return nil, err
	}

	return subject, nil
}

func getDescription(doc *goquery.Document) (string, error) {
	var objetivosNode *goquery.Selection = nil
	bold := doc.Find("b")

	for i := 0; i < bold.Length(); i++ {
		s := bold.Eq(i)
		text := s.Text() // get inner html

		if strings.TrimSpace(text) == "Objetivos" { // found
			objetivosNode = s
		}
	}

	if objetivosNode == nil {
		return "", nil
	}

	objetivosTr := objetivosNode.Closest("tr") // get tr parent
	descriptionTr := objetivosTr.Next()        // tr with description is next <tr>

	desc := strings.TrimSpace(descriptionTr.Text())
	return desc, nil
}

func getClassCredits(search *goquery.Selection) (int, error) {
	classCredits := strings.TrimSpace(search.Eq(0).Text())
	class, err := strconv.Atoi(classCredits)

	if err != nil {
		return -1, err
	}

	return class, nil
}

func getAssignCredits(search *goquery.Selection) (int, error) {
	assignCredits := strings.TrimSpace(search.Eq(1).Text())
	assign, err := strconv.Atoi(assignCredits)

	if err != nil {
		return -1, err
	}

	return assign, nil
}

func getTotalHours(search *goquery.Selection) (string, error) {
	totalHours := strings.Trim(search.Eq(2).Text(), " \n\t")
	space, err := regexp.Compile(`\s+`)
	if err != nil {
		return "", err
	}

	total := space.ReplaceAllString(totalHours, " ")
	return total, nil
}
