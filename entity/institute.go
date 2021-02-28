package entity

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
)

type Institute struct {
	Name    string   `json:"name"`
	Code    string   `json:"code"`
	Courses []Course `json:"courses"`
}

func (i Institute) Insert(DB db.Env, collection string) error {
	for _, c := range i.Courses {
		if err := c.Insert(DB, "courses"); err != nil {
			return err
		}
	}

	return nil
}
