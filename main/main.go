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
	r, err := server.SetupRouter(DB)

	if err != nil {
		log.Fatal(err)
	}

	if s, _ := os.LookupEnv("MODE"); s == "build" {
		for _, b := range builder.Builders {
			err := b.Build(DB)
			if err != nil {
				log.Fatal("failed to build: ", err)
			}
		}
		return
	}

	_ = r.Run(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
}
