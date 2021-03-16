/* package main runs the backend router*/
package main

import (
	"log"
	"os"

	"github.com/Projeto-USPY/uspy-backend/builder"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/server"
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
