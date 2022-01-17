/*package main runs the backend router*/
package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/server"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
		PadLevelText:  true,
	})
	log.SetReportCaller(true)

	config.Setup()

	if config.Env.IsProd() {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	DB := db.SetupDB()
	r, err := server.SetupRouter(DB)
	if err != nil {
		log.Fatal(err)
	}

	_ = r.Run(config.Env.Domain + ":" + config.Env.Port)
}
