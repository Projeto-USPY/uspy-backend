/* package main runs the backend router*/
package main

import (
	"log"
	"os"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/server"
)

func main() {
	DB := db.SetupDB(".env")
	r, err := server.SetupRouter(DB)
	if err != nil {
		log.Fatal(err)
	}

	_ = r.Run(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
}
