/*Package db contains useful functions related to the Firestore Database */
package db

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/Projeto-USPY/uspy-backend/config"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// Inserter will be implemented by almost all entities
type Inserter interface {
	Insert(db Database, collection string) error
}

// Updater will be implemented by almost all entities
type Updater interface {
	Update(db Database, collection string) error
}

// Writer implements Inserter and Updater (InserterUpdater is a bad name)
type Writer interface {
	Inserter
	Updater
}

// BatchObject is used for batched writes that can contain different types that implement Inserter
// Set Doc to empty string if you'd like to use a random Hash
type BatchObject struct {
	Collection string
	Doc        string

	WriteData  Writer
	UpdateData []firestore.Update

	Preconditions []firestore.Precondition
	SetOptions    []firestore.SetOption
}

// Operation is used as a generic operation to be applied on a document
//
// It is mostly used inside transactions to provide an easy way to store operations to be executed after reads
type Operation struct {
	Ref     *firestore.DocumentRef
	Method  string
	Payload interface{}

	Err error
}

// Database is passed to /server/dao functions that require DB operations
type Database struct {
	Client *firestore.Client
	Ctx    context.Context
}

// Restore restores a document with a specific hash
//
// If the document is not found, returns an error which can be checked with
// status.Code(err) == codes.NotFound
//
// Besides, the Exists method for this Ref will return false
func (db Database) Restore(documentHash string) (*firestore.DocumentSnapshot, error) {
	snap, err := db.Client.Doc(documentHash).Get(db.Ctx)
	if err != nil {
		return nil, err
	}

	return snap, nil
}

// RestoreBatch is similar to Env.Restore, but restores a batch of documents concurrently
//
// If any document is not found, the Exists method for that snap will return false
//
// It is guaranteed that snapshots are returned in the same order as passed hashes
func (db Database) RestoreBatch(documentHashes []string) ([]*firestore.DocumentSnapshot, error) {
	refs := make([]*firestore.DocumentRef, 0, len(documentHashes))
	for _, doc := range documentHashes {
		refs = append(refs, db.Client.Doc(doc))
	}

	return db.Client.GetAll(db.Ctx, refs)
}

// RestoreCollection is similar to Env.Restore, but restores all documents from a collection
//
// Collection cannot end in "/"
func (db Database) RestoreCollection(collection string) ([]*firestore.DocumentSnapshot, error) {
	snap, err := db.Client.Collection(collection).Documents(db.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	return snap, nil
}

// RestoreCollectionRefs is similar to RestoreCollection, but uses DocRefs that allow missing documents inside the query
//
// Collection cannot end in "/"
func (db Database) RestoreCollectionRefs(collection string) ([]*firestore.DocumentRef, error) {
	snap, err := db.Client.Collection(collection).DocumentRefs(db.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	return snap, nil
}

// Insert inserts an entity that implements Inserter into a DB collection
func (db Database) Insert(obj Inserter, collection string) error {
	return obj.Insert(db, collection)
}

// Update updates entity in firestore with data in object variable
func (db Database) Update(obj Updater, collection string) error {
	return obj.Update(db, collection)
}

// BatchWrite will perform operations atomically
//
// For a batch of more than 500 documents, batch write will perform each of these batches sequentially
// TODO: Apply batches concurrently
func (db Database) BatchWrite(objs []BatchObject) error {
	numObjs := len(objs)
	for i := 0; i < numObjs; i += 500 {
		last := i + 500
		if last > numObjs {
			last = numObjs
		}

		// perform batch of at maximum 500 operations
		batch := db.Client.Batch()

		for j := i; j < last; j++ {
			o := objs[j]

			if o.WriteData == nil && o.UpdateData == nil {
				return errors.New("both write data and update data are nil")
			}

			if o.Doc == "" { // create document with random hash
				batch.Set(db.Client.Collection(o.Collection).NewDoc(), o.WriteData, o.SetOptions...)
			} else {
				if o.WriteData != nil { // set operation
					batch.Set(db.Client.Collection(o.Collection).Doc(o.Doc), o.WriteData, o.SetOptions...)
				} else if o.UpdateData != nil { // update operation
					batch.Update(db.Client.Collection(o.Collection).Doc(o.Doc), o.UpdateData, o.Preconditions...)
				}
			}
		}

		if _, err := batch.Commit(db.Ctx); err != nil {
			return err
		}
	}

	return nil
}

// InitFireStore initiates the DB Environment (requires some environment variables to work)
func InitFireStore() Database {
	var DB = Database{
		Ctx: context.Background(),
	}

	if config.Env.IsUsingProjectID() {
		conf := &firebase.Config{ProjectID: config.Env.Identify()}
		app, err := firebase.NewApp(DB.Ctx, conf)
		if err != nil {
			log.Fatal(err)
		}

		DB.Client, err = app.Firestore(DB.Ctx)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		sa := option.WithCredentialsFile(config.Env.Identify())

		app, err := firebase.NewApp(DB.Ctx, nil, sa)
		if err != nil {
			log.Fatal(err)
		}

		DB.Client, err = app.Firestore(DB.Ctx)
		if err != nil {
			log.Fatal("There might be something wrong with your credentials file!")
		}
	}

	return DB
}

// SetupDB wraps the Firestore initialization
func SetupDB() Database {
	return InitFireStore()
}
