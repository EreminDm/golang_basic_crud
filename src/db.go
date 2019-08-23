package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connectionURI address for monog db localDB connected.
const (
	connectionURI = "192.168.99.100:27017"
	database      = "personal_data"
	collection    = "information"
)

// mongoClient to mongo db connections,
// mongoCollection for work with documents inside colletion.
var (
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
)

// MongodbURIConnection - connect to mongo db by URI,
// connectionURI URI for mongo db connetion.
func MongodbURIConnection() error {
	var err error
	// setting client options.
	clientOption := options.Client().ApplyURI("mongodb://" + connectionURI)
	mongoClient, err = mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Printf(`MongodbConnection func, error: %v`, err)
		return err
	}
	log.Printf("MongoDB client connected")
	mongoCollection = mongoClient.Database(database).Collection(collection)
	return nil
}

// insertPersonalData function adding data to db.
func insertPersonalData(ctx context.Context, document PersonalData) (*mongo.InsertOneResult, error) {
	result, err := mongoCollection.InsertOne(ctx, document)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

// selectAllPersonalData select all documents from db.
func selectAllPersonalData(ctx context.Context) (results *[]PersonalData, err error) {
	// no filter by default
	// Searches documents in colletion
	cursor, err := mongoCollection.Find(ctx, nil, options.Find())
	if err != nil {
		log.Printf(`Find collecion error: %v`, err)
		return nil, err
	}
	defer cursor.Close(ctx)
	// Decode documents from colletion.
	err = cursor.All(ctx, &results)
	if err != nil {
		log.Printf("Documents curser decode error: %v", err)
		return nil, err
	}
	if err := cursor.Err(); err != nil {
		log.Printf("Curser error: %v", err)
		return nil, err
	}
	return results, nil
}

// selectPersonalData select document from Mongo,
// key and value params to make filtration.
func selectPersonalData(ctx context.Context, key, value string) (result *PersonalData, err error) {
	val, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		log.Printf(`Couldn't decode object id from hex err: %v`, err)
		return nil, err
	}
	filter := bson.D{{key, val}}
	if key == `` && value == `` {
		filter = nil
	}
	err = mongoCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf(`Find colletion error: %v`, err)
		return nil, err
	}
	return result, nil
}

// deletePersonalData removes documents from Mongo.
func deletePersonalData(ctx context.Context, id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf(`Couldn't decode object id from hex err: %v`, err)
		return 0, err
	}
	filter := bson.D{{"_id", objectID}}
	delResult, err := mongoCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf(`Delete document error: %v`, err)
		return 0, err
	}
	log.Printf("Deleted %v documents in collection\n", delResult.DeletedCount)

	return delResult.DeletedCount, nil
}

// updatePersonalDataByID rewrite information in db by user id filtration.
func updatePersonalDataByID(ctx context.Context, id *primitive.ObjectID, p PersonalData) (int64, error) {
	filter := bson.D{{"_id", id}}
	update := bson.D{{
		"$in", bson.D{{"name", p.Name}, {"lastName", p.LastName}, {"phone", p.Phone}, {"email", p.Email}, {"yaerOfBirth", p.YearOfBirth}},
	}}
	updateResult, err := mongoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf(`Update %v object error: %v`, id, err)
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}
