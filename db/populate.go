package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/tpreischadt/ProjetoJupiter/scraper/professor"
	"log"
)

func PopulateProfessors(client *firestore.Client, ctx context.Context) (int, error) {
	log.Print("scraping all icmc departments")
	icmcProfessors := professor.ScrapeDepartments()
	cnt := 0
	for _, prof := range icmcProfessors {
		profDB, err := NewProfessor(prof, client, ctx)
		if err != nil {
			return -1, err
		}
		log.Printf("inserting %v%v\n", prof.Name, profDB.HashID)
		go profDB.Insert(client, ctx)
		cnt++
	}

	return cnt, nil
}

func PopulateOfferings(client *firestore.Client, ctx context.Context) (int, error) {
	log.Print("scraping all icmc departments")
	icmcProfessors := professor.ScrapeDepartments()
	cnt := 0
	for _, prof := range icmcProfessors {
		log.Printf("getting %v history\n", prof.Name)
		offerings, err := professor.GetProfessorHistory(prof.CodPes, 2010)
		if err != nil {
			return -1, err
		}

		for _, offer := range offerings {
			offerDB := NewOffering(offer)
			log.Printf("inserting %v%v\n", prof.Name, offerDB.HashID)
			go offerDB.Insert(client, ctx)
			cnt++
		}
	}

	return cnt, nil
}
