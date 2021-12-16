/*Package db contains useful functions related to the Firestore Database */
package db

import (
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/Projeto-USPY/uspy-backend/config"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// Inserter will be implemented by almost all entities
type Inserter interface {
	Insert(db Env, collection string) error
}

// Updater will be implemented by almost all entities
type Updater interface {
	Update(db Env, collection string) error
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

// Env is passed to /server/dao functions that require DB operations
type Env struct {
	Client *firestore.Client
	Ctx    context.Context
}

// Restore restores a document with a specific HashID and collection origin from Firestore
// collection cannot end in "/"
func (db Env) Restore(collection, HashID string) (*firestore.DocumentSnapshot, error) {
	snap, err := db.Client.Collection(collection).Doc(HashID).Get(db.Ctx)
	if err != nil {
		return nil, err
	}

	return snap, nil
}

// RestoreCollection is similar to Env.Restore, but restores all documents from a collection
//
// Collection cannot end in "/"
func (db Env) RestoreCollection(collection string) ([]*firestore.DocumentSnapshot, error) {
	snap, err := db.Client.Collection(collection).Documents(db.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	return snap, nil
}

// Insert inserts an entity that implements Inserter into a DB collection
func (db Env) Insert(obj Inserter, collection string) error {
	return obj.Insert(db, collection)
}

// Update updates entity in firestore with data in object variable
func (db Env) Update(obj Updater, collection string) error {
	return obj.Update(db, collection)
}

// BatchWrite will perform operations atomically
func (db Env) BatchWrite(objs []BatchObject) error {
	batch := db.Client.Batch()

	for _, o := range objs {
		if o.WriteData == nil && o.UpdateData == nil {
			return errors.New("both write data and update data are nil")
		}

		if o.Doc == "" { // create document with random hash
			batch.Set(db.Client.Collection(o.Collection).NewDoc(), o.WriteData)
		} else {
			if o.WriteData != nil { // set operation
				batch.Set(db.Client.Collection(o.Collection).Doc(o.Doc), o.WriteData)
			} else if o.UpdateData != nil { // update operation
				batch.Update(db.Client.Collection(o.Collection).Doc(o.Doc), o.UpdateData)
			}
		}
	}
	_, err := batch.Commit(db.Ctx)
	return err
}

// InitFireStore initiates the DB Environment (requires some environment variables to work)
func InitFireStore() Env {
	var DB = Env{
		Ctx: context.Background(),
	}

	if config.Env.IsUsingProjectID() {
		conf := &firebase.Config{ProjectID: config.Env.Identify()}
		app, err := firebase.NewApp(DB.Ctx, conf)
		if err != nil {
			log.Fatalln(err)
		}

		DB.Client, err = app.Firestore(DB.Ctx)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		sa := option.WithCredentialsFile(config.Env.Identify())

		app, err := firebase.NewApp(DB.Ctx, nil, sa)
		if err != nil {
			log.Fatalln(err)
		}

		DB.Client, err = app.Firestore(DB.Ctx)
		if err != nil {
			log.Println(err)
			log.Fatalln("There might be something wrong with your credentials file!")
		}
	}

	return DB
}

// SetupDB wraps the Firestore initialization
func SetupDB() Env {
	return InitFireStore()
}
