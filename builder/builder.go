// package builder contains useful functions for building the Firestore database
// Use with caution, because it can overwrite most data present in the database, including reviews and statistics
package builder

import "github.com/tpreischadt/ProjetoJupiter/db"

// Builder interface is implemented by all types that somehow populate the database
type Builder interface {
	Build(db.Env) error
}

var Builders = map[string]Builder{
	"SubjectBuilder": SubjectBuilder{},
}