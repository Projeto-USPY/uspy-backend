package db

import (
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
)

type Manager interface {
	Insert(client *firestore.Client, ctx context.Context) error
	Hash() string
}

func Restore(client *firestore.Client, ctx context.Context, collection, HashID string) (*firestore.DocumentSnapshot, error) {
	snap, err := client.Collection(collection).Doc(HashID).Get(ctx)
	if err != nil {
		return nil, err
	}

	return snap, nil
}
