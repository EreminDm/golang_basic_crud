package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PersonalData description.
type PersonalData struct {
	DocumentID  string `bson:"_id"` // as *primitive.ObjectID
	Name        string `bson:"name"`
	LastName    string `bson:"lastName"`
	Phone       string `bson:"phone,omitempty"`
	Email       string `bson:"email,omitempty"`
	YearOfBirth int    `bson:"yaerOfBirth,omitempty"`
}

// User abc
type User interface {
	SelectOne(ctx context.Context, value string) (PersonalData, error)
	SelectAll(ctx context.Context) (results []PersonalData, err error)
	Remove(ctx context.Context, id string) (int64, error)
	Update(ctx context.Context, p *PersonalData) (int64, error)
	Insert(ctx context.Context, document *PersonalData) (interface{}, error)
}

// SelectOne returns personal data for a given id,
// key and value params to make filtration.
func (M *Mongodatabase) SelectOne(ctx context.Context, value string) (result *PersonalData, err error) {
	val, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return nil, fmt.Errorf(`couldn't decode object id from hex err: %v`, err)
	}
	filter := bson.D{{"_id", val}}
	if value == `` {
		filter = nil
	}
	err = M.Person.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf(`find colletion error: %v`, err)
	}
	return result, nil
}

// Insert function adding data to db.
func (m *Mongodatabase) Insert(ctx context.Context, document *PersonalData) (interface{}, error) {
	result, err := m.Person.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf(`adding database document error: %v`, err)
	}
	return result.InsertedID, nil
}

// SelectAll select all documents from db.
func (m *Mongodatabase) SelectAll(ctx context.Context) (results *[]PersonalData, err error) {
	// no filter by default.
	// Searches documents in colletion.
	cursor, err := m.Person.Find(ctx, nil, options.Find())
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
func (m *Mongodatabase) Remove(ctx context.Context, id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf(`couldn't decode object id from hex err: %v`, err)
	}
	filter := bson.D{{"_id", objectID}}
	delResult, err := m.Person.DeleteOne(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf(`delete document error: %v`, err)
	}

	return delResult.DeletedCount, nil
}

// Update rewrite information in db by user id filtration.
func (m *Mongodatabase) Update(ctx context.Context, p *PersonalData) (int64, error) {
	docID, err := primitive.ObjectIDFromHex(p.DocumentID)
	if err != nil {
		return 0, err
	}
	filter := bson.D{{"_id", docID}}
	update := bson.D{{
		"$in", bson.D{{"name", p.Name}, {"lastName", p.LastName}, {"phone", p.Phone}, {"email", p.Email}, {"yaerOfBirth", p.YearOfBirth}},
	}}
	updateResult, err := m.Person.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf(`update %v object error: %v`, p.DocumentID, err)
	}
	return updateResult.ModifiedCount, nil
}
