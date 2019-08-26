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
	connectionURI = "192.168.99.100:27017"
	database      = "personal_data"
	collection    = "information"
)

// mongoClient to mongo db connections,
// mongoCollection for work with documents inside colletion.
var (
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
	err             error
)

// MongodbURIConnection - connect to mongo db by URI,
// connectionURI URI for mongo db connetion.
func MongodbURIConnection() error {
	// setting client options.
	clientOption := options.Client().ApplyURI("mongodb://" + connectionURI)
	mongoClient, err = mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		return fmt.Errorf(`MongodbConnection func, error: %v`, err)
	}
	log.Printf("MongoDB client connected")
	mongoCollection = mongoClient.Database(database).Collection(collection)
	return nil
}

// insertPersonalData function adding data to db.
func InsertPersonalData(ctx context.Context, document *PersonalData) (*mongo.InsertOneResult, error) {
	result, err := mongoCollection.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf(`DB document add error: %v`, err)
	}
	return result, nil
}

// SelectAllPersonalData select all documents from db.
func SelectAllPersonalData(ctx context.Context) (results *[]PersonalData, err error) {
	// no filter by default.
	// Searches documents in colletion.
	cursor, err := mongoCollection.Find(ctx, nil, options.Find())
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

// SelectPersonalData select document from Mongo,
// key and value params to make filtration.
func SelectPersonalData(ctx context.Context, key, value string) (result *PersonalData, err error) {
	val, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return nil, fmt.Errorf(`Couldn't decode object id from hex err: %v`, err)
	}
	filter := bson.D{{key, val}}
	if key == `` && value == `` {
		filter = nil
	}
	err = mongoCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf(`Find colletion error: %v`, err)
	}
	return result, nil
}

// DeletePersonalData removes documents from Mongo.
func DeletePersonalData(ctx context.Context, id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf(`Couldn't decode object id from hex err: %v`, err)
	}
	filter := bson.D{{"_id", objectID}}
	delResult, err := mongoCollection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf(`Delete document error: %v`, err)
	}
	log.Printf("Deleted %v documents in collection\n", delResult.DeletedCount)

	return delResult.DeletedCount, nil
}

// UpdatePersonalDataByID rewrite information in db by user id filtration.
func UpdatePersonalDataByID(ctx context.Context, p *PersonalData) (int64, error) {
	filter := bson.D{{"_id", p.DocumentID}}
	update := bson.D{{
		"$in", bson.D{{"name", p.Name}, {"lastName", p.LastName}, {"phone", p.Phone}, {"email", p.Email}, {"yaerOfBirth", p.YearOfBirth}},
	}}
	updateResult, err := mongoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf(`Update %v object error: %v`, p.DocumentID, err)
	}
	return updateResult.ModifiedCount, nil
}
