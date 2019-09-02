package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connectionURI address for monog db localDB connected.
const (
	collection = "information"
)

// Mongodatabase is a Mongo collections.
type Mongodatabase struct {
	Person *mongo.Collection
}

// Connect connects to mongo db by URI,
// connectionURI URI for mongo db connetion.
func Connect(ctx context.Context, connectionURI, databaseName string) error {
	// setting client options.
	clientOption := options.Client().ApplyURI("mongodb://" + connectionURI)
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return fmt.Errorf(`couldn't connect to database using uri, error: %v`, err)
	}
	NewPersonCollection(client.Database(databaseName).Collection(collection))
	return nil
}

// NewPersonCollection returns new mongo person collection.
func NewPersonCollection(coll *mongo.Collection) *Mongodatabase {
	return &Mongodatabase{
		Person: coll,
	}
}
