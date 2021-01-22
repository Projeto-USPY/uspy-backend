package main

import (
	"github.com/tpreischadt/ProjetoJupiter/server"
	"log"
	"os"
)

func main() {
	r, err := server.SetupRouter(server.SetupDB(".env"))
	if err != nil {
		log.Fatal(err)
	}
	_ = r.Run(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
}
