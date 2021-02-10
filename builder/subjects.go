// package builder contains useful functions for building the Firestore database
// Use with caution, because it can overwrite most data present in the database, including reviews and statistics
package builder

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/scraper/icmc/subject"
)

type SubjectBuilder struct{}

// SubjectBuilder.Build scrapes all icmc courses and subjects and builds them onto Firestore
func (SubjectBuilder) Build(DB db.Env) error {
	courses, err := subject.ScrapeICMCCourses()
	if err != nil {
		return err
	}

	objs := make([]db.Object, 0)
	for _, course := range courses {
		courseSubNames := make(map[string]string)
		for _, sub := range course.Subjects {
			sub.Stats = map[string]int{
				"worth_it": 0,
				"total":    0,
			}
			objs = append(objs, db.Object{Collection: "subjects", Doc: sub.Hash(), Data: sub})
			courseSubNames[sub.Code] = sub.Name
		}
		course.SubjectCodes = courseSubNames
		objs = append(objs, db.Object{Collection: "courses", Doc: course.Hash(), Data: course})
	}

	for _, o := range objs {
		var err error
		go func() {
			err = DB.Insert(o.Data, o.Collection)
		}()

		if err != nil {
			return err
		}
	}

	return nil
}
