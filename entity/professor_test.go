/* Package db contains useful functions related to the Firestore Database */
package entity

import (
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/joho/godotenv"
	"reflect"
	"testing"
)

func TestNewProfessorWithOfferings(t *testing.T) {
	_ = godotenv.Load("../.env")
	DB := db.InitFireStore("dev")
	offs := make([]Offering, 0, 3)
	for i := 0; i < 3; i++ {
		off := Offering{
			Semester:  0,
			Year:      0,
			Professor: 0,
			Subject:   "TestOffering" + fmt.Sprint(i),
		}

		offs = append(offs, off)
		err := DB.Insert(off, "offerings")
		if err != nil {
			t.Fatal(err)
		}
		snap, err := DB.Restore("offerings", off.Hash())

		if err != nil {
			t.Fatal(err)
		}

		var stored Offering
		err = snap.DataTo(&stored)

		if err != nil {
			t.Fatal(err)
		}

		t.Log("stored:", stored)

		if stored != off {
			t.Fail()
		}
	}

	prof, err := Professor{
		CodPes:     0,
		Name:       "testProfessor",
		Department: "test",
		Stats:      map[string]int{"teste": 0},
	}.WithOfferings(DB)

	if err != nil {
		t.Fatal(err)
	}

	err = DB.Insert(prof, "professors")
	if err != nil {
		t.Fatal(err)
	}

	snap, err := DB.Restore("professors", prof.Hash())

	if err != nil {
		t.Fatal(err)
	}

	var stored Professor
	err = snap.DataTo(&stored)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("stored:", stored)

	if reflect.DeepEqual(stored, prof) {
		t.Fail()
	}
}
