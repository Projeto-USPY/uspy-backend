package entity

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
)

type Institute struct {
	Name    string
	Code    string
	Courses []Course
}

func (i Institute) Insert(DB db.Env, collection string) error {
	for _, c := range i.Courses {
		if err := c.Insert(DB, "courses"); err != nil {
			return err
		}
	}

	return nil
}
