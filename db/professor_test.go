package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"testing"
)

func TestNewProfessorWithOfferings(t *testing.T) {
	_ = godotenv.Load("../.env")
	DB := InitFireStore("dev")
	offs := make([]entity.Offering, 0, 3)
	for i := 0; i < 3; i++ {
		off := entity.Offering{
			Semester:  0,
			Year:      0,
			Professor: 0,
			Subject:   "TestOffering" + fmt.Sprint(i),
		}

		offDB := NewOffering(off)
		offs = append(offs, off)
		err := DB.Insert(offDB)
		if err != nil {
			t.Fatal(err)
		}
		snap, err := DB.Restore("offerings", offDB.HashID)

		if err != nil {
			t.Fatal(err)
		}

		var stored OfferingDB
		err = snap.DataTo(&stored)

		if err != nil {
			t.Fatal(err)
		}

		t.Log("stored:", stored)

		if stored != *offDB {
			t.Fail()
		}
	}

	prof := entity.Professor{
		CodPes:     0,
		Name:       "testProfessor",
		Department: "test",
	}

	profDB, err := NewProfessorWithOfferings(prof, offs)
	if err != nil {
		t.Fatal(err)
	}

	err = DB.Insert(profDB)
	if err != nil {
		t.Fatal(err)
	}

	snap, err := DB.Restore("professors", profDB.HashID)

	if err != nil {
		t.Fatal(err)
	}

	var stored ProfessorDB
	err = snap.DataTo(&stored)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("stored:", stored)

	if stored.Professor != profDB.Professor {
		t.Fail()
	}
}
