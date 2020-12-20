package subject

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

// GetByCode returns entity.Subject with code 'code'
func GetByCode(DB db.Env, code string) (entity.Subject, error) {
	subject := entity.Subject{Code: code}
	snap, err := DB.Restore("subjects", subject.Hash())
	if err != nil {
		return entity.Subject{}, err
	}
	err = snap.DataTo(&subject)
	if err != nil {
		return entity.Subject{}, err
	}
	return subject, nil
}
