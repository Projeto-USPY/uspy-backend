package scraper

import (
	"testing"
)

func TestGetProfessorByCode(t *testing.T) {
	prof, err := getProfessorByCode(2085191) // Agma
	if err != nil {
		t.Fatal(err)
	} else if prof.Name != "Agma Juci Machado Traina" {
		t.Fail()
	}

	t.Log(prof)
}

func TestGetProfessorHistory(t *testing.T) {
	results, err := GetProfessorHistory(54946, 2010)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(results)
}
