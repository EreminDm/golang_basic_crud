package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connectionURI address for monog db localDB connected.
const (
	collection = "information"
)

// Collection type equals mongo Collection type
type Collection = mongo.Collection

// Connect - connect to mongo db by URI,
// connectionURI URI for mongo db connetion.
func Connect(ctx context.Context, connectionURI, database string) (*Collection, error) {
	// setting client options.
	clientOption := options.Client().ApplyURI("mongodb://" + connectionURI)
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, fmt.Errorf(`couldn't connect to database using uri, error: %v`, err)
	}
	return client.Database(database).Collection(collection), nil
}

// Insert function adding data to db.
func Insert(ctx context.Context, collection *Collection, document *PersonalData) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf(`DB document add error: %v`, err)
	}
	return result, nil
}

// SelectAll select all documents from db.
func SelectAll(ctx context.Context, collection *Collection) (results *[]PersonalData, err error) {
	// no filter by default.
	// Searches documents in colletion.
	cursor, err := collection.Find(ctx, nil, options.Find())
	if err != nil {
		return nil, fmt.Errorf(`Find collecion error: %v`, err)
	}
	defer cursor.Close(ctx)
	// Decode documents from colletion.
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("Documents curser decode error: %v", err)
	}
	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("Curser error: %v", err)
	}
	return results, nil
}

// SelectOne select document from Mongo,
// key and value params to make filtration.
func SelectOne(ctx context.Context, collection *Collection, key, value string) (result *PersonalData, err error) {
	val, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return nil, fmt.Errorf(`Couldn't decode object id from hex err: %v`, err)
	}
	filter := bson.D{{key, val}}
	if key == `` && value == `` {
		filter = nil
	}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf(`Find colletion error: %v`, err)
	}
	return result, nil
}

// Remove deletes document from Mongo.
func Remove(ctx context.Context, collection *Collection, id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf(`Couldn't decode object id from hex err: %v`, err)
	}
	filter := bson.D{{"_id", objectID}}
	delResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf(`Delete document error: %v`, err)
	}
	log.Printf("Deleted %v documents in collection\n", delResult.DeletedCount)

	return delResult.DeletedCount, nil
}

// Update rewrite information in db by user id filtration.
func Update(ctx context.Context, collection *Collection, p *PersonalData) (int64, error) {
	filter := bson.D{{"_id", p.DocumentID}}
	update := bson.D{{
		"$in", bson.D{{"name", p.Name}, {"lastName", p.LastName}, {"phone", p.Phone}, {"email", p.Email}, {"yaerOfBirth", p.YearOfBirth}},
	}}
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf(`Update %v object error: %v`, p.DocumentID, err)
	}
	return updateResult.ModifiedCount, nil
}
