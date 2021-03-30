package entity

import (
	"github.com/Projeto-USPY/uspy-backend/db"
)

type Institute struct {
	Name    string   `json:"name"`
	Code    string   `json:"code"`
	Courses []Course `json:"courses"`
}

func (i Institute) Insert(DB db.Env, collection string) error { return nil }
