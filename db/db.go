package db

import (
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
)

type Manager interface {
	Insert(client *firestore.Client, ctx context.Context) error
}
