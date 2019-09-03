package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// collection is database collection for monog db connected.
const (
	collection = "information"
)

// Mongodatabase is a Mongo collections.
type Mongodatabase struct {
	Person *mongo.Collection
}

// Connect connects to mongo db by URI,
// connectionURI URI for mongo db connetion,
// returns new mongo person collection
func Connect(ctx context.Context, connectionURI, databaseName string) (*Mongodatabase, error) {
	// setting client options.
	clientOption := options.Client().ApplyURI("mongodb://" + connectionURI)
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't connect to database using uri")
	}
	clt := client.Database(databaseName).Collection(collection)
	return &Mongodatabase{
		Person: clt,
	}, nil
}
