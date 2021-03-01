/* package main runs the backend router*/
package main

import (
	"log"
	"os"

	"github.com/tpreischadt/ProjetoJupiter/builder"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/server"
)

func main() {
	DB := db.SetupDB(".env")

	if s, ok := os.LookupEnv("BUILD"); ok && s == "TRUE" {
		for name, b := range builder.Builders {
			log.Println("executing builder", name)
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
