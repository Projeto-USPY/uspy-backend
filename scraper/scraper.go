package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/scraper/subject"
	"github.com/tpreischadt/ProjetoJupiter/utils"
	"regexp"
	"strings"
)

// ScrapeICMC scrapes the whole institute (every course)
func ScrapeICMC() (courses []entity.Course, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error scraping ICMC courses: %v", r)
		}
	}()

	allCoursesURL := utils.JupiterURL + "jupCursoLista?codcg=55&tipo=N"
	resp, body := utils.HTTPGetWithUTF8(allCoursesURL)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	utils.CheckPanic(err)

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
		subjs, err := subject.GetSubjects(utils.JupiterURL+courseURL, courseCode)
		utils.CheckPanic(err)

		courseObj := entity.Course{
			Name:     courseName,
			Code:     courseCode,
			Subjects: subjs,
		}

		courses = append(courses, courseObj)
	})

	return courses, nil
}
