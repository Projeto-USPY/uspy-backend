/* package main runs the backend router*/
package main

import (
	"log"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/server"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config.Setup()
}

func main() {
	DB := db.SetupDB()
	r, err := server.SetupRouter(DB)
	if err != nil {
		log.Fatal(err)
	}

	_ = r.Run(config.Env.Domain + ":" + config.Env.Port)
}
