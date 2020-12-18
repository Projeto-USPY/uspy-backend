package populator

import (
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/scraper/icmc/subject"
	"log"
)

func PopulateICMCSubjects(DB db.Env) (int, error) {
	log.Println("scraping icmc courses")
	courses, err := subject.ScrapeICMCCourses()
	if err != nil {
		return 0, err
	}

	cntCourses, cntSubjects := 0, 0
	for _, course := range courses {
		log.Println("inserting course", course.Name)
		err := course.Insert(DB, "courses")
		if err != nil {
			return 0, nil
		}
		for _, sub := range course.Subjects {
			log.Println("inserting subjects from course", course.Name)
			go sub.Insert(DB, fmt.Sprintf("courses/%s/subjects", course.Hash()))
			cntSubjects++
		}
		cntCourses++
	}
	return cntCourses + cntSubjects, nil
}
