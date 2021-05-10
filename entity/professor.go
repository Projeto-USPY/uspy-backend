/* package entity contains structs that will be used for backend input validation and DB operations */
package entity

import (
	"github.com/Projeto-USPY/uspy-backend/db"
)

// entity.Professor represents a professor
// Example: {1234567, "Fulano da Silva", Stats, Offerings}
type Professor struct {
	CodPes string
	Name   string

	Offerings []Offering
}

func (prof Professor) Insert(DB db.Env, collection string) error {
	return nil
}

func (prof Professor) Update(DB db.Env, collection string) error {
	return nil
}
