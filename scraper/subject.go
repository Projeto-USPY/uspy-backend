package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
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

	subject := entity.Subject{
		Code:           sc.Code,
		CourseCode:     sc.CourseCode,
		Specialization: sc.Specialization,
		Name:           doc.Find("span.txt_arial_10pt_black > b").Text(),
		Stats: map[string]int{
			"total":    0,
			"worth_it": 0,
		},
	}

	search := doc.Find("tr[valign=\"TOP\"][align=\"LEFT\"] > td > font > span[class=\"txt_arial_8pt_gray\"]")
	classCredits := strings.TrimSpace(search.Eq(0).Text())
	if class, err := strconv.Atoi(classCredits); err != nil {
		return subject, err
	} else {
		subject.ClassCredits = class
	}

	assignCredits := strings.TrimSpace(search.Eq(1).Text())
	if assign, err := strconv.Atoi(assignCredits); err != nil {
		return subject, err
	} else {
		subject.AssignCredits = assign
	}

	totalHours := strings.Trim(search.Eq(2).Text(), " \n\t")
	if space, err := regexp.Compile(`\s+`); err != nil {
		return subject, err
	} else {
		total := space.ReplaceAllString(totalHours, " ")
		subject.TotalHours = total
	}

	return subject, nil
}
