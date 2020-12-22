package main

import (
	"github.com/tpreischadt/ProjetoJupiter/server"
	"os"
)

func main() {
	r := server.SetupRouter(server.SetupDB(".env"))
	_ = r.Run(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
}
