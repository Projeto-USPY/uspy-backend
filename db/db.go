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

type Manager interface {
	Insert(client *firestore.Client, ctx context.Context) error
	Hash() string
}

func Restore(client *firestore.Client, ctx context.Context, collection, HashID string) (*firestore.DocumentSnapshot, error) {
	snap, err := client.Collection(collection).Doc(HashID).Get(ctx)
	if err != nil {
		return nil, err
	}

	return snap, nil
}

func InitFireStore(envRelPath string) (*firestore.Client, context.Context) {
	_ = godotenv.Load(envRelPath)

	ctx := context.Background()
	var client *firestore.Client

	if mode, ok := os.LookupEnv("MODE"); ok && mode != "prod" { // dev or build
		if key, ok := os.LookupEnv("FIRESTORE_KEY"); ok {
			sa := option.WithCredentialsFile(key)
			app, err := firebase.NewApp(ctx, nil, sa)
			if err != nil {
				log.Fatalln(err)
			}

			client, err = app.Firestore(ctx)
			if err != nil {
				log.Fatalln(err)
			}

			// TODO: Refactor this to optimize DB Reads
			if mode == "build" { // populate and exit
				func() {
					cnt, err := PopulateOfferings(client, ctx)
					if err != nil {
						_ = client.Close()
						log.Fatalln("failed to build: ", err)
					} else {
						log.Println("total: ", cnt)
					}
				}()

				func() {
					cnt, err := PopulateProfessors(client, ctx)
					if err != nil {
						_ = client.Close()
						log.Fatalln("failed to build: ", err)
					} else {
						log.Println("total: ", cnt)
					}
				}()

				_ = client.Close()
				os.Exit(0)
			}

		} else {
			log.Fatal("FIRESTORE_KEY path not specified in .env file")
		}
	} else if ok && mode == "prod" { // production
		if id, ok := os.LookupEnv("PROJECT_ID"); ok {
			conf := &firebase.Config{ProjectID: id}
			app, err := firebase.NewApp(ctx, conf)
			if err != nil {
				log.Fatalln(err)
			}

			client, err = app.Firestore(ctx)
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		log.Fatalln(".env runtime MODE invalid or not specified")
	}

	return client, ctx
}
