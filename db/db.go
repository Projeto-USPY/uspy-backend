/* Package db contains useful functions related to the Firestore Database */
package db

import (
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"os"
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

// Object is used for batched writes that can contain different types that implement Inserter
// Set Doc to empty string if you'd like to use a random Hash
type Object struct {
	Collection string
	Doc        string
	Data       Writer
}

// Env is passed to /server/models functions that require DB operations
type Env struct {
	Client *firestore.Client
	Ctx    context.Context
}

// Env.Restore restores a document with a specific HashID and collection origin from Firestore
// collection cannot end in "/"
func (db Env) Restore(collection, HashID string) (*firestore.DocumentSnapshot, error) {
	snap, err := db.Client.Collection(collection).Doc(HashID).Get(db.Ctx)
	if err != nil {
		return nil, err
	}

	return snap, nil
}

// Env.RestoreCollection is similar to Env.Restore, but restores all documents from a collection
// collection cannot end in "/"
func (db Env) RestoreCollection(collection string) ([]*firestore.DocumentSnapshot, error) {
	snap, err := db.Client.Collection(collection).Documents(db.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	return snap, nil
}

// Env.Insert inserts an entity that implements Inserter into a DB collection
func (db Env) Insert(obj Inserter, collection string) error {
	return obj.Insert(db, collection)
}

// Env.Update performs one or more updates in a single document
func (db Env) Update(doc, col string, updates []firestore.Update) error {
	_, err := db.Client.Collection(col).Doc(doc).Update(db.Ctx, updates)
	return err
}

// Env.BatchWrite will perform inserts atomically
func (db Env) BatchWrite(objs []Object) error {
	batch := db.Client.Batch()

	for _, o := range objs {
		if o.Doc == "" { // create document with random hash
			batch.Set(db.Client.Collection(o.Collection).NewDoc(), o.Data)
		} else {
			batch.Set(db.Client.Collection(o.Collection).Doc(o.Doc), o.Data)
		}
	}
	_, err := batch.Commit(db.Ctx)
	return err
}

// InitFirestore receives a mode dev/prod and initiates the DB Environment
func InitFireStore(LOCAL string) Env {
	var DB = Env{
		Ctx: context.Background(),
	}

	if LOCAL == "FALSE" {
		if id, ok := os.LookupEnv("PROJECT_ID"); ok {
			conf := &firebase.Config{ProjectID: id}
			app, err := firebase.NewApp(DB.Ctx, conf)
			if err != nil {
				log.Fatalln(err)
			}

			DB.Client, err = app.Firestore(DB.Ctx)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatal("missing env variable PROJECT_ID")
		}
	} else {
		if key, ok := os.LookupEnv("FIRESTORE_KEY"); ok {
			sa := option.WithCredentialsFile(key)
			app, err := firebase.NewApp(DB.Ctx, nil, sa)
			if err != nil {
				log.Fatalln(err)
			}

			DB.Client, err = app.Firestore(DB.Ctx)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatal("FIRESTORE_KEY path not specified in .env file")
		}
	}

	return DB
}

func SetupDB(envPath string) Env {
	_ = godotenv.Load(envPath)
	return InitFireStore(os.Getenv("LOCAL"))
}
