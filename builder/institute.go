package builder

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/scraper"
	"log"
	"sync"
)

type InstituteBuilder struct{}

// InstituteBuilder Build scrapes all icmc courses and subjects and builds them onto Firestore
func (InstituteBuilder) Build(DB db.Env) error {
	log.Println("scraping institute")
	instituteObj := scraper.NewInstituteScraper("55")
	if instituteObj, err := instituteObj.Start(); err != nil {
		return err
	} else {
		log.Println("done")
		objs := make([]db.Object, 0)

		var institute = instituteObj.(entity.Institute)

		for _, course := range institute.Courses {
			for _, sub := range course.Subjects {
				objs = append(objs, db.Object{Collection: "subjects", Doc: sub.Hash(), Data: sub})
			}
			objs = append(objs, db.Object{Collection: "courses", Doc: course.Hash(), Data: course})
		}

		var wg sync.WaitGroup
		errors := make(chan error, 10000)
		for _, o := range objs {
			wg.Add(1)
			go func(obj db.Object, group *sync.WaitGroup) {
				defer group.Done()
				errors <- DB.Insert(obj.Data, obj.Collection)
			}(o, &wg)

			log.Printf("inserting %v into %v\n", o.Doc, o.Collection)
		}

		wg.Wait()
		close(errors)

		for e := range errors {
			if e != nil {
				return e
			}
		}

		log.Printf("inserted %d total objects\n", len(objs))
	}

	return nil
}
