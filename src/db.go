package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//MongodbURIConnection - connect to mongo db by URI
//connectionURI URI for mongo db connetion
func MongodbURIConnection() error {
	var err error
	// setting client options
	clientOption := options.Client().ApplyURI("mongodb://" + ConnectionURI)
	// clientOption.Auth{
	// }
	Mongo, err = mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Printf(`MongodbConnection func, error: %v`, err)
		return err
	}
	log.Printf("MongoDB client connected")
	return nil
}

// PingMongo use for get actual connetion information
func PingMongo() error {
	//	defer cancel()
	err := Mongo.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf(`No connetion to DB: %v`, err)
		return err
	}
	return nil
}

//InsertPersonalData function adding data to DB
func InsertPersonalData(document PersonalData) (interface{}, error) {
	fmt.Println(document)
	err := PingMongo()
	if err != nil {
		err = MongodbURIConnection()
		if err != nil {
			log.Printf(`Couldn't connect to db, err: %v`, err)
			return nil, err
		}
	}
	collection := Mongo.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result.InsertedID, nil
}

// SelectAllPersonalData select all documents from
func SelectAllPersonalData() (interface{}, error) {
	err := PingMongo()
	if err != nil {
		err = MongodbURIConnection()
		if err != nil {
			log.Printf(`Couldn't connect to db, err: %v`, err)
			return nil, err
		}
	}
	//no filter by default
	filter := bson.M{}
	var (
		results []PersonalData
	)
	//get mongo collection
	collection := Mongo.Database(Database).Collection(Collection)
	ctx := context.TODO()

	//Searches documents in colletion
	cursor, err := collection.Find(ctx, filter, options.Find())
	if err != nil {
		log.Printf(`Find collecion error: %v`, err)
		return &[]PersonalData{}, err
	}
	defer cursor.Close(ctx)
	//Decode documents from colletion
	err = cursor.All(ctx, &results)
	if err != nil {
		log.Printf("Documents curser decode error: %v", err)
		return results, err
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Curser error: %v", err)
		return results, err
	}
	return results, nil
}

//SelectPersonalData select document from Mongo
//key and value params to make filtration
func SelectPersonalData(key, value string) (interface{}, error) {
	err := PingMongo()
	if err != nil {
		err = MongodbURIConnection()
		if err != nil {
			log.Printf(`Couldn't connect to db, err: %v`, err)
			return nil, err
		}
	}
	val, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		log.Printf(`Couldn't decode object id from hex err: %v`, err)
		return nil, err
	}
	filter := bson.M{key: val}
	if key == `` && value == `` {
		filter = bson.M{}
	}
	var result PersonalData
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//get mongo collection
	collection := Mongo.Database(Database).Collection(Collection)
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf(`Find colletion error: %v`, err)
		return &PersonalData{}, err
	}
	return result, nil
}

//DeletePersonalData removes documents from Mongo
func DeletePersonalData(id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf(`Couldn't decode object id from hex err: %v`, err)
		return 0, err
	}
	err = PingMongo()
	if err != nil {
		err = MongodbURIConnection()
		if err != nil {
			log.Printf(`Couldn't connect to db, err: %v`, err)
			return 0, err
		}
	}
	filter := bson.M{"_id": objectID}
	//get mongo collection
	collection := Mongo.Database(Database).Collection(Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	delResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf(`Delete document error: %v`, err)
		return 0, err
	}
	log.Printf("Deleted %v documents in collection\n", delResult.DeletedCount)

	return delResult.DeletedCount, nil
}

//UpdatePersonalDataByID rewrite information in db by user id filtration
func UpdatePersonalDataByID(id *primitive.ObjectID, p PersonalData) (int64, error) {
	err := PingMongo()
	if err != nil {
		err = MongodbURIConnection()
		if err != nil {
			log.Printf(`Couldn't connect to db, err: %v`, err)
			return 0, err
		}
	}
	//dataJSon, _ := json.Marshal(personalData)
	// dataBSON, err := bson.MarshalJSON(personalData)
	// if err != nil {
	// 	log.Printf(`Couldn't masrshal data, err: %v`, err)
	// 	return 0, err
	// }
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"name": p.Name, "lastName": p.LastName, "phone": p.Phone, "email": p.Email, "yaerOfBirth": p.YearOfBirth}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//get mongo collection
	collection := Mongo.Database(Database).Collection(Collection)
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf(`Update %v object error: %v`, id, err)
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}
