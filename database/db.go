package database

import (
	"context"
	"fmt"
	"log"

	"github.com/EreminDm/golang_basic_crud/crud"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SelectOne returns personal data for a given id,
// key and value params to make filtration.
func (m *mongodatabase) SelectOne(ctx context.Context, key, value string) (result *crud.PersonalData, err error) {
	val, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return nil, fmt.Errorf(`couldn't decode object id from hex err: %v`, err)
	}
	filter := bson.D{{key, val}}
	if key == `` && value == `` {
		filter = nil
	}
	err = m.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf(`find colletion error: %v`, err)
	}
	return result, nil
}

// Insert function adding data to db.
func (m *mongodatabase) Insert(ctx context.Context, document *crud.PersonalData) (interface{}, error) {
	result, err := m.Collection.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf(`adding database document error: %v`, err)
	}
	return result.InsertedID, nil
}

// SelectAll select all documents from db.
func (m *mongodatabase) SelectAll(ctx context.Context) (results *[]crud.PersonalData, err error) {
	// no filter by default.
	// Searches documents in colletion.
	cursor, err := m.Collection.Find(ctx, nil, options.Find())
	if err != nil {
		return nil, fmt.Errorf(`find collecion error: %v`, err)
	}
	defer cursor.Close(ctx)
	// Decode documents from colletion.
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("documents curser decode error: %v", err)
	}
	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("curser error: %v", err)
	}
	return results, nil
}

// Remove deletes document from Mongo.
func (m *mongodatabase) Remove(ctx context.Context, id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf(`couldn't decode object id from hex err: %v`, err)
	}
	filter := bson.D{{"_id", objectID}}
	delResult, err := m.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf(`delete document error: %v`, err)
	}
	log.Printf("deleted %v documents in collection\n", delResult.DeletedCount)

	return delResult.DeletedCount, nil
}

// Update rewrite information in db by user id filtration.
func (m *mongodatabase) Update(ctx context.Context, p *crud.PersonalData) (int64, error) {
	filter := bson.D{{"_id", p.DocumentID}}
	update := bson.D{{
		"$in", bson.D{{"name", p.Name}, {"lastName", p.LastName}, {"phone", p.Phone}, {"email", p.Email}, {"yaerOfBirth", p.YearOfBirth}},
	}}
	updateResult, err := m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf(`update %v object error: %v`, p.DocumentID, err)
	}
	return updateResult.ModifiedCount, nil
}
