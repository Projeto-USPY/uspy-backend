package populator

import (
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
		courseSubHashes := make([]string, 0, 300)
		for _, sub := range course.Subjects {
			sub.Stats = map[string]int{
				"qtWorthIt": 0,
				"qtReviews": 0,
			}
			log.Println("inserting subjects from course", course.Name)
			go DB.Insert(sub, "subjects")
			courseSubHashes = append(courseSubHashes, sub.Hash())
			cntSubjects++
		}
		course.SubjectHashes = courseSubHashes
		err := DB.Insert(course, "courses")
		log.Println("inserting course", course.Name)
		if err != nil {
			return 0, nil
		}
		cntCourses++
	}
	return cntCourses + cntSubjects, nil
}
