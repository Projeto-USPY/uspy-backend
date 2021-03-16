package scraper

import (
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
	"github.com/PuerkitoBio/goquery"
	"io"
	"regexp"
	"strings"
)

type InstituteScraper struct {
	URLMask string
	Code    string
}

func NewInstituteScraper(institute string) InstituteScraper {
	return InstituteScraper{
		URLMask: DefaultInstituteURLMask,
		Code:    institute,
	}
}

func (sc InstituteScraper) Start() (db.Manager, error) {
	URL := fmt.Sprintf(sc.URLMask, sc.Code)
	return Start(sc, URL)
}

func (sc InstituteScraper) Scrape(reader io.Reader) (obj db.Manager, err error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	institute := entity.Institute{
		Name:    strings.TrimSpace(doc.Find("span > b").Text()),
		Code:    sc.Code,
		Courses: make([]entity.Course, 0, 50),
	}

	coursesHref := doc.Find("td[valign=\"top\"] a.link_gray")
	for i := 0; i < coursesHref.Length(); i++ {
		// follow every course href
		node := coursesHref.Eq(i)
		if courseCode, courseSpec, err := getCourseIdentifiers(node); err != nil {
			return nil, err
		} else {
			courseScraper := NewCourseScraper(sc.Code, courseCode, courseSpec)
			if course, err := courseScraper.Start(); err != nil {
				return nil, err
			} else {
				institute.Courses = append(institute.Courses, course.(entity.Course))
			}
		}
	}
	return institute, nil
}

func getCourseIdentifiers(node *goquery.Selection) (string, string, error) {
	if courseURL, exists := node.Attr("href"); exists {
		// get course code and specialization code
		regexCode := regexp.MustCompile(`codcur=(\d+)&codhab=(\d+)`)
		courseCodeMatches := regexCode.FindStringSubmatch(courseURL)

		if len(courseCodeMatches) < 3 {
			return "", "", ErrorCourseNotExist
		}

		courseCode, courseSpec := courseCodeMatches[1], courseCodeMatches[2]
		return courseCode, courseSpec, nil
	} else {
		return "", "", ErrorCourseNotExist
	}
}
