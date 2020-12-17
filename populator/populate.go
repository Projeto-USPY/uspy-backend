package populator

import (
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/scraper"
	"github.com/tpreischadt/ProjetoJupiter/scraper/icmc/professor"
	"log"
)

func PopulateICMCOfferings(DB db.Env) (int, error) {
	log.Print("scraping all icmc departments")
	icmcProfessors := professor.ScrapeDepartments()
	cntOffs, cntProfs := 0, 0
	for _, prof := range icmcProfessors {
		log.Printf("getting %v history\n", prof.Name)
		offerings, err := scraper.GetProfessorHistory(prof.CodPes, 2010)
		if err != nil {
			return -1, err
		}

		for _, offer := range offerings {
			offerDB := db.NewOffering(offer)
			log.Printf("inserting %v offerings\n", prof.Name)
			go DB.Insert(offerDB)
			cntOffs++
		}

		log.Printf("inserting professor %v\n", prof.Name)
		profDB, err := db.NewProfessorWithOfferings(prof, offerings)
		go DB.Insert(profDB)
		cntProfs++
	}

	log.Printf("cntOffs: %v, cntProfs: %v\n", cntOffs, cntProfs)
	return cntOffs + cntProfs, nil
}
