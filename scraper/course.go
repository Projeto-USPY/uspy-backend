package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"io"
)

type CourseScraper struct {
	URLMask        string
	InstituteCode  string
	Code           string
	Specialization string
}

func NewCourseScraper(institute, course, spec string) CourseScraper {
	return CourseScraper{
		URLMask:        DefaultCourseURLMask,
		Code:           course,
		Specialization: spec,
		InstituteCode:  institute,
	}
}
func (sc CourseScraper) Start() (db.Manager, error) {
	URL := fmt.Sprintf(sc.URLMask, sc.InstituteCode, sc.Code, sc.Specialization)
	return Start(sc, URL)
}

func (sc CourseScraper) Scrape(reader io.Reader) (db.Manager, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	course := entity.Course{
		Name:         doc.Find("td > font:nth-child(2) > span").Last().Text(),
		Code:         sc.Code,
		Subjects:     make([]entity.Subject, 0, 1000),
		SubjectCodes: make(map[string]string, 0),
	}

	// Get Subjects

	return course, nil
}
