package db

import (
	"context"

	"cloud.google.com/go/firestore"
)

func GetEmulator() (Env, error) {
	env := Env{Ctx: context.Background()}

	if client, err := firestore.NewClient(env.Ctx, "test"); err != nil {
		return Env{}, err
	} else {
		env.Client = client
	}

	return env, nil
}
