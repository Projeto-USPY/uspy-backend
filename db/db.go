package db

import (
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"os"
)

type Manager interface {
	Insert(db Env, collection string) error
}

type Object struct {
	Collection string
	Doc        string
	Data       Manager
}

type Env struct {
	Client *firestore.Client
	Ctx    context.Context
}

func (db Env) Restore(collection, HashID string) (*firestore.DocumentSnapshot, error) {
	snap, err := db.Client.Collection(collection).Doc(HashID).Get(db.Ctx)
	if err != nil {
		return nil, err
	}

	return snap, nil
}

func (db Env) RestoreCollection(collection string) ([]*firestore.DocumentSnapshot, error) {
	snap, err := db.Client.Collection(collection).Documents(db.Ctx).GetAll()
	if err != nil {
		return nil, err
	}

	return snap, nil
}

func (db Env) Insert(obj Manager, collection string) error {
	err := obj.Insert(db, collection)
	if err != nil {
		return err
	}
	return nil
}

func (db Env) BatchWrite(objs []Object) error {
	batch := db.Client.Batch()

	for _, o := range objs {
		if o.Doc == "" {
			batch.Set(db.Client.Collection(o.Collection).NewDoc(), o.Data)
		} else {
			batch.Set(db.Client.Collection(o.Collection).Doc(o.Doc), o.Data)
		}
	}
	_, err := batch.Commit(db.Ctx)
	return err
}

func InitFireStore(mode string) Env {
	var DB = Env{
		Ctx: context.Background(),
	}

	if mode == "prod" {
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
	} else { // build or dev
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
