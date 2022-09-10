package emulator

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/db/mock"
)

// ClearDatabase wipes all data from the emulator database
//
// It panics if the DELETE request fails to be created or executed
func ClearDatabase() {
	domain := os.Getenv("FIRESTORE_EMULATOR_HOST")

	if req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("http://%s/emulator/v1/projects/test/databases/(default)/documents", domain),
		nil,
	); err != nil {
		panic("could not create wipe database request: " + err.Error())
	} else {
		client := &http.Client{}
		if _, err := client.Do(req); err != nil {
			panic("could not wipe database with DELETE request: " + err.Error())
		}
	}
}

// MustGet returns the database environment for the testing emulator
//
// It is similar to get, but panics in case the environment is not able to be initialized
func MustGet() db.Database {
	// clear the database if it already exists
	ClearDatabase()

	if emu, err := Get(); err != nil {
		panic("failed to get emulator while running MustGet:" + err.Error())
	} else {
		return emu
	}
}

// Get returns the database environment for the testing emulator
func Get() (testDB db.Database, getError error) {
	testDB = db.Database{Ctx: context.Background()}

	client, err := firestore.NewClient(testDB.Ctx, "test")
	if err != nil {
		return db.Database{}, err
	}

	testDB.Client = client

	if err := mock.SetupMockData(testDB); err != nil {
		getError = err
		return
	}

	return
}
