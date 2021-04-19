/* Package db contains useful functions related to the Firestore Database */
package entity

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/joho/godotenv"
	"testing"
)

func TestNewOffering(t *testing.T) {
	off := Offering{
		Semester:  0,
		Year:      0,
		Professor: 0,
		Subject:   "TestOffering",
	}

	_ = godotenv.Load("../.env")
	DB := db.InitFireStore("dev")
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
