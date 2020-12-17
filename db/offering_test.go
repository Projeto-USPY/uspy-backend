package db

import (
	"github.com/joho/godotenv"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"testing"
)

func TestNewOffering(t *testing.T) {
	off := entity.Offering{
		Semester:  0,
		Year:      0,
		Professor: 0,
		Subject:   "TestOffering",
	}
	offDB := NewOffering(off)
	hash := offDB.Hash()

	_ = godotenv.Load("../.env")
	DB := InitFireStore("dev")
	err := DB.Insert(offDB)
	if err != nil {
		t.Fatal(err)
	}
	snap, err := DB.Restore("offerings", hash)

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
