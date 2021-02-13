/* package main runs the backend router*/
package main

import (
	"github.com/tpreischadt/ProjetoJupiter/builder"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/server"
	"log"
	"os"
)

func main() {
	DB := db.SetupDB(".env")

	if s, ok := os.LookupEnv("BUILD"); ok && s == "true" {
		for _, b := range builder.Builders {
			err := b.Build(DB)
			if err != nil {
				log.Fatal("failed to build: ", err)
			}
		}
		return
	}

	r, err := server.SetupRouter(DB)
	if err != nil {
		log.Fatal(err)
	}

	_ = r.Run(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
}
