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

		exists := make(map[string]bool)
		hashes := make([]string, 0, 200)
		for _, offer := range offerings {
			if _, ok := exists[offer.Hash()]; ok {
				continue
			}

			log.Printf("inserting %v offering\n", prof.Name)
			go DB.Insert(offer, "offerings")

			exists[offer.Hash()] = true
			hashes = append(hashes, offer.Hash())
			cntOffs++
		}

		prof.Offerings = hashes
		prof.Stats = map[string]int{
			"sumDidactics": 0,
			"sumRigorous":  0,
			"sumWorthIt":   0,
		}

		log.Printf("inserting professor %v\n", prof.Name)
		go DB.Insert(prof, "professors")
		cntProfs++
	}

	log.Printf("cntOffs: %v, cntProfs: %v\n", cntOffs, cntProfs)
	return cntOffs + cntProfs, nil
}
