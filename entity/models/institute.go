package models

import (
	"github.com/Projeto-USPY/uspy-backend/db"
)

type Institute struct {
	Name    string
	Code    string
	Courses []Course

	Professors []Professor
}

func (i Institute) Insert(DB db.Env, collection string) error { return nil }

func (i Institute) Update(DB db.Env, collection string) error { return nil }
