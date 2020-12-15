package professors

import (
	"fmt"
	"testing"
)

func TestScrapeDepartments(t *testing.T) {
	result := ScrapeDepartments()
	fmt.Println(result)
}

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
	results, err := GetProfessorHistory(2085191, 2010)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(results)
}
