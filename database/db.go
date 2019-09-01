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

//Mongodatabase is a Mongo
type mongodatabase struct {
	Collection *mongo.Collection
}


// Connect - connect to mongo db by URI,
// connectionURI URI for mongo db connetion.
func Connect(ctx context.Context, connectionURI, databaseName string) (error) {
	// setting client options.
	clientOption := options.Client().ApplyURI("mongodb://" + connectionURI)
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return fmt.Errorf(`couldn't connect to database using uri, error: %v`, err)
	}
	database.NewDatabaseConnection(client.Database(databaseName).Collection(collection))
	return nil
}

// NewDatabaseConnection asd.
func NewDatabaseConnection(coll *mongo.Collection) *mongodatabase {
	return &mongodatabase{
		Collection: coll,
	}
}

	