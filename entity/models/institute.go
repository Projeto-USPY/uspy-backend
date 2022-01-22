package models

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// Institute represents an institute or collection of courses
type Institute struct {
	Name string `firestore:"name"`
	Code string `firestore:"code"`

	// Attributes only used to nest collected data by uspy-scraper
	Courses    []Course    `firestore:"-"`
	Professors []Professor `firestore:"-"`
}

// NewInstituteFromController is a constructor. It takes a institute controller and returns a model.
func NewInstituteFromController(inst *controllers.Institute) *Institute {
	return &Institute{Code: inst.Code}
}

// Hash returns SHA256(code)
func (i Institute) Hash() string {
	return utils.SHA256(i.Code)
}

// Insert sets an institute to a given collection. This is usually /institutes
func (i Institute) Insert(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(i.Hash()).Set(DB.Ctx, i)
	return err
}

// Update sets an institute to a given collection. This is usually /institutes
func (i Institute) Update(DB db.Env, collection string) error {
	_, err := DB.Client.Collection(collection).Doc(i.Hash()).Set(DB.Ctx, i)
	return err
}
